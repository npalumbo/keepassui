package ui_test

import (
	"errors"
	"image/color"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

type ColouredStager struct {
	colour    string
	container *fyne.Container
}

func CreateColouredStager(colour string, rgba color.RGBA) ColouredStager {
	return ColouredStager{colour: colour, container: container.NewStack(canvas.NewRectangle(rgba))}
}

func (c ColouredStager) GetStageName() string {
	return c.colour
}

func (c ColouredStager) GetPaintedContainer() *fyne.Container {
	return c.container
}

func (c ColouredStager) ExecuteOnTakeOver() {
	c.container.Add(canvas.NewText(c.GetStageName(), color.Black))
}

func TestStageManager(t *testing.T) {
	mainStage := container.NewStack()

	stageManager := ui.CreateStageManager(mainStage)

	w := test.NewWindow(mainStage)
	w.Resize(fyne.Size{Width: 600, Height: 600})

	redStager := CreateColouredStager("red", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	blueStager := CreateColouredStager("blue", color.RGBA{R: 0, G: 0, B: 255, A: 255})
	stageManager.RegisterStager(redStager)
	stageManager.RegisterStager(blueStager)

	test.AssertImageMatches(t, "StageManager_BeforeTakeOver.png", w.Canvas().Capture())

	err := stageManager.TakeOver("red")

	assert.Nil(t, err)

	test.AssertImageMatches(t, "StageManager_RedTakeOver.png", w.Canvas().Capture())

	err = stageManager.TakeOver("blue")
	assert.Nil(t, err)

	test.AssertImageMatches(t, "StageManager_BlueTakeOver.png", w.Canvas().Capture())

	err = stageManager.TakeOver("green")
	assert.Equal(t, errors.New("Unknown stager: green"), err, "TakeOver green should error as there is no green stager")

}
