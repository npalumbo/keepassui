package ui

import (
	"keepassui/internal/secretsdb"
	"log/slog"

	"fyne.io/fyne/v2"
)

//go:generate mockgen -destination=../mocks/ui/mock_addentryview.go -source=./addentryview.go
type EntryUpdater interface {
	AddEntry(templateEntry *secretsdb.SecretEntry, secretsDB *secretsdb.SecretsDB)
	ModifyEntry(templateEntry *secretsdb.SecretEntry)
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

func (a *AddEntryView) AddEntry(templateEntry *secretsdb.SecretEntry, secretsDB *secretsdb.SecretsDB) {
	addOrModify(a, templateEntry, secretsDB)
}

func addOrModify(a *AddEntryView, templateEntry *secretsdb.SecretEntry, secretsDB *secretsdb.SecretsDB) {
	secretForm := CreateSecretForm(false)
	secretForm.UpdateForm(*templateEntry)
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
		if secretsDB != nil {
			secretsDB.AddSecretEntry(*templateEntry)
		}
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

func (a *AddEntryView) ModifyEntry(templateEntry *secretsdb.SecretEntry) {
	addOrModify(a, templateEntry, nil)
}

func CreateAddEntryView(previousStageName string, stageManager StagerController) AddEntryView {

	return AddEntryView{
		SecretForm:        nil,
		stageManager:      stageManager,
		previousStageName: previousStageName,
	}
}
