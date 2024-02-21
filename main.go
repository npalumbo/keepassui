package main

import (
	"keepassui/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func main() {

	a := app.NewWithID("com.keepassui")
	w := a.NewWindow("Keepass UI")
	mainContent := ui.CreateMainContent(w, a.Storage())

	w.SetContent(mainContent.MakeUI())
	w.Resize(fyne.NewSize(600, 600))
	err := mainContent.StagerController.TakeOver(mainContent.HomeView.GetStageName())
	if err != nil {
		dialog.ShowError(err, w)
	}
	w.ShowAndRun()
}
