package ui

import (
	"keepassui/internal/keepass"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type MainContent struct {
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	KeyAccordion         KeyAccordion
	detailedView         DetailedView
}

type ToSecretReaderFn func(d DBPathAndPassword) keepass.SecretReader

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewBorder(
		container.NewVBox(m.DBFileEntry.Container),
		m.detailedView.container, nil, nil, m.KeyAccordion.accordionWidget,
	)
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	masterPasswordDialog := CreateDialog(parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	detailedView := CreateDetailedView()
	keyAccordion := CreatekeyAccordion(masterPasswordDialog.dbPathAndPassword, detailedView, parent, CreateKeepassSecretReaderFromDBPathAndPassword)
	masterPasswordDialog.AddListener(&keyAccordion)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		KeyAccordion:         keyAccordion,
		detailedView:         *detailedView,
	}
}

func CreateKeepassSecretReaderFromDBPathAndPassword(d DBPathAndPassword) keepass.SecretReader {
	return keepass.CipheredKeepassDB{DBBytes: d.ContentInBytes, Password: d.Password, UriID: d.UriID}
}
