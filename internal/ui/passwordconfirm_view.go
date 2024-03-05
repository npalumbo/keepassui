package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PasswordConfirmView struct {
	fileSaver        FileSaver
	container        *fyne.Container
	stagerController StagerController
	parent           fyne.Window
}

func CreatePasswordConfirmView(fileSaver FileSaver, stagerController StagerController, parent fyne.Window) PasswordConfirmView {
	container := createContainer(fileSaver, stagerController, parent)
	return PasswordConfirmView{fileSaver: fileSaver, container: container, stagerController: stagerController, parent: parent}

}

func createContainer(fileSaver FileSaver, stagerController StagerController, parent fyne.Window) *fyne.Container {
	firstPass := widget.NewPasswordEntry()
	firstPass.Validator = createValidator("pass")
	firstPassItem := widget.NewFormItem("password", firstPass)
	secondPass := widget.NewPasswordEntry()
	secondPass.Validator = func(s string) error {
		if secondPass.Text == "" {
			return errors.New("repeated password can't be empty")
		}
		if firstPass.Text != secondPass.Text {
			return errors.New("introduce the same password in both fields")
		}
		return nil
	}
	secondPassItem := widget.NewFormItem("confirm password", secondPass)

	form := widget.NewForm(firstPassItem, secondPassItem)

	form.OnSubmit = func() {
		fileSaver.ShowForMasterPassword(firstPass.Text)
	}
	form.OnCancel = func() {
		err := stagerController.TakeOver("Home")
		if err != nil {
			dialog.ShowError(err, parent)
		}
	}

	container := container.NewBorder(widget.NewLabel("Set up master password"), nil, nil, nil, form)
	return container
}

func (p *PasswordConfirmView) GetPaintedContainer() *fyne.Container {
	return p.container
}

func (p *PasswordConfirmView) GetStageName() string {
	return "PasswordConfirm"
}

func (p *PasswordConfirmView) ExecuteOnTakeOver() {
	p.container = createContainer(p.fileSaver, p.stagerController, p.parent)
}
