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
	content           binding.Bytes
	dialog            *dialog.FormDialog
	passwordEntry     *widget.Entry
	formItems         []*widget.FormItem
	parent            fyne.Window
}

type DBPathAndPassword struct {
	Path     string
	Password string
}

func CreateDialog(parent fyne.Window) MasterPasswordDialog {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("KeyPass DB password")
	formItems = append(formItems, widget.NewFormItem("password", passwordEntry))
	dbPathAndPassword := binding.NewUntyped()
	content := binding.NewBytes()

	return MasterPasswordDialog{
		dbPathAndPassword: dbPathAndPassword,
		content:           content,
		dialog:            nil,
		passwordEntry:     passwordEntry,
		formItems:         formItems,
		parent:            parent,
	}

}

func (m *MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.dbPathAndPassword.AddListener(l)
}

func (m *MasterPasswordDialog) ShowDialog(path binding.URI, contentInBytes *[]byte) {
	m.dialog = dialog.NewForm("Enter master password", "Confirm", "Cancel", m.formItems, func(valid bool) {
		if valid {
			err := m.content.Set(*contentInBytes)
			if err != nil {
				slog.Error("Error updating DB bytes", err)
			}
			if m.passwordEntry.Text != "" {
				pathURI, err := path.Get()
				if err == nil {
					err = m.dbPathAndPassword.Set(DBPathAndPassword{
						Path:     pathURI.Path(),
						Password: m.passwordEntry.Text,
					})
				}
				if err != nil {
					slog.Error("Error updating Path and Password", err)

				}

			} else {
				slog.Error("You have to enter a password")
			}
		} else {
			slog.Error("invalid password")
		}
	}, m.parent)

	m.dialog.Resize(fyne.NewSize(400, 100))
	m.dialog.Show()
}
