package containers

import (
	"DomainMan/pkg/api"
	"DomainMan/pkg/database"
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/gui/utils"
	"DomainMan/pkg/models"
	"DomainMan/pkg/mq"
	"DomainMan/pkg/random"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/http"
	"time"
)

func initWebAPIButton(inst *instance.Instance) (webAPILabel *widget.Label, webAPIButton *widget.Button) {
	webAPILabel = widget.NewLabel("Web API")
	webAPIStarted := false
	onWebAPIStarting := func() {
		webAPIButton.Disable()
		webAPIButton.SetText("Starting...")
	}
	onWebAPIStarted := func() {
		webAPIButton.Enable()
		webAPIButton.SetText("Stop")
		webAPIButton.SetIcon(theme.MediaStopIcon())
		webAPILabel.SetText(fmt.Sprintf("Web API  (Listening: %s)", inst.Settings.WebAPI.Listen))
		webAPIStarted = true
	}
	onWebAPIStopping := func() {
		webAPIButton.Disable()
		webAPIButton.SetText("Stopping...")
	}
	onWebAPIStopped := func() {
		webAPIButton.Enable()
		webAPIButton.SetText("Start")
		webAPIButton.SetIcon(theme.MediaPlayIcon())
		webAPILabel.SetText("Web API")
		webAPIStarted = false
	}
	webAPIButton = widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func() {
		if !webAPIStarted {
			onWebAPIStarting()
			mq.Init()
			api.Init()
			go func() {
				err := api.Run(inst.Settings.WebAPI.Listen)
				if err != http.ErrServerClosed && err != nil {
					utils.ShowMessage(inst, "Web API Error", fmt.Sprintf("Error when starting web API: %v", err))
					onWebAPIStopped()
				}
			}()
			onWebAPIStarted()
		} else {
			onWebAPIStopping()
			_ = mq.MQ.Stop()
			err := api.Stop()
			if err != nil {
				utils.ShowMessage(inst, "Web API Error", fmt.Sprintf("Error when stopping web API: %v", err))
				onWebAPIStarted()
			} else {
				onWebAPIStopped()
			}
		}
	})
	return
}

func initWebAPIAccessTokenGenerateButton(inst *instance.Instance) (webAPIAccessTokenCreateButton *widget.Button) {
	webAPIAccessTokenCreateButton = widget.NewButtonWithIcon("Generate", theme.ContentAddIcon(), func() {
		accessToken := models.AccessToken{
			Name:  fmt.Sprintf("dm-gui@%s", time.Now().Format(time.RFC3339Nano)),
			Token: random.String(models.AccessTokenLength, random.DictAlphaNumber),
		}
		if v := database.DB.Create(&accessToken); v.Error != nil {
			utils.ShowMessage(inst, "Database Error", fmt.Sprintf("Error when creating access token: %v", v.Error))
		} else {
			utils.ShowSimpleWindowComputed(inst, "Web API Access Token", func(window fyne.Window) fyne.CanvasObject {
				var copyButton *widget.Button
				copyButton = widget.NewButton("Copy", func() {
					window.Clipboard().SetContent(accessToken.Token)
				})
				copyButton.Importance = widget.HighImportance
				return container.NewVBox(
					widget.NewLabel(fmt.Sprintf("Access Token: %s", accessToken.Token)),
					copyButton,
					widget.NewButton("Close", func() {
						window.Close()
					}),
				)
			})
		}
	})
	return
}

func SettingsRoot(inst *instance.Instance) fyne.CanvasObject {
	title := canvas.NewText("Settings", theme.ForegroundColor())
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24

	webAPILabel, webAPIButton := initWebAPIButton(inst)
	webAPIAccessTokenCreateButton := initWebAPIAccessTokenGenerateButton(inst)
	return container.NewVBox(
		title,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Database Settings"),
			widget.NewButtonWithIcon("Edit", theme.NavigateNextIcon(), func() {
				NewSettingsDatabaseWindow(inst, false).Show()
			}),
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			webAPILabel,
			webAPIButton,
			widget.NewLabel("Web API Settings"),
			widget.NewButtonWithIcon("Edit", theme.NavigateNextIcon(), func() {
				NewSettingsWebAPIWindow(inst).Show()
			}),
			widget.NewLabel("Web API Access Token"),
			webAPIAccessTokenCreateButton,
		),
		widget.NewSeparator(),
	)
}
