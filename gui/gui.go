package gui

import (
	"errors"
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SatisfactoryServerManager/SSMAgentManager/agent"
	"github.com/SatisfactoryServerManager/SSMAgentManager/customwidgets"
	"github.com/SatisfactoryServerManager/SSMAgentManager/mylayout"
)

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.RGBA{50, 50, 50, 255}
	}

	if name == theme.ColorNameDisabled {
		if variant == theme.VariantLight {
			return color.RGBA{180, 180, 180, 255}
		}
		return color.RGBA{120, 120, 120, 255}
	}

	if name == theme.ColorNameInputBorder {
		if variant == theme.VariantLight {
			return color.RGBA{100, 100, 100, 255}
		}

		return color.Black
	}

	if name == theme.ColorNameInputBackground {
		if variant == theme.VariantLight {
			return color.RGBA{220, 220, 220, 255}
		}
		return color.RGBA{20, 20, 20, 255}
	}

	if name == theme.ColorNamePrimary {
		return color.RGBA{11, 146, 204, 255}
	}

	if name == theme.ColorNameOverlayBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.RGBA{50, 50, 50, 255}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var MainWindow fyne.Window
var MainApp fyne.App

func emptyValidator(s string) (err error) {
	if s == "" {
		return errors.New("Empty String!")
	}

	return nil
}

func Init() {
	MainApp = app.NewWithID("com.ssm.ssmagentmanager")
}

func SetupGUI() {

	MainApp.Settings().SetTheme(&myTheme{})
	MainWindow = MainApp.NewWindow("SSM Agent Manager")

	tabs := container.NewAppTabs(
		container.NewTabItem("Home", BuildHomeTabContent()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	var agents = []agent.Agent{{Name: "Test"}, {Name: "DockerTest"}}
	for _, agent := range agents {
		tabs.Append(agent.GetAgentTabItem())
	}

	content := container.New(layout.NewVBoxLayout(), BuildTopBar(), tabs)

	MainWindow.SetContent(content)
	MainWindow.Resize(fyne.NewSize(800, 600))
	MainWindow.ShowAndRun()
}

func OpenCreateAgentDialog() {

	AgentNameBox := widget.NewEntry()
	AgentNameBox.SetPlaceHolder("Agent Name")
	AgentNameBox.Validator = emptyValidator
	AgentPortBox := customwidgets.NewNumericalEntry()
	AgentPortBox.Text = "0"
	AgentPortBox.SetPlaceHolder("Port Offset from 15777")

	AgentFileLocationBtn := widget.NewButtonWithIcon("Select Data Directory", theme.FolderIcon(), func() {
		log.Println("Open Data Dir dialog")
		dialog.NewFolderOpen(func(location fyne.ListableURI, err error) {
			fmt.Println(location)
		}, MainWindow).Show()
	})

	AgentFileLocationBtn.Importance = widget.HighImportance
	AgentFileLocationBtn.Disable()

	AgentMemoryBox := customwidgets.NewNumericalEntry()
	AgentMemoryBox.Text = "0"

	AgentTypeSelect := widget.NewSelect([]string{
		"Docker",
		"Standalone",
	}, func(val string) {
		switch val {
		case "Docker":
			AgentFileLocationBtn.Disable()
		case "Standalone":
			AgentFileLocationBtn.Enable()
		}
	})

	var newDialog dialog.Dialog

	formItems := []*widget.FormItem{
		{Text: "Agent Name:", Widget: AgentNameBox},
		{Text: "Agent Port Offset:", Widget: AgentPortBox},
		{Text: "Agent Type:", Widget: AgentTypeSelect},
		{Text: "Agent Data Directory:", Widget: AgentFileLocationBtn},
		{Text: "Agent Memory:", Widget: AgentMemoryBox},
	}

	log.Println("Create Agent Button Pressed")

	newDialog = dialog.NewForm("Create Agent", "Create", "Cancel", formItems, func(t bool) {}, MainWindow)
	newDialog.Resize(fyne.NewSize(500, 400))
	newDialog.Show()
}

func BuildTopBar() *fyne.Container {

	title := canvas.NewText("SSM Agent Manager", theme.PrimaryColor())
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	title.TextSize = 20
	titlePos := fyne.NewPos(10, 10)
	title.Move(titlePos)

	createAgentBtn := widget.NewButtonWithIcon("Create Agent", theme.ContentAddIcon(), OpenCreateAgentDialog)

	createAgentBtn.Importance = widget.HighImportance
	createAgentBtn.Move(fyne.NewPos(650, 10))

	topBar := container.New(&mylayout.BlankLayout{}, title, createAgentBtn)
	return topBar
}

func BuildHomeTabContent() *fyne.Container {
	SSMURLBox := widget.NewEntry()
	SSMURLBox.PlaceHolder = "https://ssmcloud.hostxtra.co.uk"
	SSMURLBox.Validator = emptyValidator

	SSMURLBox.Text = MainApp.Preferences().StringWithFallback("ssmurl", "https://ssmcloud.hostxtra.co.uk")

	SSMAPIKeyBox := widget.NewEntry()
	SSMAPIKeyBox.PlaceHolder = "API-XXXXXXXXXXXXX"
	SSMAPIKeyBox.Validator = emptyValidator
	SSMAPIKeyBox.Text = MainApp.Preferences().String("ssmapikey")

	SSMUserBox := widget.NewEntry()
	SSMUserBox.Text = MainApp.Preferences().String("ssmuser")
	SSMPassBox := widget.NewPasswordEntry()
	SSMPassBox.Text = MainApp.Preferences().String("ssmpass")

	form := widget.NewForm(
		widget.NewFormItem("SSM Cloud URL:", SSMURLBox),
		widget.NewFormItem("SSM Cloud API Key:", SSMAPIKeyBox),
		widget.NewFormItem("SSM Cloud Email:", SSMUserBox),
		widget.NewFormItem("SSM Cloud Password:", SSMPassBox),
	)
	form.OnSubmit = func() {
		MainApp.Preferences().SetString("ssmurl", SSMURLBox.Text)
		MainApp.Preferences().SetString("ssmapikey", SSMAPIKeyBox.Text)
		MainApp.Preferences().SetString("ssmuser", SSMUserBox.Text)
		MainApp.Preferences().SetString("ssmpass", SSMPassBox.Text)
	}

	form.SubmitText = "Save"

	form.Move(fyne.NewPos(0, 20))
	form.Resize(fyne.NewSize(780, form.MinSize().Height))

	formBox := container.New(mylayout.NewFullWidthLayout(), form)

	return formBox
}
