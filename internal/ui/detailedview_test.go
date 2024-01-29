package ui

import (
	"keepassui/internal/keepass"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestCreateDetailedView_Hidden(t *testing.T) {
	detailedView := CreateDetailedView()
	w := test.NewWindow(container.NewWithoutLayout())
	w.SetContent(detailedView.container)

	test.AssertImageMatches(t, "detailedView_Create.png", w.Canvas().Capture())
}

func TestUpdateDetailes_Shown(t *testing.T) {
	detailedView := CreateDetailedView()
	w := test.NewWindow(container.NewWithoutLayout())
	w.SetContent(detailedView.container)

	secretEntry := keepass.SecretEntry{
		Title:    "title",
		Group:    "path",
		Username: "username",
		Password: "password",
		Url:      "url",
		Notes:    "notes",
	}
	detailedView.UpdateDetails(secretEntry)

	w.Resize(fyne.Size{Width: 300, Height: 300})
	test.AssertImageMatches(t, "detailedView_UpdateDetails_Shown.png", w.Canvas().Capture())
}
