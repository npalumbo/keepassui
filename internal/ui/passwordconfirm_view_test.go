package ui_test

import (
	mocks_ui "keepassui/internal/mocks/ui"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePasswordConfirmView(t *testing.T) {

	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})

	test.AssertImageMatches(t, "CreatePasswordConfirmView.png", w.Canvas().Capture())
}

func TestPasswordConfirmViewTapOnCancelCallsStagerControllerHome(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, stagerController, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	stagerController.EXPECT().TakeOver("Home").Times(1)

	buttonCancel := findButton(objects, "Cancel")

	test.Tap(buttonCancel)
}

func TestPasswordConfirmViewValidationBothFieldsEmpty(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	form := findForm(objects)

	assert.Equal(t, "pass cannot be empty", form.Validate().Error())
}

func TestPasswordConfirmViewValidationOnlyConfirmPasswordFieldEmpty(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	form := findForm(objects)

	passwordEntry := form.Items[0].Widget.(*widget.Entry)
	test.Type(passwordEntry, "aPassword")
	assert.Equal(t, "repeated password can't be empty", form.Validate().Error())
}

func TestPasswordConfirmViewValidationOnlyPasswordFieldEmpty(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	form := findForm(objects)

	confirmPasswordEntry := form.Items[1].Widget.(*widget.Entry)
	test.Type(confirmPasswordEntry, "aPassword")
	assert.Equal(t, "pass cannot be empty", form.Validate().Error())
}

func TestPasswordConfirmViewValidationPasswordsNotMatching(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	form := findForm(objects)

	passwordEntry := form.Items[0].Widget.(*widget.Entry)
	test.Type(passwordEntry, "aPassword")
	confirmPasswordEntry := form.Items[1].Widget.(*widget.Entry)
	test.Type(confirmPasswordEntry, "aConfirmPassword")
	assert.Equal(t, "introduce the same password in both fields", form.Validate().Error())
}

func TestPasswordConfirmViewValidationPasswordsMatching(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(nil, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	form := findForm(objects)

	passwordEntry := form.Items[0].Widget.(*widget.Entry)
	confirmPasswordEntry := form.Items[1].Widget.(*widget.Entry)

	test.Type(passwordEntry, "aPassword")
	test.Type(confirmPasswordEntry, "aPassword")

	assert.Nil(t, form.Validate())
}

func TestPasswordConfirmViewTapOnSubmitCallsFileSaver(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fileSaver := mocks_ui.NewMockFileSaver(mockCtrl)

	w := test.NewWindow(container.NewWithoutLayout())
	passwordConfirmView := ui.CreatePasswordConfirmView(fileSaver, nil, w)

	w.SetContent(passwordConfirmView.GetPaintedContainer())

	objects := test.LaidOutObjects(passwordConfirmView.GetPaintedContainer())

	form := findForm(objects)

	passwordEntry := form.Items[0].Widget.(*widget.Entry)
	passwordEntry.Text = "aPassword"
	confirmPasswordEntry := form.Items[1].Widget.(*widget.Entry)
	confirmPasswordEntry.Text = "aPassword"

	fileSaver.EXPECT().ShowForMasterPassword("aPassword").Times(1)
	form.OnSubmit()
}

func findForm(objects []fyne.CanvasObject) *widget.Form {
	var form *widget.Form

	for _, v := range objects {
		formFound, ok := v.(*widget.Form)
		if ok {
			form = formFound
		}
	}
	return form
}

func findButton(objects []fyne.CanvasObject, matchingText string) *widget.Button {
	var button *widget.Button

	for _, v := range objects {
		buttonFound, ok := v.(*widget.Button)
		if ok && buttonFound.Text == matchingText {
			button = buttonFound
		}
	}
	return button
}
