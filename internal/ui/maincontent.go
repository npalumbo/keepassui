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

type SecretReaderResolver interface {
	GetSecretReader(d DBPathAndPassword) secretsreader.SecretReader
}

type SecretReaderFactory struct {
}

var DefaultSecretReaderFactory SecretReaderFactory = SecretReaderFactory{}

func (s SecretReaderFactory) GetSecretReader(d DBPathAndPassword) secretsreader.SecretReader {
	return secretsreader.CipheredKeepassDB{DBBytes: d.ContentInBytes, Password: d.Password, UriID: d.UriID}
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewStack(container.NewBorder(m.DBFileEntry.Container, nil, nil, nil, m.stageManager.currentViewContainer))
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	dbPathAndPassword := &DBPathAndPassword{}
	masterPasswordDialog := CreateDialog(dbPathAndPassword, parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	currentContainer := container.NewStack()
	stageManager := CreateStageManager(currentContainer)
	detailedView := CreateDetailedView("NavView", stageManager)
	addEntryView := CreateAddEntryView("NavView", stageManager)
	navView := CreateNavView(dbPathAndPassword, &addEntryView, &detailedView, parent, &stageManager, DefaultSecretReaderFactory)

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
