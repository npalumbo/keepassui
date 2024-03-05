package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HomeView struct {
	DefaultStager
	dbFileEntry      *DBFileEntry
	stagerController StagerController
	homeContainer    *fyne.Container
}

func CreateHomeView(dbFileEntry *DBFileEntry, stagerController StagerController, parent fyne.Window) HomeView {
	button := widget.NewButtonWithIcon("New KeepassDB", theme.DocumentCreateIcon(),
		func() {
			err := stagerController.TakeOver("PasswordConfirm")
			if err != nil {
				dialog.ShowError(err, parent)
			}
		},
	)
	homeContainer := container.NewBorder(container.NewVBox(container.NewPadded(dbFileEntry.Container), container.NewPadded(button)), nil, nil, nil, nil)

	return HomeView{
		DefaultStager:    DefaultStager{},
		dbFileEntry:      dbFileEntry,
		stagerController: stagerController,
		homeContainer:    homeContainer,
	}
}

func (h *HomeView) GetPaintedContainer() *fyne.Container {
	return h.homeContainer
}

func (h *HomeView) GetStageName() string {
	return "Home"
}
