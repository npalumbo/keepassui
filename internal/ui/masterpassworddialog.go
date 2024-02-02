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
	dbPathAndPassword *DBPathAndPassword
	dialog            *dialog.FormDialog
	passwordEntry     *widget.Entry
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
		dbPathAndPassword: dbPathAndPassword,
		dialog:            nil,
		passwordEntry:     passwordEntry,
		formItems:         formItems,
		parent:            parent,
		notify:            binding.NewString(),
	}

}

func (m *MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.notify.AddListener(l)
}

func (m *MasterPasswordDialog) ShowDialog(uriID string, contentInBytes *[]byte) {
	m.dialog = dialog.NewForm("Enter master password", "Confirm", "Cancel", m.formItems, func(valid bool) {
		if valid {
			m.dbPathAndPassword.ContentInBytes = *contentInBytes
			m.dbPathAndPassword.UriID = uriID
			m.dbPathAndPassword.Password = m.passwordEntry.Text
			err := m.notify.Set(uniuri.New())
			if err != nil {
				slog.Error("Error notifying changes to listener", err)
			}
		}
	}, m.parent)

	m.dialog.Resize(fyne.NewSize(400, 100))
	m.dialog.Show()
}
