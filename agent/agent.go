package agent

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SatisfactoryServerManager/SSMAgentManager/customwidgets"
	"github.com/SatisfactoryServerManager/SSMAgentManager/mylayout"
)

type Agent struct {
	Name             string
	PortOffset       int
	InstallDirectory string
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

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Agent Name:", Widget: AgentNameBox},
			{Text: "Agent Port Offset:", Widget: AgentPortBox},
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
