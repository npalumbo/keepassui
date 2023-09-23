package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
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
					log.Printf("Error updating Path and Password %v", err)
				}
			}
		} else {
			log.Println("Didn't get password")
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
