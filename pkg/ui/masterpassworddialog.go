package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log/slog"
)

type MasterPasswordDialog struct {
	dbPathAndPassword binding.Untyped
	path              binding.String
	dialog            dialog.Dialog
}

func CreateDialog(path binding.String, parent fyne.Window) MasterPasswordDialog {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("KeyPass DB password")
	formItems = append(formItems, widget.NewFormItem("password", passwordEntry))
	dbPathAndPassword := binding.NewUntyped()

	dialog := dialog.NewForm("Enter master password", "Confirm", "Cancel", formItems, func(valid bool) {
		if valid {
			if passwordEntry.Text != "" {
				path, _ := path.Get()
				err := dbPathAndPassword.Set(Data{
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
		dialog:            dialog,
	}

}

func (m MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.dbPathAndPassword.AddListener(l)
}

type Data struct {
	Path     string
	Password string
}

func (m MasterPasswordDialog) ShowDialog() {
	m.dialog.Show()
}
