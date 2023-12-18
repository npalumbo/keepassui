package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainContent struct {
	DBFileEntry          DBFileEntry
	MasterPasswordDialog MasterPasswordDialog
	KeyList              KeyList
	LoadDBButton         *widget.Button
}

func (m *MainContent) MakeUI() fyne.CanvasObject {
	return container.NewVBox(
		container.NewVBox(m.DBFileEntry.Container, m.LoadDBButton),
		widget.NewSeparator(),
		m.KeyList.listWidget)
}

func CreateMainContent(parent fyne.Window, stor fyne.Storage) MainContent {
	dbFileEntry := CreateDBFileEntry(parent)
	masterPasswordDialog := CreateDialog(dbFileEntry.PathBinding, dbFileEntry.ContentInBytes, parent)
	keyList := CreatekeyList(masterPasswordDialog.dbPathAndPassword, masterPasswordDialog.content, parent)
	loadFileButton := CreateLoadDBButton(masterPasswordDialog)
	masterPasswordDialog.AddListener(&keyList)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		KeyList:              keyList,
		LoadDBButton:         loadFileButton,
	}
}

func CreateLoadDBButton(masterPasswordDialog MasterPasswordDialog) *widget.Button {
	return widget.NewButton("Load KeePass DB", func() {
		masterPasswordDialog.ShowDialog()
	})
}
