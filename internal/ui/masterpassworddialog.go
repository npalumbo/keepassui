package ui

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type MasterPasswordDialog struct {
	dbPathAndPassword binding.Untyped
	dialog            *dialog.FormDialog
	passwordEntry     *widget.Entry
	formItems         []*widget.FormItem
	parent            fyne.Window
}

type DBPathAndPassword struct {
	UriID          string
	ContentInBytes []byte
	Password       string
}

func CreateDialog(parent fyne.Window) MasterPasswordDialog {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("KeyPass DB password")
	formItems = append(formItems, widget.NewFormItem("password", passwordEntry))
	dbPathAndPassword := binding.NewUntyped()

	return MasterPasswordDialog{
		dbPathAndPassword: dbPathAndPassword,
		dialog:            nil,
		passwordEntry:     passwordEntry,
		formItems:         formItems,
		parent:            parent,
	}

}

func (m *MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.dbPathAndPassword.AddListener(l)
}

func (m *MasterPasswordDialog) ShowDialog(uriID string, contentInBytes *[]byte) {
	m.dialog = dialog.NewForm("Enter master password", "Confirm", "Cancel", m.formItems, func(valid bool) {
		if valid {
			err := m.dbPathAndPassword.Set(DBPathAndPassword{
				UriID:          uriID,
				Password:       m.passwordEntry.Text,
				ContentInBytes: *contentInBytes,
			})
			if err != nil {
				slog.Error("Error updating Path and Password", err)
			}
		}
	}, m.parent)

	m.dialog.Resize(fyne.NewSize(400, 100))
	m.dialog.Show()
}
