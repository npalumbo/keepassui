package ui

import (
	"errors"
	"keepassui/internal/keepass"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyAccordion struct {
	dbPathAndPassword binding.Untyped
	accordionWidget   *widget.Accordion
	detailedView      *DetailedView
	parent            fyne.Window
	createReader      func(dbPathAndPassword DBPathAndPassword) keepass.SecretReader
}

func (k *KeyAccordion) DataChanged() {
	o, err := k.dbPathAndPassword.Get()
	if err != nil {
		dialog.ShowError(err, k.parent)
		return
	}
	if o == nil {
		return
	}

	d, ok := o.(DBPathAndPassword)
	if !ok {
		dialog.ShowError(errors.New("Could not cast dbPathAndPassword to DBPathAndPassword"), k.parent)
		return
	}

	secretReader := k.createReader(d)

	secretsDB, err := secretReader.ReadEntriesFromContentGroupedByPath()

	if err != nil {
		dialog.ShowError(errors.New("Error reading secrets: "+err.Error()), k.parent)
		return
	}

	for _, path := range secretsDB.PathsInOrder {
		listOfSecretsForPath := secretsDB.EntriesByPath[path]
		newList, err := createList(listOfSecretsForPath, k.detailedView, k.parent)

		if err != nil {
			dialog.ShowError(err, k.parent)
			return
		}
		accordionItem := widget.NewAccordionItem(path, newList)
		k.accordionWidget.Append(accordionItem)
	}
}

func createList(listOfSecretsForPath []keepass.SecretEntry, detailedView *DetailedView, parent fyne.Window) (*widget.List, error) {
	untypedList := binding.NewUntypedList()
	newList := widget.NewListWithData(untypedList,
		func() fyne.CanvasObject {
			copyButton := widget.NewButtonWithIcon("copy", theme.ContentCopyIcon(), func() {})
			showInfoButton := widget.NewButtonWithIcon("details", theme.InfoIcon(), func() {})
			buttons := container.NewHBox(copyButton, showInfoButton)

			return container.NewBorder(nil, nil, widget.NewLabel(""), buttons, nil)
		},
		func(lii binding.DataItem, co fyne.CanvasObject) {
			box := co.(*fyne.Container)
			di := lii.(binding.Untyped)
			untyped, err := di.Get()

			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			secret := untyped.(keepass.SecretEntry)
			objects := box.Objects
			label := objects[0].(*widget.Label)
			label.SetText(secret.Title)
			buttons := objects[1].(*fyne.Container)
			copyPasswordButton := buttons.Objects[0].(*widget.Button)
			copyPasswordButton.OnTapped = func() {
				parent.Clipboard().SetContent(secret.Password)
			}
			showInfoButton := buttons.Objects[1].(*widget.Button)
			showInfoButton.OnTapped = func() {
				detailedView.UpdateDetails(secret)
			}
		})

	for _, v := range listOfSecretsForPath {
		err := untypedList.Append(v)
		if err != nil {
			return nil, err
		}
	}

	return newList, nil
}

func CreatekeyAccordion(dbPathAndPassword binding.Untyped, detailedView *DetailedView, parent fyne.Window, createReader ToSecretReaderFn) KeyAccordion {
	return KeyAccordion{
		dbPathAndPassword: dbPathAndPassword,
		accordionWidget:   widget.NewAccordion(),
		detailedView:      detailedView,
		parent:            parent,
		createReader:      createReader,
	}
}
