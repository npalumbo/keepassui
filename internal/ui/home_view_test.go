package ui_test

import (
	mocks_ui "keepassui/internal/mocks/ui"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/mock/gomock"
)

func TestCreateHomeView(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())

	dbFileEntry := ui.CreateDBFileEntry(nil, w)
	homeView := ui.CreateHomeView(&dbFileEntry, nil, w)

	w.SetContent(homeView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})

	test.AssertImageMatches(t, "CreateHomeView.png", w.Canvas().Capture())
}

func TestHomeViewTapOnNewDBCallsTakeOverPasswordConfirm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())

	dbFileEntry := ui.CreateDBFileEntry(nil, w)
	homeView := ui.CreateHomeView(&dbFileEntry, stagerController, w)

	w.SetContent(homeView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})

	stagerController.EXPECT().TakeOver("PasswordConfirm").Times(1)

	objects := test.LaidOutObjects(homeView.GetPaintedContainer())

	for _, v := range objects {
		button, ok := v.(*widget.Button)
		if ok && button.Text == "New KeepassDB" {
			test.Tap(button)
		}
	}
}
