package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log/slog"
)

type DBFileEntry struct {
	Container   *fyne.Container
	PathBinding binding.String
}

func CreateDBFileEntry(parent fyne.Window) DBFileEntry {
	pathBinding := binding.NewString()
	findFileButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		fileOpen := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err == nil && dir != nil {
				err = pathBinding.Set(dir.URI().Path())
				if err != nil {
					slog.Error("Error setting path: %s", dir.URI().Path(), err)
				}
			}
		}, parent)
		fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
		fileOpen.Show()
	})

	kdbxFilePathEntry := widget.NewEntryWithData(pathBinding)
	kdbxFilePathEntry.Resize(fyne.NewSize(700, kdbxFilePathEntry.Size().Height))
	kdbxFilePathEntry.PlaceHolder = "Path to db.kdbx"

	return DBFileEntry{
		Container:   container.NewBorder(nil, nil, nil, findFileButton, kdbxFilePathEntry),
		PathBinding: pathBinding,
	}
}
