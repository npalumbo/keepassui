package ui

import (
	"keepassui/internal/keepass"
	"log/slog"

	"fyne.io/fyne/v2"
)

//go:generate mockgen -destination=../mocks/ui/mock_addentryview.go -source=./addentryview.go
type EntryUpdater interface {
	AddEntry(templateEntry *keepass.SecretEntry, secretsDB *keepass.SecretsDB)
}

type AddEntryView struct {
	DefaultStager
	SecretForm        *SecretForm
	stageManager      StagerController
	previousStageName string
}

func (a *AddEntryView) GetPaintedContainer() *fyne.Container {
	return a.SecretForm.FormContainer
}

func (a *AddEntryView) GetStageName() string {
	return "AddEntry"
}

func (a *AddEntryView) AddEntry(templateEntry *keepass.SecretEntry, secretsDB *keepass.SecretsDB) {
	secretForm := CreateSecretForm(false)
	a.SecretForm = &secretForm
	secretForm.DetailsForm.OnCancel = func() {
		err := a.stageManager.TakeOver(a.previousStageName)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	secretForm.DetailsForm.Refresh()
	secretForm.DetailsForm.OnSubmit = func() {
		secretForm.UpdateEntry(templateEntry)
		secretsDB.AddSecretEntry(*templateEntry)
		err := a.stageManager.TakeOver(a.previousStageName)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	err := a.stageManager.TakeOver(a.GetStageName())
	if err != nil {
		slog.Error(err.Error())
	}
}

func CreateAddEntryView(previousStageName string, stageManager StagerController) AddEntryView {

	return AddEntryView{
		SecretForm:        nil,
		stageManager:      stageManager,
		previousStageName: previousStageName,
	}
}
