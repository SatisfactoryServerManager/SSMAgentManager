package gui

import (
	"errors"
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
var MainTabs *container.AppTabs

func emptyValidator(s string) (err error) {
	if s == "" {
		return errors.New("Empty String!")
	}

	return nil
}

func Init() {
	MainApp = app.NewWithID("com.ssm.ssmagentmanager")
}

func SetupTabs() {
	MainTabs = container.NewAppTabs(
		container.NewTabItem("Home", BuildHomeTabContent()),
	)

	RefreshTabs()
}

func RefreshTabs() {
	var tabItems = []*container.TabItem{
		container.NewTabItem("Home", BuildHomeTabContent()),
	}

	for _, agent := range agent.AllAgents.Agents {
		tabItems = append(tabItems, agent.GetAgentTabItem())
	}

	MainTabs.SetItems(tabItems)

	MainTabs.Refresh()
}

func SetupGUI() {

	MainApp.Settings().SetTheme(&myTheme{})
	MainWindow = MainApp.NewWindow("SSM Agent Manager")

	SetupTabs()

	content := container.New(layout.NewVBoxLayout(), BuildTopBar(), MainTabs)

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

	var agentDataDir = ""
	AgentFileLocationBtn := widget.NewButtonWithIcon("Select Data Directory", theme.FolderIcon(), func() {
		log.Println("Open Data Dir dialog")
		dialog.NewFolderOpen(func(location fyne.ListableURI, err error) {
			if location != nil {
				agentDataDir = location.Path()
			}
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
		{Text: "Agent Memory (GB):", Widget: AgentMemoryBox},
	}

	log.Println("Create Agent Button Pressed")

	newDialog = dialog.NewForm("Create Agent", "Create", "Cancel", formItems, func(t bool) {
		if t {
			portOffset, _ := AgentPortBox.GetValue()
			memory, _ := AgentMemoryBox.GetValue()

			_, err := agent.CreateNewAgent(
				AgentNameBox.Text,
				AgentTypeSelect.Selected,
				portOffset,
				memory,
				agentDataDir,
				MainApp.Preferences(),
			)

			if err != nil {
				errorDiag := dialog.NewError(err, MainWindow)
				errorDiag.Show()
				log.Printf("Error creating agent, with error %s\r\n", err.Error())
				return
			}

			RefreshTabs()
		}
	}, MainWindow)

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

	var testBtn *widget.Button

	SSMURLBox := widget.NewEntry()
	SSMURLBox.PlaceHolder = "https://ssmcloud.hostxtra.co.uk"
	SSMURLBox.Validator = emptyValidator
	SSMURLBox.Text = MainApp.Preferences().StringWithFallback("ssmurl", "https://ssmcloud.hostxtra.co.uk")

	SSMAPIKeyBox := widget.NewEntry()
	SSMAPIKeyBox.PlaceHolder = "API-XXXXXXXXXXXXX"
	SSMAPIKeyBox.Validator = emptyValidator
	SSMAPIKeyBox.Text = MainApp.Preferences().String("ssmapikey")

	connectionDetailsText := widget.NewRichTextWithText("Enter the SSM Cloud connection details\nSSM URL - The SSM Cloud URL\nSSM API Key - The SSM User API Key, used to create new agents")
	connectionDetailsText.Move(fyne.NewPos(10, 10))

	form := widget.NewForm(
		widget.NewFormItem("SSM Cloud URL:", SSMURLBox),
		widget.NewFormItem("SSM Cloud API Key:", SSMAPIKeyBox),
	)

	form.SubmitText = "Save"
	form.OnSubmit = func() {
		MainApp.Preferences().SetString("ssmurl", SSMURLBox.Text)
		MainApp.Preferences().SetString("ssmapikey", SSMAPIKeyBox.Text)
		MainApp.Preferences().SetBool("testedconnection", false)
		testBtn.Enable()
	}

	form.Move(fyne.NewPos(0, 110))
	form.Resize(fyne.NewSize(780, form.MinSize().Height))

	seperator := widget.NewSeparator()
	seperator.Move(fyne.NewPos(0, 80))
	seperator.Resize(fyne.NewSize(100, 5))

	seperator2 := widget.NewSeparator()
	seperator2.Move(fyne.NewPos(0, 280))
	seperator2.Resize(fyne.NewSize(100, 5))

	testText := widget.NewRichTextWithText("Using the button below, test the connection to SSM Cloud")
	testText.Move(fyne.NewPos(10, 290))

	testBtn = widget.NewButtonWithIcon("Test Connection", theme.MediaReplayIcon(), func() {
		MainApp.Preferences().SetBool("testedconnection", true)
		testBtn.Disable()
	})
	testBtn.Importance = widget.HighImportance
	testBtn.Move(fyne.NewPos(20, 330))

	if MainApp.Preferences().Bool("testedconnection") {
		testBtn.Disable()
	} else {
		testBtn.Enable()
	}

	buttonBox := container.New(mylayout.NewBlankLayout(), testBtn)

	formBox := container.New(mylayout.NewFullWidthLayout(), connectionDetailsText, seperator, form, seperator2, testText, buttonBox)

	return formBox
}
