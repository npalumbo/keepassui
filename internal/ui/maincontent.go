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
	keyAccordion := CreatekeyAccordion(masterPasswordDialog.dbPathAndPassword, masterPasswordDialog.content, detailedView, parent, func(contentInBytes []byte, password string) keepass.SecretReader {
		return keepass.CipheredKeepassDB{ContentInBytes: contentInBytes, Password: password}
	})
	masterPasswordDialog.AddListener(&keyAccordion)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		KeyAccordion:         keyAccordion,
		detailedView:         *detailedView,
	}
}
