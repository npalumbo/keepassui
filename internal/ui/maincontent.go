package ui

import (
	"keepassui/internal/keepass"

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

type ToSecretReaderFn func(d DBPathAndPassword) keepass.SecretReader

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewStack(container.NewBorder(m.DBFileEntry.Container, nil, nil, nil, m.stageManager.currentViewContainer))
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	dbPathAndPassword := &DBPathAndPassword{}
	masterPasswordDialog := CreateDialog(dbPathAndPassword, parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	currentContainer := container.NewStack()
	stageManager := CreateStageManager(currentContainer)
	detailedView := CreateDetailedView(stageManager)
	addEntryView := CreateAddEntryView(stageManager)
	navView := CreateNavView(dbPathAndPassword, &addEntryView, &detailedView, parent, &stageManager, CreateKeepassSecretReaderFromDBPathAndPassword)

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

func CreateKeepassSecretReaderFromDBPathAndPassword(d DBPathAndPassword) keepass.SecretReader {
	return keepass.CipheredKeepassDB{DBBytes: d.ContentInBytes, Password: d.Password, UriID: d.UriID}
}
