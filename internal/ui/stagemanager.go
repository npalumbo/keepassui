package ui

import "fyne.io/fyne/v2"

type StageManager struct {
	currentViewContainer *fyne.Container
	stagerMap            map[string]Stager
}

//go:generate mockgen -destination=../mocks/ui/mock_stagemanager.go -source=./stagemanager.go

type StagerController interface {
	TakeOver(name string)
}

type DefaultStager struct {
}

type Stager interface {
	GetPaintedContainer() *fyne.Container
	ExecuteOnResume()
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

func (s StageManager) TakeOver(name string) {
	s.currentViewContainer.RemoveAll()
	stager := s.stagerMap[name]
	container := stager.GetPaintedContainer()
	container.Refresh()
	s.currentViewContainer.Add(container)
	stager.ExecuteOnResume()
	s.currentViewContainer.Refresh()
}

func (d *DefaultStager) ExecuteOnResume() {

}
