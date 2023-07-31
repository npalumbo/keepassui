package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"keepassui/pkg/ui"
)

func main() {

	a := app.New()
	w := a.NewWindow("Keepass UI")
	mainView := &ui.MainContent{}
	w.SetContent(mainView.MakeUI(w))
	w.Resize(fyne.NewSize(600, 600))

	w.ShowAndRun()
}
