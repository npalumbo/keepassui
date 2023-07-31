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
}

func CreateDialog(path binding.String) MasterPasswordDialog {
	return MasterPasswordDialog{
		path:              path,
		dbPathAndPassword: binding.NewUntyped(),
	}
}

func (m MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.dbPathAndPassword.AddListener(l)
}

type Data struct {
	Path     string
	Password string
}

func (m MasterPasswordDialog) ShowDialog(parent fyne.Window) {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("KeyPass DB password")
	formItems = append(formItems, widget.NewFormItem("password", passwordEntry))
	form := dialog.NewForm("Enter master password", "Confirm", "Cancel", formItems, func(valid bool) {
		if valid {
			if passwordEntry.Text != "" {
				path, _ := m.path.Get()
				err := m.dbPathAndPassword.Set(Data{
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
	form.Resize(fyne.NewSize(400, 100))
	form.Show()
}
