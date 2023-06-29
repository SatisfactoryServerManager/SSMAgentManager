package agent

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SatisfactoryServerManager/SSMAgentManager/customwidgets"
	"github.com/SatisfactoryServerManager/SSMAgentManager/mylayout"
	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	AllAgents Agents
)

type Agents struct {
	Agents []Agent `json:"agents"`
}

type Agent struct {
	Name             string `json:"name"`
	PortOffset       int    `json:"portOffset"`
	AgentType        string `json:"type"`
	Memory           int    `json:"memory"`
	InstallDirectory string `json:"installDir"`
	DataDirectory    string `json:"dataDir"`
	Installed        bool   `json:"installed"`
	DockerID         string `json:"dockerId"`
	APIKey           string `json:"apikey"`
}

func (a *Agent) GetAgentTabItem() *container.TabItem {
	return container.NewTabItem(a.Name, a.GetAgentTabContent())
}

func (a *Agent) GetAgentTabContent() *fyne.Container {
	title := canvas.NewText("SSM Agent - "+a.Name, theme.TextColor())
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	title.TextSize = 16
	title.Move(fyne.NewPos(0, 20))

	AgentNameBox := widget.NewEntry()
	AgentNameBox.Text = a.Name
	AgentNameBox.Disable()

	AgentPortBox := customwidgets.NewNumericalEntry()
	AgentPortBox.Text = strconv.Itoa(a.PortOffset)

	AgentTypeBox := widget.NewSelect([]string{"Docker", "Standalone"}, func(string) {})

	AgentTypeBox.Disable()

	if a.AgentType == "docker" {
		AgentTypeBox.SetSelected("Docker")
	} else {
		AgentTypeBox.SetSelected("Standalone")
	}

	AgentMemoryBox := customwidgets.NewNumericalEntry()
	AgentMemoryBox.SetValue(a.Memory / 1024 / 1024 / 1024)

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Agent Name:", Widget: AgentNameBox},
			{Text: "Agent Port Offset:", Widget: AgentPortBox},
			{Text: "Agent Type:", Widget: AgentTypeBox},
			{Text: "Agent Memory (GB):", Widget: AgentMemoryBox},
		},
		SubmitText: "Update",
		OnSubmit: func() { // optional, handle form submission
			a.PortOffset, _ = strconv.Atoi(AgentPortBox.Text)
		},
	}

	form.Move(fyne.NewPos(0, 20))
	form.Resize(fyne.NewSize(760, form.MinSize().Height))

	formBox := container.New(mylayout.NewFullWidthLayout(), form)

	vbox := container.New(layout.NewVBoxLayout(), title, formBox)

	vbox.Move(fyne.NewPos(10, 10))
	content := container.New(mylayout.NewFullWidthLayout(), vbox)
	return content
}

func LoadAgents(prefs fyne.Preferences) {
	agentsString := prefs.StringWithFallback("agentsJson", "{\"agents\":[]}")

	err := json.Unmarshal([]byte(agentsString), &AllAgents)
	if err != nil {
		panic(err)
	}

	SaveAgents(prefs)
}

func SaveAgents(prefs fyne.Preferences) {
	b, err := json.Marshal(AllAgents)
	if err != nil {
		panic(err)
	}
	prefs.SetString("agentsJson", string(b))
}

func CreateNewAgent(name string,
	agentType string,
	portOffset int,
	memory int,
	dataDirectory string,
	prefs fyne.Preferences,
) (*Agent, error) {

	if prefs.String("ssmurl") == "" || prefs.String("ssmapikey") == "" {
		return nil, errors.New("ssm url or ssm apikey is not set")
	}

	if !prefs.Bool("testedconnection") {
		return nil, errors.New("test connection before creating an agent")
	}

	for _, agentObj := range AllAgents.Agents {
		if agentObj.Name == name {
			return nil, errors.New("agent already exists with the same name")
		}
	}

	agent := Agent{}

	agent.Name = name
	agent.AgentType = strings.ToLower(agentType)

	if agent.AgentType == "standalone" {

		agent.DataDirectory = filepath.Join(dataDirectory, agent.Name)

		var installBaseDirectory = ""

		if runtime.GOOS == "windows" {
			installBaseDirectory, _ = filepath.Abs("C:\\Program Files\\SSM\\Agents")
		} else if runtime.GOOS == "linux" {
			installBaseDirectory, _ = filepath.Abs("/opt/SSM/Agents")
		}

		agent.InstallDirectory = filepath.Join(installBaseDirectory, agent.Name)
	} else if agent.AgentType == "docker" {
		if memory == 0 {
			return nil, errors.New("agent memory must be greater than 0")
		}

		agent.Memory = memory * 1024 * 1024 * 1024

		err := PullDockerImage()
		if err != nil {
			return nil, err
		}

		err = CreateDockerContainer(prefs, &agent)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("unknown agent type")
	}

	agent.Installed = true

	AllAgents.Agents = append(AllAgents.Agents, agent)
	SaveAgents(prefs)

	return &agent, nil

}

func PullDockerImage() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	log.Println("Pulling docker image ssmagent")
	_, err = cli.ImagePull(ctx, "docker.io/mrhid6/ssmagent:latest", types.ImagePullOptions{})
	if err != nil {
		return err
	}

	return nil

}

func CreateDockerContainer(prefs fyne.Preferences, agent *Agent) error {
	log.Println("Creating SSM Agent docker container...")

	var serverPort = 15777 + agent.PortOffset
	var beaconPort = 15000 + agent.PortOffset
	var port = 7777 + agent.PortOffset

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	var envStrings = []string{
		"SSM_NAME=" + agent.Name,
		"SSM_URL=" + prefs.String("ssmurl"),
		"SSM_APIKEY=",
	}

	var portBindings = nat.PortMap{
		"15777/udp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(serverPort),
			},
		},
		"15000/udp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(beaconPort),
			},
		},
		"7777/udp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(port),
			},
		},
	}

	var exposedPorts = nat.PortSet{
		"15777/udp": struct{}{},
		"15000/udp": struct{}{},
		"7777/udp":  struct{}{},
	}

	resp, err := cli.ContainerCreate(ctx, &dockerContainer.Config{
		Image:        "mrhid6/ssmagent:latest",
		Tty:          true,
		Env:          envStrings,
		ExposedPorts: exposedPorts,
	}, &dockerContainer.HostConfig{
		PortBindings: portBindings,
	}, nil, nil, agent.Name)

	if err != nil {
		return err
	}
	agent.DockerID = resp.ID

	log.Println("SSM Agent Docker container created successfully")
	return nil
}
