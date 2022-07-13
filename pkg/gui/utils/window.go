package utils

import (
	"DomainMan/pkg/gui/instance"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowSimpleWindow(inst *instance.Instance, title string, content fyne.CanvasObject) fyne.Window {
	window := inst.App.NewWindow(title)
	window.SetContent(content)
	window.CenterOnScreen()
	window.Show()
	return window
}

func ShowSimpleWindowComputed(inst *instance.Instance, title string, contentFunc func(window fyne.Window) fyne.CanvasObject) fyne.Window {
	window := inst.App.NewWindow(title)
	window.SetContent(contentFunc(window))
	window.CenterOnScreen()
	window.Show()
	return window
}

func ShowMessage(inst *instance.Instance, title string, message string) {
	window := inst.App.NewWindow(title)
	window.SetContent(container.NewVBox(widget.NewLabel(message), widget.NewButton("OK", func() {
		window.Close()
	})))
	window.CenterOnScreen()
	window.Show()
}

func ShowSizedWindow(inst *instance.Instance, title string, content fyne.CanvasObject, width, height float32) fyne.Window {
	window := inst.App.NewWindow(title)
	window.SetContent(content)
	window.Resize(fyne.NewSize(width, height))
	window.CenterOnScreen()
	window.Show()
	return window
}

func ShowSizedWindowComputed(inst *instance.Instance, title string, contentFunc func(window fyne.Window) fyne.CanvasObject, width, height float32) fyne.Window {
	window := inst.App.NewWindow(title)
	window.SetContent(contentFunc(window))
	window.Resize(fyne.NewSize(width, height))
	window.CenterOnScreen()
	window.Show()
	return window
}
