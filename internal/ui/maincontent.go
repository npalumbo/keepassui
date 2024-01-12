package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainContent struct {
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	KeyAccordion         KeyAccordion
	LoadDBButton         *widget.Button
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewBorder(
		container.NewVBox(m.DBFileEntry.Container, m.LoadDBButton),
		container.NewVBox(widget.NewSeparator(), m.KeyAccordion.cont), nil, nil, m.KeyAccordion.accordionWidget,
	)
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	dbFileEntry := CreateDBFileEntry(parent)
	masterPasswordDialog := CreateDialog(dbFileEntry.PathBinding, dbFileEntry.ContentInBytes, parent)
	keyTable := CreatekeyAccordion(masterPasswordDialog.dbPathAndPassword, masterPasswordDialog.content, parent)
	loadFileButton := CreateLoadDBButton(masterPasswordDialog)
	masterPasswordDialog.AddListener(&keyTable)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		KeyAccordion:         keyTable,
		LoadDBButton:         loadFileButton,
	}
}

func CreateLoadDBButton(masterPasswordDialog MasterPasswordDialog) *widget.Button {
	return widget.NewButton("Load KeePass DB", func() {
		masterPasswordDialog.ShowDialog()
	})
}
