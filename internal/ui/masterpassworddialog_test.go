package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestMasterPasswordDialog_Render(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	masterPasswordDialog := CreateDialog(&DBPathAndPassword{}, w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://path", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())
}

func TestMasterPasswordDialog_fillIn_And_Submit(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	dbPathAndPassword := &DBPathAndPassword{}
	masterPasswordDialog := CreateDialog(dbPathAndPassword, w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.Type(masterPasswordDialog.passwordEntry, "thePassword")

	if masterPasswordDialog.dbPathAndPassword.UriID != "" {
		t.Error("UriID from DBPathAndPassword should be empty string on start")
	}

	masterPasswordDialog.dialog.Submit()

	if masterPasswordDialog.dbPathAndPassword.UriID == "" {
		t.Error("UriID from DBPathAndPassword should not be empty string after submit")
	}

	data := *masterPasswordDialog.dbPathAndPassword
	if data.UriID == "" || data.UriID != "file://fakeKeypassDBFilePath" {
		t.Error("Expecting this URI: file://fakeKeypassDBFilePath")
	}
	if data.Password == "" || data.Password != "thePassword" {
		t.Error("Expecting password to be thePassword")
	}
}

func TestMasterPasswordDialog_Calls_Listener(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	dbPathAndPassword := &DBPathAndPassword{}
	masterPasswordDialog := CreateDialog(dbPathAndPassword, w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

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
	dbPathAndPassword              *DBPathAndPassword
}

func (f *fakeListener) DataChanged() {
	data := *f.dbPathAndPassword
	if data.UriID == "file://fakeKeypassDBFilePath" && data.Password == "thePassword" {
		f.dataHasChangedToExpectedValues = true
	}
}
