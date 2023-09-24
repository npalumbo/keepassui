package ui

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
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

func CreateMainContent(parent fyne.Window) MainContent {
	dbFileEntry := CreateDBFileEntry(parent)
	masterPasswordDialog := CreateDialog(dbFileEntry.PathBinding, parent)
	keyList := CreatekeyList(masterPasswordDialog.dbPathAndPassword, parent)
	loadFileButton := CreateLoadDBButton(dbFileEntry.PathBinding, masterPasswordDialog, parent)
	masterPasswordDialog.AddListener(&keyList)

	return MainContent{
		DBFileEntry:          dbFileEntry,
		MasterPasswordDialog: masterPasswordDialog,
		KeyList:              keyList,
		LoadDBButton:         loadFileButton,
	}
}

func CreateLoadDBButton(pathBinding binding.String, masterPasswordDialog MasterPasswordDialog, parent fyne.Window) *widget.Button {
	return widget.NewButton("Load KeePass DB", func() {
		filePath, err := pathBinding.Get()
		if err != nil {
			dialog.ShowError(err, parent)
		} else if !fileExists(filePath) {
			dialog.ShowError(errors.New("file does not exist"), parent)
			return
		} else {
			masterPasswordDialog.ShowDialog()
		}
	})
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
