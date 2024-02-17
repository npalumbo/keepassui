package ui

import (
	"keepassui/internal/secretsreader"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContent struct {
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	NavView              NavView
	detailedView         DetailedView
	stageManager         StageManager
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewStack(container.NewBorder(m.DBFileEntry.Container, nil, nil, nil, m.stageManager.currentViewContainer))
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	dbPathAndPassword := &secretsreader.DBPathAndPassword{}
	masterPasswordDialog := CreateDialog(dbPathAndPassword, parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	currentContainer := container.NewStack()
	stageManager := CreateStageManager(currentContainer)
	detailedView := CreateDetailedView("NavView", stageManager)
	addEntryView := CreateAddEntryView(dbPathAndPassword, "NavView", stageManager)
	navView := CreateNavView(dbPathAndPassword, &addEntryView, &detailedView, parent, &stageManager)

	stageManager.RegisterStager(&navView)
	stageManager.RegisterStager(&addEntryView)
	stageManager.RegisterStager(&detailedView)

	masterPasswordDialog.AddListener(&navView)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		NavView:              navView,
		detailedView:         detailedView,
		stageManager:         stageManager,
	}
}
