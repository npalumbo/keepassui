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
}

type ToSecretReaderFn func(d DBPathAndPassword) keepass.SecretReader

func (m *MainContent) MakeUI() fyne.CanvasObject {
	border := container.NewBorder(nil,
		m.NavView.detailedView.container, nil, nil, m.NavView.fullContainer,
	)
	return container.NewBorder(m.DBFileEntry.Container, nil, nil, nil, border)

}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	masterPasswordDialog := CreateDialog(parent)
	dbFileEntry := CreateDBFileEntry(&masterPasswordDialog, parent)
	detailedView := CreateDetailedView()
	navView := CreateNavView(masterPasswordDialog.dbPathAndPassword, detailedView, parent, CreateKeepassSecretReaderFromDBPathAndPassword)
	masterPasswordDialog.AddListener(&navView)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		NavView:              navView,
		detailedView:         *detailedView,
	}
}

func CreateKeepassSecretReaderFromDBPathAndPassword(d DBPathAndPassword) keepass.SecretReader {
	return keepass.CipheredKeepassDB{DBBytes: d.ContentInBytes, Password: d.Password, UriID: d.UriID}
}
