package ui

import (
	"keepassui/internal/keepass"

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
	d.stageManager.TakeOver(d.GetStageName())
}

func CreateDetailedView(previousStageName string, stageManager StageManager) DetailedView {
	secretForm := CreateSecretForm(true)
	secretForm.DetailsForm.OnSubmit = func() {
		stageManager.TakeOver("NavView")
	}
	secretForm.DetailsForm.SubmitText = "Back"
	secretForm.DetailsForm.Refresh()

	return DetailedView{
		stageManager:      stageManager,
		secretForm:        &secretForm,
		previousStageName: previousStageName,
	}
}
