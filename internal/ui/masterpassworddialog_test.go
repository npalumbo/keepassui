package ui_test

import (
	"keepassui/internal/secretsreader"
	"keepassui/internal/ui"
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func TestMasterPasswordDialog_Render(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	masterPasswordDialog, _ := ui.CreateDialog(w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://path", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())
}

func TestMasterPasswordDialog_fillIn_And_Submit(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes := make([]byte, 5)
	masterPasswordDialog, secretsReader := ui.CreateDialog(w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.Type(masterPasswordDialog.PasswordEntry, "thePassword")

	if secretsReader.GetUriID() != "" {
		t.Error("UriID from secretsReader should be empty string on start")
	}

	masterPasswordDialog.Dialog.Submit()

	if secretsReader.GetUriID() == "" {
		t.Error("UriID from secretsReader should not be empty string after submit")
	}

	if secretsReader.GetUriID() != "file://fakeKeypassDBFilePath" {
		t.Error("Expecting this URI: file://fakeKeypassDBFilePath")
	}
}

func TestMasterPasswordDialog_Should_Not_Show_A_Previously_Entered_Password(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes, err := os.ReadFile("testdata/files/db.kdbx")

	if err != nil {
		t.Fatal("Could not find test DB")
	}
	masterPasswordDialog, _ := ui.CreateDialog(w)
	w.Resize(fyne.NewSize(600, 600))

	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show.png", w.Canvas().Capture())

	test.Type(masterPasswordDialog.PasswordEntry, "keepassui")

	masterPasswordDialog.Dialog.Submit()

	// Second time ShowDialog is called it should not have the previous password
	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.AssertImageMatches(t, "masterPasswordDialog_Show_Second_Time.png", w.Canvas().Capture())
}

func TestMasterPasswordDialog_Calls_Listener(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	contentInBytes, err := os.ReadFile("testdata/files/db.kdbx")

	if err != nil {
		t.Fatal("Could not find test DB")
	}
	masterPasswordDialog, secretsReader := ui.CreateDialog(w)
	w.Resize(fyne.NewSize(600, 600))
	masterPasswordDialog.ShowDialog("file://fakeKeypassDBFilePath", &contentInBytes)

	test.Type(masterPasswordDialog.PasswordEntry, "keepassui")

	listener := &fakeListener{dataHasChangedToExpectedValues: false, secretsReader: secretsReader.(*secretsreader.DefaultSecretsReader)}

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
	secretsReader                  *secretsreader.DefaultSecretsReader
}

func (f *fakeListener) DataChanged() {
	data := *f.secretsReader
	if data.UriID == "file://fakeKeypassDBFilePath" && data.Password == "keepassui" {
		f.dataHasChangedToExpectedValues = true
	}
}
