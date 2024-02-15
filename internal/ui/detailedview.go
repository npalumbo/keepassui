package ui

import (
	"keepassui/internal/keepass"
	"log/slog"

	"fyne.io/fyne/v2"
)

type DetailedView struct {
	DefaultStager
	secretForm        *SecretForm
	stageManager      StageManager
	previousStageName string
}

func (d *DetailedView) GetPaintedContainer() *fyne.Container {
	return d.secretForm.FormContainer
}

func (d *DetailedView) GetStageName() string {
	return "DetailedView"
}

func (d *DetailedView) ShowDetails(secretEntry keepass.SecretEntry) {
	d.secretForm.UpdateForm(secretEntry)
	err := d.stageManager.TakeOver(d.GetStageName())
	if err != nil {
		slog.Error(err.Error())
	}
}

func CreateDetailedView(previousStageName string, stageManager StageManager) DetailedView {
	secretForm := CreateSecretForm(true)
	secretForm.DetailsForm.OnSubmit = func() {
		err := stageManager.TakeOver(previousStageName)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	secretForm.DetailsForm.SubmitText = "Back"
	secretForm.DetailsForm.Refresh()

	return DetailedView{
		stageManager:      stageManager,
		secretForm:        &secretForm,
		previousStageName: previousStageName,
	}
}
