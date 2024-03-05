package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContent struct {
	HomeView             HomeView
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	NavView              NavView
	StagerController     StagerController
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return m.StagerController.GetContainer()
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	masterPasswordDialog, secretsReader := CreateDialog(parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	currentContainer := container.NewStack()
	stageManager := CreateStageManager(currentContainer)
	addEntryView := CreateAddEntryView(secretsReader, "NavView", stageManager)
	navView := CreateNavView(secretsReader, &addEntryView, parent, &stageManager)
	fileSaver := CreateFileSaver(secretsReader, stageManager, parent)
	passwordConfirmView := CreatePasswordConfirmView(fileSaver, stageManager, parent)
	homeView := CreateHomeView(&dbFileEntry, stageManager, parent)

	stageManager.RegisterStager(&homeView)
	stageManager.RegisterStager(&navView)
	stageManager.RegisterStager(&addEntryView)
	stageManager.RegisterStager(&passwordConfirmView)

	masterPasswordDialog.AddListener(&navView)
	fileSaver.AddListener(&navView)

	return MainContent{
		HomeView:             homeView,
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		NavView:              navView,
		StagerController:     stageManager,
	}
}
