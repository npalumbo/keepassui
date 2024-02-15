package ui

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dchest/uniuri"
)

type MasterPasswordDialog struct {
	DbPathAndPassword *DBPathAndPassword
	Dialog            *dialog.FormDialog
	PasswordEntry     *widget.Entry
	formItems         []*widget.FormItem
	parent            fyne.Window
	notify            binding.String
}

type DBPathAndPassword struct {
	UriID          string
	ContentInBytes []byte
	Password       string
}

func CreateDialog(dbPathAndPassword *DBPathAndPassword, parent fyne.Window) MasterPasswordDialog {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("KeyPass DB password")
	formItems = append(formItems, widget.NewFormItem("password", passwordEntry))

	return MasterPasswordDialog{
		DbPathAndPassword: dbPathAndPassword,
		Dialog:            nil,
		PasswordEntry:     passwordEntry,
		formItems:         formItems,
		parent:            parent,
		notify:            binding.NewString(),
	}

}

func (m *MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.notify.AddListener(l)
}

func (m *MasterPasswordDialog) ShowDialog(uriID string, contentInBytes *[]byte) {
	m.Dialog = dialog.NewForm("Enter master password", "Confirm", "Cancel", m.formItems, func(valid bool) {
		if valid {
			m.DbPathAndPassword.ContentInBytes = *contentInBytes
			m.DbPathAndPassword.UriID = uriID
			m.DbPathAndPassword.Password = m.PasswordEntry.Text
			m.PasswordEntry.Text = ""
			err := m.notify.Set(uniuri.New())
			if err != nil {
				slog.Error("Error notifying changes to listener", err)
			}
		}
	}, m.parent)

	m.Dialog.Resize(fyne.NewSize(400, 100))
	m.Dialog.Show()
}
