package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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
	"github.com/SatisfactoryServerManager/SSMAgentManager/utils"
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

func (a *Agent) GetAgentTabItem(deleteAgentFunc func(agentname string) func()) *container.TabItem {
	return container.NewTabItem(a.Name, a.GetAgentTabContent(deleteAgentFunc))
}

func (a *Agent) GetAgentTabContent(deleteAgentFunc func(agentname string) func()) *fyne.Container {
	title := canvas.NewText("SSM Agent - "+a.Name, theme.ForegroundColor())
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

	AgentMemoryBox := customwidgets.NewNumericalEntry()
	AgentMemoryBox.SetValue(a.Memory / 1024 / 1024 / 1024)

	AgentTypeBox := widget.NewSelect([]string{"Docker", "Standalone"}, func(string) {})

	AgentTypeBox.Disable()

	if a.AgentType == "docker" {
		AgentTypeBox.SetSelected("Docker")
		AgentMemoryBox.Enable()
	} else {
		AgentTypeBox.SetSelected("Standalone")
		AgentMemoryBox.Disable()
	}

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

	deleteButton := widget.NewButtonWithIcon("Delete Agent", theme.DeleteIcon(), deleteAgentFunc(a.Name))

	deleteButton.Importance = widget.DangerImportance
	deleteButton.Move(fyne.NewPos(0, 300))
	buttonBox := container.New(mylayout.NewBlankLayout(), deleteButton)

	content := container.New(mylayout.NewFullWidthLayout(), vbox, buttonBox)
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
	agent.PortOffset = portOffset

	if agent.AgentType == "standalone" {

		if dataDirectory == "" {
			dataDirectory, _ = filepath.Abs("/SSM/data")
		}

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

	} else {
		return nil, errors.New("unknown agent type")
	}

	type newAgent struct {
		APIKey string `json:"apiKey"`
	}
	resModel := newAgent{}

	err := utils.SendPostRequest(prefs, "/api/v1/servers", agent, &resModel)
	if err != nil {
		return nil, err
	}
	agent.APIKey = resModel.APIKey

	if agent.AgentType == "standalone" {

		if utils.CheckFileExists(agent.InstallDirectory) {
			os.RemoveAll(agent.InstallDirectory)
		}

		if utils.CheckFileExists(agent.DataDirectory) {
			os.RemoveAll(agent.DataDirectory)
		}

		utils.CreateFolder(agent.InstallDirectory)
		utils.CreateFolder(agent.DataDirectory)

		err := DownloadAgent(prefs, &agent)
		if err != nil {
			return nil, err
		}

		if runtime.GOOS == "linux" {
			err := CreateLinuxAgentServiceFile(prefs, &agent)
			if err != nil {
				return nil, err
			}
		}

	} else if agent.AgentType == "docker" {

		err := PullDockerImage()
		if err != nil {
			return nil, err
		}

		err = CreateDockerContainer(prefs, &agent)
		if err != nil {
			return nil, err
		}

	}

	agent.Installed = true

	AllAgents.Agents = append(AllAgents.Agents, agent)
	SaveAgents(prefs)

	return &agent, nil

}

func DeleteAgent(AgentName string, prefs fyne.Preferences) error {
	index := -1
	var agent *Agent
	for idx := range AllAgents.Agents {
		agent = &AllAgents.Agents[idx]
		if AgentName == agent.Name {
			index = idx
			break
		}
	}

	if agent == nil {
		return errors.New("agent was not found")
	}

	if agent.AgentType == "docker" {
		err := DeleteDockerContainer(prefs, agent)
		if err != nil {
			return err
		}
	} else {
		err := os.RemoveAll(agent.InstallDirectory)
		if err != nil {
			return err
		}

		err = os.RemoveAll(agent.DataDirectory)
		if err != nil {
			return err
		}
	}

	if index > -1 {
		fmt.Printf("Deleting Agent %s\r\n", agent.Name)
		AllAgents.Agents = RemoveAgentFromArray(AllAgents.Agents, index)

		SaveAgents(prefs)
	}
	return nil
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
		"SSM_APIKEY=" + agent.APIKey,
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

func DeleteDockerContainer(prefs fyne.Preferences, agent *Agent) error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	return cli.ContainerRemove(ctx, agent.DockerID, types.ContainerRemoveOptions{RemoveVolumes: true})
}

func RemoveAgentFromArray(a []Agent, i int) []Agent {
	// Remove the element at index i from a.
	copy(a[i:], a[i+1:])  // Shift a[i+1:] left one index.
	a[len(a)-1] = Agent{} // Erase last element (write zero value).
	a = a[:len(a)-1]      // Truncate slice.

	return a
}

func DownloadAgent(prefs fyne.Preferences, agent *Agent) error {
	type gitRelease struct {
		TagVersion string `json:"tag_name"`
	}

	resModel := gitRelease{}
	err := utils.SendGetRequestURL(
		prefs,
		"https://api.github.com/repos/SatisfactoryServerManager/SSMAgent/releases/latest",
		&resModel,
	)

	if err != nil {
		return err
	}

	fmt.Printf("Latest Agent Version: %s\r\n", resModel.TagVersion)

	downloadUrl := fmt.Sprintf("https://github.com/SatisfactoryServerManager/SSMAgent/releases/download/%s/", resModel.TagVersion)
	osFile := ""
	if runtime.GOOS == "linux" {
		osFile = "SSMAgent-Linux-amd64.zip"
	} else if runtime.GOOS == "windows" {
		osFile = "SSMAgent-Windows-x64.zip"
	} else {
		return errors.New("unknown os")
	}
	downloadUrl += osFile

	err = utils.DownloadFile(prefs, downloadUrl, filepath.Join(agent.InstallDirectory, "SSMAgent.zip"))
	if err != nil {
		return err
	}
	return nil
}

func CreateLinuxAgentServiceFile(prefs fyne.Preferences, agent *Agent) error {
	if runtime.GOOS != "linux" {
		return errors.New("can only create service file on linux")
	}

	serviceContent := fmt.Sprintf(`
[Unit]
Description=SSM Agent Daemon - %s
After=network.target

[Service]
User=ssm
Group=ssm

Type=simple
WorkingDirectory=%s
ExecStart=%s/SSMAgent -name=%s -p=%d -url=%s -apikey=%s -datadir="%s"
TimeoutStopSec=20
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
	`,
		agent.Name,
		agent.InstallDirectory,
		agent.InstallDirectory,
		agent.Name,
		agent.PortOffset,
		prefs.String("ssmurl"),
		agent.APIKey,
		agent.DataDirectory,
	)

	rawServiceFile := []byte(serviceContent)

	serviceFilePath := filepath.Join("/etc/systemd/system", "SSMAgent@"+agent.Name+".service")

	err := os.WriteFile(serviceFilePath, rawServiceFile, 0777)
	if err != nil {
		return err
	}

	return nil
}
