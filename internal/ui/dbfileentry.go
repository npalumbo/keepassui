package ui

import (
	"io"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DBFileEntry struct {
	Container      *fyne.Container
	PathBinding    binding.String
	ContentInBytes *[]byte
	findFileButton *widget.Button
	fileOpenDialog *dialog.FileDialog
}

func CreateDBFileEntry(parent fyne.Window) DBFileEntry {
	var byteContent []byte
	pathBinding := binding.NewString()

	fileOpen := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
		if err == nil && dir != nil {
			err = pathBinding.Set(dir.URI().Path())
			if err == nil {
				byteContent, err = io.ReadAll(dir)
			}
			if err != nil {
				slog.Error("Error setting path: %s", dir.URI().Path(), err)
			}
		}
	}, parent)
	fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	findFileButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		fileOpen.Show()
	})

	kdbxFilePathEntry := widget.NewEntryWithData(pathBinding)
	kdbxFilePathEntry.Resize(fyne.NewSize(700, kdbxFilePathEntry.Size().Height))
	kdbxFilePathEntry.PlaceHolder = "Path to db.kdbx"

	return DBFileEntry{
		Container:      container.NewBorder(nil, nil, nil, findFileButton, kdbxFilePathEntry),
		PathBinding:    pathBinding,
		ContentInBytes: &byteContent,
		findFileButton: findFileButton,
		fileOpenDialog: fileOpen,
	}
}
