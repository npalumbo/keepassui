package ui

import (
	"fyne.io/fyne/v2"
	"keepassui/internal/keepass"
)

type DetailedView struct {
	DefaultStager
	secretForm   *SecretForm
	stageManager StageManager
}

/*func (d *DetailedView) ExecuteOnResume() {
	//TODO implement me
	panic("implement me")
}*/

func (d *DetailedView) GetPaintedContainer() *fyne.Container {
	return d.secretForm.FormContainer
}

func (d *DetailedView) GetStageName() string {
	return "DetailedView"
}

func (d *DetailedView) ShowDetails(secretEntry keepass.SecretEntry) {
	d.secretForm.UpdateForm(secretEntry)
	d.stageManager.TakeOver("DetailedView")
}

func CreateDetailedView(stageManager StageManager) DetailedView {
	secretForm := CreateForm(true)
	secretForm.detailsForm.OnSubmit = func() {
		stageManager.TakeOver("NavView")
	}
	secretForm.detailsForm.SubmitText = "Back"
	secretForm.detailsForm.Refresh()

	return DetailedView{
		stageManager: stageManager,
		secretForm:   &secretForm,
	}
}
