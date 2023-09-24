package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"testing"
)

func TestMainContentAppStarted(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mainContent := CreateMainContent(w)

	w.SetContent(mainContent.MakeUI())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "appStarted.png", w.Canvas().Capture())
}
