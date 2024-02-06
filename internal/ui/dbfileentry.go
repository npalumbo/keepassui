package ui

import (
	"io"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DBFileEntry struct {
	Container      *fyne.Container
	findFileButton *widget.Button
	fileOpenDialog *dialog.FileDialog
}

func CreateDBFileEntry(masterPasswordDialog *MasterPasswordDialog, parent fyne.Window) DBFileEntry {
	var byteContent []byte

	fileOpen := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
		if err == nil && dir != nil {

			fileURI := dir.URI()
			byteContent, err = io.ReadAll(dir)
			defer dir.Close()
			if err == nil {
				go (*masterPasswordDialog).ShowDialog(fileURI.String(), &byteContent)
			}
			if err != nil {
				slog.Error("Error setting path: %s", fileURI.Path(), err)
			}
		}
	}, parent)
	fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	findFileButton := widget.NewButtonWithIcon("Load Keepass file", theme.SearchIcon(), func() {
		fileOpen.Show()
	})

	return DBFileEntry{
		Container:      container.NewStack(findFileButton),
		findFileButton: findFileButton,
		fileOpenDialog: fileOpen,
	}
}
