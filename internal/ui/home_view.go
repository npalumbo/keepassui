package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type HomeView struct {
	DefaultStager
	dbFileEntry      *DBFileEntry
	StagerController StagerController
	homeContainer    *fyne.Container
}

func CreateHomeView(dbFileEntry *DBFileEntry, StagerController StagerController) HomeView {
	homeContainer := container.NewBorder(container.NewPadded(dbFileEntry.Container), nil, nil, nil, nil)
	return HomeView{DefaultStager: DefaultStager{}, dbFileEntry: dbFileEntry, StagerController: StagerController, homeContainer: homeContainer}
}

func (a *HomeView) GetPaintedContainer() *fyne.Container {
	return a.homeContainer
}

func (a *HomeView) GetStageName() string {
	return "Home"
}
