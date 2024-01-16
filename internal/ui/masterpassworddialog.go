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
	path              binding.String
	dialog            *dialog.FormDialog
	passwordEntry     *widget.Entry
}

type DBPathAndPassword struct {
	Path     string
	Password string
}

func CreateDialog(path binding.String, contentInBytes *[]byte, parent fyne.Window) MasterPasswordDialog {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("KeyPass DB password")
	formItems = append(formItems, widget.NewFormItem("password", passwordEntry))
	dbPathAndPassword := binding.NewUntyped()
	content := binding.NewBytes()

	dialog := dialog.NewForm("Enter master password", "Confirm", "Cancel", formItems, func(valid bool) {
		if valid {
			err := content.Set(*contentInBytes)
			if err != nil {
				slog.Error("Error updating DB bytes", err)
			}
			if passwordEntry.Text != "" {
				path, _ := path.Get()
				err := dbPathAndPassword.Set(DBPathAndPassword{
					Path:     path,
					Password: passwordEntry.Text,
				})
				if err != nil {
					slog.Error("Error updating Path and Password", err)

				}

			} else {
				slog.Error("You have to enter a password")
			}
		} else {
			slog.Error("invalid password")
		}
	}, parent)
	dialog.Resize(fyne.NewSize(400, 100))

	return MasterPasswordDialog{
		path:              path,
		dbPathAndPassword: dbPathAndPassword,
		content:           content,
		dialog:            dialog,
		passwordEntry:     passwordEntry,
	}

}

func (m MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.dbPathAndPassword.AddListener(l)
}

func (m MasterPasswordDialog) ShowDialog() {
	m.dialog.Show()
}
