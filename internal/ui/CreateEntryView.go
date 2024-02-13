package ui

import (
	"fyne.io/fyne/v2"
	"keepassui/internal/keepass"
)

type AddEntryView struct {
	DefaultStager
	secretForm   *SecretForm
	stageManager StageManager
}

func (a *AddEntryView) GetPaintedContainer() *fyne.Container {
	return a.secretForm.FormContainer
}

func (a *AddEntryView) GetStageName() string {
	return "AddEntry"
}

func (a *AddEntryView) AddEntry(templateEntry *keepass.SecretEntry, secretsDB *keepass.SecretsDB) {
	secretForm := CreateForm(false)
	a.secretForm = &secretForm
	secretForm.detailsForm.OnCancel = func() {
		a.stageManager.TakeOver("NavView")
	}
	secretForm.detailsForm.Refresh()
	secretForm.detailsForm.OnSubmit = func() {
		secretForm.UpdateEntry(templateEntry)
		secretsDB.AddSecretEntry(*templateEntry)
		a.stageManager.TakeOver("NavView")
	}
	a.stageManager.TakeOver(a.GetStageName())
}

func CreateAddEntryView(stageManager StageManager) AddEntryView {

	return AddEntryView{
		secretForm:   nil,
		stageManager: stageManager,
	}
}
