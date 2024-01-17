package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
)

func TestMasterPasswordDialog_Render(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	path := binding.NewString()
	contentInBytes := make([]byte, 5)
	masterPasswordDialog := CreateDialog(path, &contentInBytes, w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog()

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())

}

func TestMasterPasswordDialog_fillIn_And_Submit(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	path := binding.NewString()
	err := path.Set("fakeKeypassDBFilePath")
	if err != nil {
		t.Error()
	}
	contentInBytes := make([]byte, 5)
	masterPasswordDialog := CreateDialog(path, &contentInBytes, w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog()

	test.Type(masterPasswordDialog.passwordEntry, "thePassword")

	rawData, _ := masterPasswordDialog.dbPathAndPassword.Get()
	if rawData != nil {
		t.Error()
	}

	masterPasswordDialog.dialog.Submit()

	rawData, _ = masterPasswordDialog.dbPathAndPassword.Get()
	if rawData == nil {
		t.Error()
	}
	data := rawData.(DBPathAndPassword)
	if data.Path == "" || data.Path != "fakeKeypassDBFilePath" {
		t.Error()
	}
	if data.Password == "" || data.Password != "thePassword" {
		t.Error()
	}
}

func TestMasterPasswordDialog_Calls_Listener(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	path := binding.NewString()
	err := path.Set("fakeKeypassDBFilePath")
	if err != nil {
		t.Error()
	}
	contentInBytes := make([]byte, 5)
	masterPasswordDialog := CreateDialog(path, &contentInBytes, w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog()

	test.Type(masterPasswordDialog.passwordEntry, "thePassword")

	listener := &fakeListener{dataHasChangedToExpectedValues: false, dbPathAndPassword: masterPasswordDialog.dbPathAndPassword}

	masterPasswordDialog.AddListener(listener)

	if listener.dataHasChangedToExpectedValues == true {
		t.Error()
	}

	masterPasswordDialog.dialog.Submit()

	time.Sleep(10 * time.Millisecond)

	if listener.dataHasChangedToExpectedValues == false {
		t.Error()
	}

}

type fakeListener struct {
	dataHasChangedToExpectedValues bool
	dbPathAndPassword              binding.Untyped
}

func (f *fakeListener) DataChanged() {
	rawData, _ := f.dbPathAndPassword.Get()
	data := rawData.(DBPathAndPassword)
	if data.Path == "fakeKeypassDBFilePath" && data.Password == "thePassword" {
		f.dataHasChangedToExpectedValues = true
	}

}
