package ui

import (
	"keepassui/internal/keepass"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestCreateDetailedView_Hidden(t *testing.T) {
	stageManager := CreateStageManager(container.NewStack())
	detailedView := CreateDetailedView(stageManager)
	w := test.NewWindow(container.NewWithoutLayout())
	w.SetContent(detailedView.secretForm.FormContainer)

	test.AssertImageMatches(t, "detailedView_Create.png", w.Canvas().Capture())
}

func TestUpdateDetailed_Shown(t *testing.T) {
	stageManager := CreateStageManager(container.NewStack())
	detailedView := CreateDetailedView(stageManager)
	stageManager.RegisterStager(&detailedView)
	w := test.NewWindow(container.NewWithoutLayout())
	w.SetContent(detailedView.secretForm.FormContainer)

	secretEntry := keepass.SecretEntry{
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
