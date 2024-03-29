package ui

import (
	"keepassui/internal/secretsdb"
	"keepassui/internal/secretsreader"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
)

//go:generate mockgen -destination=../mocks/ui/mock_addentryview.go -source=./addentryview.go
type EntryUpdater interface {
	AddEntry(templateEntry *secretsdb.SecretEntry)
	ModifyEntry(templateEntry *secretsdb.SecretEntry)
}

type AddEntryView struct {
	DefaultStager
	secretsReader     secretsreader.SecretReader
	SecretForm        *SecretForm
	stageManager      StagerController
	previousStageName string
	parent            fyne.Window
}

func (a *AddEntryView) GetPaintedContainer() *fyne.Container {
	return a.SecretForm.FormContainer
}

func (a *AddEntryView) GetStageName() string {
	return "AddEntry"
}

func (a *AddEntryView) AddEntry(templateEntry *secretsdb.SecretEntry) {
	addOrModify(a, templateEntry, false)
}

func addOrModify(a *AddEntryView, templateEntry *secretsdb.SecretEntry, modify bool) {
	secretForm := CreateSecretForm()
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
		originalTitle := templateEntry.Title
		originalGroup := templateEntry.Group
		originalIsGroup := templateEntry.IsGroup
		secretForm.UpdateEntry(templateEntry)
		if !modify {
			a.secretsReader.AddSecretEntry(*templateEntry)
		} else {
			a.secretsReader.ModifySecretEntry(originalTitle, originalGroup, originalIsGroup, *templateEntry)
		}

		err := a.secretsReader.Save()
		if err != nil {
			slog.Error(err.Error())
		}

		err = a.stageManager.TakeOver(a.previousStageName)
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
	addOrModify(a, templateEntry, true)
}

func CreateAddEntryView(secretsreader secretsreader.SecretReader, previousStageName string, stageManager StagerController, parent fyne.Window) AddEntryView {

	return AddEntryView{
		secretsReader:     secretsreader,
		SecretForm:        nil,
		stageManager:      stageManager,
		previousStageName: previousStageName,
		parent:            parent,
	}
}

func (a *AddEntryView) ExecuteOnTakeOver() {
	a.parent.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == mobile.KeyBack {
			err := a.stageManager.TakeOver(a.previousStageName)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	})
}
