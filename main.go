package main

import (
	"keepassui/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {

	a := app.NewWithID("com.keepassui")
	w := a.NewWindow("Keepass UI")
	mainContent := ui.CreateMainContent(w, a.Storage())

	w.SetContent(mainContent.MakeUI())
	w.Resize(fyne.NewSize(600, 600))

	w.ShowAndRun()
}
