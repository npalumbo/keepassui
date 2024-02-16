package ui_test

import (
	"keepassui/internal/secretsdb"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestCreateDetailedView_EmptyContent(t *testing.T) {
	stageManager := ui.CreateStageManager(container.NewStack())
	detailedView := ui.CreateDetailedView("", stageManager)
	w := test.NewWindow(container.NewWithoutLayout())
	w.SetContent(detailedView.GetPaintedContainer())

	w.Resize(fyne.Size{Width: 300, Height: 300})
	test.AssertImageMatches(t, "detailedView_Create.png", w.Canvas().Capture())
}

func TestShowDetails_HasContent(t *testing.T) {
	stageManager := ui.CreateStageManager(container.NewStack())
	detailedView := ui.CreateDetailedView("", stageManager)
	stageManager.RegisterStager(&detailedView)
	w := test.NewWindow(container.NewWithoutLayout())
	w.SetContent(detailedView.GetPaintedContainer())

	secretEntry := secretsdb.SecretEntry{
		Title:    "title",
		Group:    "path",
		Username: "username",
		Password: "password",
		Url:      "url",
		Notes:    "notes",
	}
	detailedView.ShowDetails(secretEntry)

	w.Resize(fyne.Size{Width: 300, Height: 300})
	test.AssertImageMatches(t, "detailedView_UpdateDetails_Shown.png", w.Canvas().Capture())
}
