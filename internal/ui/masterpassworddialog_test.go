package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
)

func TestMasterPasswordDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	path := binding.NewString()
	masterPasswordDialog := CreateDialog(path, w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog()

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())

}
