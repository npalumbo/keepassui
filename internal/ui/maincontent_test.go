package ui_test

import (
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestMainContentAppStarted(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	app := test.NewApp()
	mainContent := ui.CreateMainContent(w, app.Storage())

	w.SetContent(mainContent.MakeUI())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "mainContent_Show.png", w.Canvas().Capture())
}
