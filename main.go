package main

import (
	"embed"
	"keepassui/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
)

//go:embed translation
var translations embed.FS

func main() {
	err := lang.AddTranslationsFS(translations, "translation")
	if err != nil {
		panic("Error initialising translations")
	}

	a := app.NewWithID("com.keepassui")
	w := a.NewWindow(lang.L("Window Title"))
	mainContent := ui.CreateMainContent(w, a.Storage())

	w.SetContent(mainContent.MakeUI())
	w.Resize(fyne.NewSize(600, 600))
	err = mainContent.StagerController.TakeOver(mainContent.HomeView.GetStageName())
	if err != nil {
		dialog.ShowError(err, w)
	}
	w.ShowAndRun()
}
