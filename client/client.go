package main

// a simple cydrive client only for test

import (
	"fmt"
	"fyne.io/fyne"
	fyneApp "fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/yah01/CyDrive/model"
	"net/http"
	"net/http/cookiejar"
)

var (
	cookieJar *cookiejar.Jar
	client    *http.Client
	user      = model.User{
		Username: "test",
		Password: "testCyDrive",
	}

	baseUrl string
)

func init() {
	cookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           cookieJar,
		Timeout:       0,
	}
}

var serverAddress = "127.0.0.1"

var (
	app         = fyneApp.New()
	window      = app.NewWindow("CyDrive")
	fileList    = fyne.NewContainer()
	loginButton = widget.NewButton("Login", func() {
		Login(user.Username, user.Password)
	})
	listButton = widget.NewButton("List", func() {
		ListRemoteDir("")
	})

	taskListTab  = widget.NewTabItem("Task", widget.NewLabel("Task List"))
	driveTab     = widget.NewTabItem("Drive", widget.NewLabel("File List"))
	settingTab   = widget.NewTabItem("Setting", widget.NewLabel("Setting Items"))
	tabContainer = widget.NewTabContainer(taskListTab, driveTab, settingTab)
)

func main() {
	tabContainer.SetTabLocation(widget.TabLocationBottom)
	baseUrl = fmt.Sprintf("http://%s:6454", serverAddress)
	app.Settings().SetTheme(theme.LightTheme())
	Login(user.Username, user.Password)
	ListRemoteDir("")
	tabContainer.SelectTab(driveTab)

	window.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widget.NewLabelWithStyle("CyDrive",
			fyne.TextAlignCenter, fyne.TextStyle{
				Bold:   true,
				Italic: true,
			}),

		layout.NewSpacer(),

		fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(4),
			fileList.Objects...,
		),

		layout.NewSpacer(),

		fyne.NewContainerWithLayout(layout.NewCenterLayout(),
			tabContainer,
		),
	))

	window.ShowAndRun()
}
