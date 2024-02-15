package ui_test

import (
	"keepassui/internal/ui"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestMasterPasswordDialog_Render(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	masterPasswordDialog := ui.CreateDialog(&ui.DBPathAndPassword{}, w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://path", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())
}

func TestMasterPasswordDialog_fillIn_And_Submit(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	dbPathAndPassword := &ui.DBPathAndPassword{}
	masterPasswordDialog := ui.CreateDialog(dbPathAndPassword, w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.Type(masterPasswordDialog.PasswordEntry, "thePassword")

	if masterPasswordDialog.DbPathAndPassword.UriID != "" {
		t.Error("UriID from DBPathAndPassword should be empty string on start")
	}

	masterPasswordDialog.Dialog.Submit()

	if masterPasswordDialog.DbPathAndPassword.UriID == "" {
		t.Error("UriID from DBPathAndPassword should not be empty string after submit")
	}

	data := *masterPasswordDialog.DbPathAndPassword
	if data.UriID == "" || data.UriID != "file://fakeKeypassDBFilePath" {
		t.Error("Expecting this URI: file://fakeKeypassDBFilePath")
	}
	if data.Password == "" || data.Password != "thePassword" {
		t.Error("Expecting password to be thePassword")
	}
}

func TestMasterPasswordDialog_Should_Not_Show_A_Previously_Entered_Password(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	dbPathAndPassword := &ui.DBPathAndPassword{}
	masterPasswordDialog := ui.CreateDialog(dbPathAndPassword, w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())

	test.Type(masterPasswordDialog.PasswordEntry, "thePassword")

	masterPasswordDialog.Dialog.Submit()

	// Second time ShowDialog is called it should not have the previous password
	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show_Second_Time.png", w.Canvas().Capture())
}

func TestMasterPasswordDialog_Calls_Listener(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	dbPathAndPassword := &ui.DBPathAndPassword{}
	masterPasswordDialog := ui.CreateDialog(dbPathAndPassword, w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.Type(masterPasswordDialog.PasswordEntry, "thePassword")

	listener := &fakeListener{dataHasChangedToExpectedValues: false, dbPathAndPassword: masterPasswordDialog.DbPathAndPassword}

	masterPasswordDialog.AddListener(listener)

	if listener.dataHasChangedToExpectedValues == true {
		t.Error()
	}

	masterPasswordDialog.Dialog.Submit()

	time.Sleep(10 * time.Millisecond)

	if listener.dataHasChangedToExpectedValues == false {
		t.Error()
	}
}

type fakeListener struct {
	dataHasChangedToExpectedValues bool
	dbPathAndPassword              *ui.DBPathAndPassword
}

func (f *fakeListener) DataChanged() {
	data := *f.dbPathAndPassword
	if data.UriID == "file://fakeKeypassDBFilePath" && data.Password == "thePassword" {
		f.dataHasChangedToExpectedValues = true
	}
}
