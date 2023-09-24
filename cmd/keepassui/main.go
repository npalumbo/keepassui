package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"keepassui/internal/ui"
)

func main() {

	a := app.New()
	w := a.NewWindow("Keepass UI")
	mainContent := ui.CreateMainContent(w)

	w.SetContent(mainContent.MakeUI())
	w.Resize(fyne.NewSize(600, 600))

	w.ShowAndRun()
}
