package ui

import (
	"keepassui/internal/keepass"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainContent struct {
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	KeyAccordion         KeyAccordion
	detailedView         DetailedView
	LoadDBButton         *widget.Button
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewBorder(
		container.NewVBox(m.DBFileEntry.Container, m.LoadDBButton),
		m.detailedView.container, nil, nil, m.KeyAccordion.accordionWidget,
	)

}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	dbFileEntry := CreateDBFileEntry(parent)
	masterPasswordDialog := CreateDialog(dbFileEntry.PathBinding, dbFileEntry.ContentInBytes, parent)
	detailedView := CreateDetailedView()
	keyAccordion := CreatekeyAccordion(masterPasswordDialog.dbPathAndPassword, masterPasswordDialog.content, detailedView, parent, func(contentInBytes []byte, password string) keepass.SecretReader {
		return keepass.CipheredKeepassDB{ContentInBytes: contentInBytes, Password: password}
	})
	loadFileButton := CreateLoadDBButton(masterPasswordDialog)
	masterPasswordDialog.AddListener(&keyAccordion)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		KeyAccordion:         keyAccordion,
		detailedView:         *detailedView,
		LoadDBButton:         loadFileButton,
	}
}

func CreateLoadDBButton(masterPasswordDialog MasterPasswordDialog) *widget.Button {
	return widget.NewButton("Load KeePass DB", func() {
		masterPasswordDialog.ShowDialog()
	})
}
