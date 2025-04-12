package ui

import (
	"keepassui/internal/secretsdb"
	"keepassui/internal/secretsreader"
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type NavView struct {
	stagerController        StagerController
	navAndListContainer     *fyne.Container
	navTop                  *fyne.Container
	GroupCreateButton       *widget.Button
	SecretEntryCreateButton *widget.Button
	GoBackButton            *widget.Button
	LockDBButton            *widget.Button
	breadCrumbs             *fyne.Container
	generalButtons          *fyne.Container
	ListPanel               *fyne.Container
	addEntryView            EntryUpdater
	parent                  fyne.Window
	secretsReader           secretsreader.SecretReader
	currentPath             string
}

func (n *NavView) DataChanged() {
	if n.secretsReader.GetUriID() == "" {
		return
	}

	n.UpdateNavView(n.secretsReader.GetFirstPath())

	if n.stagerController != nil {
		err := n.stagerController.TakeOver(n.GetStageName())
		if err != nil {
			slog.Error(err.Error())
		}
	}

}

func (n *NavView) UpdateNavView(path string) {

	list, err := createListNav(path, n.parent, n)
	list.Refresh()

	if err != nil {
		dialog.ShowError(err, n.parent)
		return
	}

	n.breadCrumbs.RemoveAll()

	pathComponents := strings.Split(path, "|")
	pathAcc := []string{}
	for i, group := range pathComponents {
		pathAcc = append(pathAcc, group)
		computedPath := strings.Join(pathAcc, "|")
		pathComponentButton := widget.NewButton(group, func() {
			n.UpdateNavView(computedPath)
		})
		if i == len(pathComponents)-1 {
			pathComponentButton.Disable()
		}
		n.breadCrumbs.Add(pathComponentButton)
	}

	if len(pathComponents) < 2 {
		n.GoBackButton.Disable()
	} else {
		n.GoBackButton.Enable()
		n.GoBackButton.OnTapped = func() {
			index := strings.LastIndex(path, "|")
			backPath := path[:index]
			n.UpdateNavView(backPath)
		}
	}

	n.ListPanel.RemoveAll()
	n.ListPanel.Add(list)
	n.ListPanel.Refresh()
	n.currentPath = path

	n.GroupCreateButton.OnTapped = func() {
		groupNameEntry := widget.NewEntry()
		groupNameEntry.Validator = createValidator("Group")
		form := dialog.NewForm(lang.L("Add new group"), lang.L("Confirm"), lang.L("Cancel"), []*widget.FormItem{widget.NewFormItem(lang.L("Name"), groupNameEntry)}, func(valid bool) {
			if valid {
				newGroup := secretsdb.SecretEntry{Path: pathComponents, Group: path, Title: groupNameEntry.Text, IsGroup: true}
				n.secretsReader.AddSecretEntry(newGroup)
				n.UpdateNavView(path)
				err := n.secretsReader.Save()
				if err != nil {
					dialog.ShowError(err, n.parent)
				}
			}
		}, n.parent)
		form.SetOnClosed(func() {
			n.ExecuteOnTakeOver()
		})
		n.disableBackButton()
		form.Show()
	}

	n.SecretEntryCreateButton.OnTapped = func() {
		templateEntry := secretsdb.SecretEntry{Path: pathComponents, Group: path, IsGroup: false}
		n.addEntryView.AddEntry(&templateEntry)
	}
}

func createListNav(path string, parent fyne.Window, navView *NavView) (*widget.List, error) {
	listOfSecretsForPath := navView.secretsReader.GetEntriesForPath(path)
	untypedList := binding.NewUntypedList()
	newList := widget.NewListWithData(untypedList,
		func() fyne.CanvasObject {
			copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {})
			openGroupButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {})
			editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {})
			deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})
			buttons := container.NewHBox(
				container.NewPadded(copyButton),
				container.NewPadded(openGroupButton),
				container.NewPadded(editButton),
				container.NewPadded(deleteButton),
			)
			templateLabel := widget.NewLabel("template")
			iconAndLabel := container.NewHBox(widget.NewIcon(theme.FolderIcon()), templateLabel)
			container := container.NewBorder(nil, nil, iconAndLabel, buttons, nil)
			return container
		},
		func(lii binding.DataItem, co fyne.CanvasObject) {
			box := co.(*fyne.Container)
			di := lii.(binding.Untyped)
			untyped, err := di.Get()

			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			secret := untyped.(secretsdb.SecretEntry)
			objects := box.Objects
			iconAndLabel := objects[0].(*fyne.Container)

			icon := iconAndLabel.Objects[0].(*widget.Icon)

			label := iconAndLabel.Objects[1].(*widget.Label)

			buttons := objects[1].(*fyne.Container)
			copyPasswordPaddedContainer := buttons.Objects[0].(*fyne.Container)
			copyPasswordButton := copyPasswordPaddedContainer.Objects[0].(*widget.Button)
			openGroupPadddedContainer := buttons.Objects[1].(*fyne.Container)
			openGroupButton := openGroupPadddedContainer.Objects[0].(*widget.Button)
			editButton := buttons.Objects[2].(*fyne.Container).Objects[0].(*widget.Button)
			deleteButton := buttons.Objects[3].(*fyne.Container).Objects[0].(*widget.Button)
			deleteButton.OnTapped = func() {
				deleted := navView.secretsReader.DeleteSecretEntry(secret)
				if deleted {
					navView.UpdateNavView(secret.Group)
					err := navView.secretsReader.Save()
					if err != nil {
						dialog.ShowError(err, navView.parent)
					}
				}
			}

			label.SetText(secret.Title)
			if secret.IsGroup {
				icon.SetResource(theme.FolderIcon())
				copyPasswordPaddedContainer.Hide()

				editButton.OnTapped = func() {
					groupNameEntry := widget.NewEntry()
					groupNameEntry.Text = secret.Title
					groupNameEntry.Validator = createValidator("Group")
					form := dialog.NewForm(lang.L("Change group name"), lang.L("Confirm"), lang.L("Cancel"), []*widget.FormItem{widget.NewFormItem(lang.L("Name"), groupNameEntry)}, func(valid bool) {
						if valid {
							originalTitle := secret.Title
							secret.Title = groupNameEntry.Text
							navView.secretsReader.ModifySecretEntry(originalTitle, secret.Group, secret.IsGroup, secret)
							err := navView.secretsReader.Save()
							if err != nil {
								dialog.ShowError(err, navView.parent)
							}
							navView.UpdateNavView(path)
						}
					}, navView.parent)
					form.SetOnClosed(func() {
						navView.ExecuteOnTakeOver()
					})
					navView.disableBackButton()
					form.Show()
				}

				openGroupButton.OnTapped = func() {
					navView.UpdateNavView(strings.Join([]string{secret.Group, secret.Title}, "|"))
				}
			} else {
				openGroupPadddedContainer.Hide()
				icon.SetResource(theme.DocumentIcon())
				copyPasswordButton.OnTapped = func() {
					//nolint:staticcheck // parent.Clipboard is deprecated, but fixing requires a significant refactor
					parent.Clipboard().SetContent(secret.Password)
				}
				editButton.OnTapped = func() {
					navView.addEntryView.ModifyEntry(&secret)
				}

			}

		})

	newList.OnSelected = func(id widget.ListItemID) {
		newList.UnselectAll()
		secretRaw, err := untypedList.GetValue(id)
		if err == nil {
			secret, ok := secretRaw.(secretsdb.SecretEntry)
			if ok {
				if secret.IsGroup {
					navView.UpdateNavView(strings.Join([]string{secret.Group, secret.Title}, "|"))
				} else {
					navView.addEntryView.ModifyEntry(&secret)
				}
			}
		}
	}

	for _, v := range listOfSecretsForPath {
		err := untypedList.Append(v)
		if err != nil {
			return nil, err
		}
	}

	return newList, nil
}

func CreateNavView(secretsReader secretsreader.SecretReader, addEntryView EntryUpdater, parent fyne.Window, stagerController StagerController) NavView {

	breadCrumbs := container.NewHBox()
	generalButtons := container.NewHBox()

	goBackButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {

	})

	breadcrumbsWithBackButton := container.NewBorder(nil, nil, breadCrumbs, goBackButton, nil)

	navTop := container.NewBorder(container.NewPadded(
		generalButtons,
	), container.NewPadded(breadcrumbsWithBackButton), nil, nil, nil)

	groupCreateButton := widget.NewButtonWithIcon(lang.L("new group"), theme.FolderNewIcon(), func() {

	})
	secretEntryCreateButton := widget.NewButtonWithIcon(lang.L("new secret"), theme.DocumentCreateIcon(), func() {

	})

	lockDBButton := widget.NewButtonWithIcon(lang.L("lock db"), theme.LogoutIcon(), func() {
		goToHomeView(stagerController, parent)
	})

	generalButtons.Add(container.NewPadded(groupCreateButton))
	generalButtons.Add(container.NewPadded(secretEntryCreateButton))

	listPanel := container.NewPadded()
	navAndListContainer := container.NewBorder(container.NewVBox(navTop, widget.NewSeparator()), container.NewPadded(lockDBButton), nil, nil, listPanel)

	return NavView{
		stagerController:        stagerController,
		navAndListContainer:     navAndListContainer,
		breadCrumbs:             breadCrumbs,
		ListPanel:               listPanel,
		addEntryView:            addEntryView,
		parent:                  parent,
		secretsReader:           secretsReader,
		currentPath:             "",
		GoBackButton:            goBackButton,
		LockDBButton:            lockDBButton,
		generalButtons:          generalButtons,
		navTop:                  navTop,
		GroupCreateButton:       groupCreateButton,
		SecretEntryCreateButton: secretEntryCreateButton,
	}

}

func (n *NavView) GetPaintedContainer() *fyne.Container {
	return n.navAndListContainer
}

func (n *NavView) GetStageName() string {
	return "NavView"
}

func (n *NavView) ExecuteOnTakeOver() {
	n.parent.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == mobile.KeyBack {
			if n.GoBackButton.Disabled() {
				goToHomeView(n.stagerController, n.parent)
			} else {
				n.GoBackButton.OnTapped()
			}
		}
	})
	n.UpdateNavView(n.currentPath)
}

func (n *NavView) disableBackButton() {
	n.parent.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
	})
}
