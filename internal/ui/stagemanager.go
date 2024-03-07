package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type StageManager struct {
	currentViewContainer *fyne.Container
	stagerMap            map[string]Stager
}

//go:generate mockgen -destination=../mocks/ui/mock_stagemanager.go -source=./stagemanager.go

type StagerController interface {
	TakeOver(name string) error
	RegisterStager(stager Stager)
	GetContainer() *fyne.Container
}

type DefaultStager struct {
}

type Stager interface {
	GetPaintedContainer() *fyne.Container
	ExecuteOnTakeOver()
	GetStageName() string
}

func CreateStageManager(currentViewContainer *fyne.Container) StageManager {
	return StageManager{
		currentViewContainer: currentViewContainer,
		stagerMap:            make(map[string]Stager),
	}
}

func (s StageManager) RegisterStager(stager Stager) {
	s.stagerMap[stager.GetStageName()] = stager
}

func (s StageManager) TakeOver(name string) error {
	stager, ok := s.stagerMap[name]

	if !ok {
		return errors.New("Unknown stager: " + name)
	}

	s.currentViewContainer.RemoveAll()
	container := stager.GetPaintedContainer()
	container.Refresh()
	s.currentViewContainer.Add(container)
	stager.ExecuteOnTakeOver()
	s.currentViewContainer.Refresh()

	return nil
}

func (s StageManager) GetContainer() *fyne.Container {
	return s.currentViewContainer
}

func (d *DefaultStager) ExecuteOnTakeOver() {

}

func handleErrorAndGoToHomeView(err error, parent fyne.Window, stagerController StagerController) {
	dialog.ShowError(err, parent)
	goToHomeView(stagerController, parent)
}

func goToHomeView(stagerController StagerController, parent fyne.Window) {
	err := stagerController.TakeOver("Home")
	if err != nil {
		dialog.ShowError(err, parent)
	}
}
