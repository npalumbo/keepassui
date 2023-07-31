package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MainContent struct {
	dbFileEntry          DBFileEntry
	loadFileButton       *widget.Button
	masterPasswordDialog MasterPasswordDialog
	keyList              KeyList
}

func (dbOpener *MainContent) MakeUI(parent fyne.Window) fyne.CanvasObject {

	dbOpener.dbFileEntry = CreateDBFileEntry(parent)
	dbOpener.masterPasswordDialog = CreateDialog(dbOpener.dbFileEntry.PathBinding)
	dbOpener.keyList = CreatekeyList(dbOpener.masterPasswordDialog.dbPathAndPassword)
	dbOpener.masterPasswordDialog.AddListener(&dbOpener.keyList)
	dbOpener.loadFileButton = widget.NewButton("Load KeePass DB", func() {
		dbOpener.masterPasswordDialog.ShowDialog(parent)
	})

	return container.NewVBox(
		container.NewVBox(dbOpener.dbFileEntry.Container, dbOpener.loadFileButton),
		widget.NewSeparator(),
		container.NewMax(dbOpener.keyList.listWidget))
}
