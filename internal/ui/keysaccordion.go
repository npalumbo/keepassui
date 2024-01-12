package ui

import (
	"errors"
	"keepassui/internal/keepass"
	"log/slog"
	"reflect"

	keepassuiwidget "keepassui/internal/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyAccordion struct {
	dbPathAndPassword binding.Untyped
	content           binding.Bytes
	accordionWidget   *widget.Accordion
	details           *widget.Form
	parent            fyne.Window
	cont              *fyne.Container
}

func (k *KeyAccordion) DataChanged() {
	o, err := k.dbPathAndPassword.Get()
	if err == nil && o != nil {
		d, ok := o.(Data)
		contentBytes, errBytes := k.content.Get()
		if ok {
			if errBytes != nil {
				dialog.ShowError(errors.New("Error reading secrets: "+errBytes.Error()), k.parent)
			} else if d.Password != "" {
				secrets, err := keepass.ReadEntriesFromContent(contentBytes, d.Password)

				if err != nil {
					slog.Error("Error reading secret entries", err)
					dialog.ShowError(errors.New("Error reading secrets: "+err.Error()), k.parent)
				} else {

					groupedPath := make(map[string][]keepass.SecretEntry)
					for _, p := range secrets {
						groupedPath[p.Path] = append(groupedPath[p.Path], p)
						slog.Info("path ", "path", p.Path)
					}
					for _, key := range reflect.ValueOf(groupedPath).MapKeys() {

						listOfSecrets := groupedPath[key.String()]

						untypedList := binding.NewUntypedList()
						newList := widget.NewListWithData(untypedList,
							func() fyne.CanvasObject {

								copyButton := widget.NewButtonWithIcon("copy", theme.ContentCopyIcon(), func() {})

								showInfoButton := widget.NewButtonWithIcon("details", theme.InfoIcon(), func() {})

								buttons := container.NewHBox(copyButton, showInfoButton)

								box := container.NewBorder(nil, nil, widget.NewLabel(""), buttons, nil)

								return box
							},
							func(lii binding.DataItem, co fyne.CanvasObject) {
								box := co.(*fyne.Container)
								di := lii.(binding.Untyped)
								untyped, err := di.Get()

								if err != nil {
									dialog.ShowError(err, k.parent)
									return
								} else {
									secret := untyped.(keepass.SecretEntry)
									objects := box.Objects
									label := objects[0].(*widget.Label)
									label.SetText(secret.Title)
									buttons := objects[1].(*fyne.Container)
									copyPasswordButton := buttons.Objects[0].(*widget.Button)
									copyPasswordButton.OnTapped = func() {
										k.parent.Clipboard().SetContent(secret.Password)
									}
									showInfoButton := buttons.Objects[1].(*widget.Button)
									showInfoButton.OnTapped = func() {
										labelTitle := k.details.Items[0].Widget.(*widget.Label)
										labelUsernameValue := k.details.Items[1].Widget.(*widget.Label)
										labelPasswordValue := k.details.Items[2].Widget.(*widget.Entry)
										labelUrl := k.details.Items[3].Widget.(*widget.Label)
										labelNotes := k.details.Items[4].Widget.(*widget.Label)
										labelTitle.SetText(secret.Title)
										labelUsernameValue.SetText(secret.Username)
										labelPasswordValue.SetText(secret.Password)
										labelPasswordValue.Password = true
										labelUrl.SetText(secret.Url)
										labelNotes.SetText(secret.Notes)
										k.details.Refresh()
										k.cont.Show()
									}
								}

							})

						for _, v := range listOfSecrets {
							err := untypedList.Append(v)
							if err != nil {
								slog.Error("err", "", err)
								return
							}
						}

						accordionItem := widget.NewAccordionItem(key.String(), newList)

						k.accordionWidget.Append(accordionItem)

					}

				}
			}
		}
	}
}

func CreatekeyAccordion(dbPathAndPassword binding.Untyped, content binding.Bytes, parent fyne.Window) KeyAccordion {
	accordionWidget := widget.NewAccordion()

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.ActionItem = keepassuiwidget.NewPasswordRevealerNotDisabled(passwordEntry)
	passwordEntry.Disable()

	title := widget.NewFormItem("Title", widget.NewLabel(""))
	username := widget.NewFormItem("Username", widget.NewLabel(""))
	password := widget.NewFormItem("Password", passwordEntry)
	url := widget.NewFormItem("Url", widget.NewLabel(""))
	notes := widget.NewFormItem("Notes", widget.NewLabel(""))

	details := widget.NewForm(title, username, password, url, notes)

	closeDetails := widget.NewButtonWithIcon("Close", theme.CancelIcon(), func() {

	})
	closeDetailsContainer := container.NewVBox(closeDetails, details)
	closeDetails.OnTapped = func() {
		closeDetailsContainer.Hide()
	}

	closeDetailsContainer.Hide()

	return KeyAccordion{
		dbPathAndPassword: dbPathAndPassword,
		content:           content,
		accordionWidget:   accordionWidget,
		details:           details,
		cont:              closeDetailsContainer,
		parent:            parent,
	}
}
