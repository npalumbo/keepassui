package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContent struct {
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	NavView              NavView
	stageManager         StageManager
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewStack(container.NewBorder(m.DBFileEntry.Container, nil, nil, nil, m.stageManager.currentViewContainer))
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	masterPasswordDialog, secretsReader := CreateDialog(parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	currentContainer := container.NewStack()
	stageManager := CreateStageManager(currentContainer)
	addEntryView := CreateAddEntryView(secretsReader, "NavView", stageManager)
	navView := CreateNavView(secretsReader, &addEntryView, parent, &stageManager)

	stageManager.RegisterStager(&navView)
	stageManager.RegisterStager(&addEntryView)

	masterPasswordDialog.AddListener(&navView)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		NavView:              navView,
		stageManager:         stageManager,
	}
}
