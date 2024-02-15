package ui_test

import (
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
)

func TestCreateDBFileEntry(t *testing.T) {
	mainContainer := container.NewVBox()
	w := test.NewWindow(mainContainer)

	dbFileEntry := ui.CreateDBFileEntry(nil, w)

	mainContainer.Add(dbFileEntry.Container)
	w.Resize(fyne.NewSize(600, 600))

	w.ShowAndRun()

	test.AssertImageMatches(t, "createDBFileEntry.png", w.Canvas().Capture())

	test.Tap(dbFileEntry.FindFileButton)

	URI, err := storage.ParseURI("file://testdata//files")
	if err != nil {
		t.FailNow()
	}
	listableURI, err := storage.ListerForURI(URI)
	if err != nil {
		t.FailNow()
	}
	dbFileEntry.FileOpenDialog.SetLocation(listableURI)

	test.AssertImageMatches(t, "dBFileEntry_FindFile_Tapped.png", w.Canvas().Capture())

}
