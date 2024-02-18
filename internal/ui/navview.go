package ui

import (
	"errors"
	"keepassui/internal/secretsdb"
	"keepassui/internal/secretsreader"
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type NavView struct {
	stageManager            *StageManager
	navAndListContainer     *fyne.Container
	navTop                  *fyne.Container
	SaveButton              *widget.Button
	GroupCreateButton       *widget.Button
	SecretEntryCreateButton *widget.Button
	GoBackButton            *widget.Button
	breadCrumbs             *fyne.Container
	generalButtons          *fyne.Container
	listPanel               *fyne.Container
	addEntryView            EntryUpdater
	parent                  fyne.Window
	secretsReader           secretsreader.SecretReader
	currentPath             string
}

func (n *NavView) DataChanged() {
	if n.secretsReader.GetUriID() == "" {
		return
	}

	err := n.secretsReader.ReadEntriesFromContentGroupedByPath()

	if err != nil {
		dialog.ShowError(errors.New("Error reading secrets: "+err.Error()), n.parent)
		return
	}

	n.UpdateNavView(n.secretsReader.GetFirstPath())

	n.SaveButton.OnTapped = func() {
		bytes, err := n.secretsReader.WriteDBBytes()

		if err != nil {
			dialog.ShowError(err, n.parent)
			return
		}
		fileSaveDialog := createFileSaveDialog(bytes, n.secretsReader.GetUriID(), n.parent)

		if fileSaveDialog != nil {
			fileSaveDialog.Show()
		}
	}

	if n.stageManager != nil {
		err := n.stageManager.TakeOver(n.GetStageName())
		if err != nil {
			slog.Error(err.Error())
		}
	}

}

func createFileSaveDialog(bytes []byte, originalURI string, parent fyne.Window) *dialog.FileDialog {
	fileSaveDialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {

		if err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if uc == nil {
			return
		}
		defer uc.Close()
		_, writeerr := uc.Write(bytes)

		if writeerr != nil {
			dialog.ShowError(writeerr, parent)
			return
		}

	}, parent)

	fURI, uriErr := storage.ParseURI(originalURI)
	if uriErr != nil {
		dialog.ShowError(uriErr, parent)
		return nil
	}

	locationURI, err := getLocationURI(fURI)
	if err != nil {
		dialog.ShowError(err, parent)
		return nil
	}
	if locationURI != nil {
		fileSaveDialog.SetLocation(locationURI)
		fileSaveDialog.SetFileName(fURI.Name())
	}

	fileSaveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	return fileSaveDialog
}

func getLocationURI(fURI fyne.URI) (fyne.ListableURI, error) {
	if !fyne.CurrentDevice().IsMobile() {
		listable, err := storage.CanList(fURI)

		if err != nil {
			slog.Error(err.Error())
		}
		// if full URI is not listable, attempt with parent
		if !listable {
			locationURI, err := storage.Parent(fURI)
			if err == nil {
				listable, err = storage.CanList(locationURI)
			}
			if err == nil && listable {
				listableURI, err := storage.ListerForURI(locationURI)
				if err != nil {
					return nil, err
				} else {
					return listableURI, nil
				}
			}
		}
	}
	return nil, nil
}

func (n *NavView) UpdateNavView(path string) {

	list, err := createListNav(path, n.parent, n)
	list.Refresh()

	if err != nil {
		dialog.ShowError(err, n.parent)
		return
	}

	n.breadCrumbs.RemoveAll()
	// n.breadCrumbs.Add()

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

	n.listPanel.RemoveAll()
	n.listPanel.Add(list)
	n.listPanel.Refresh()
	n.currentPath = path

	n.GroupCreateButton.OnTapped = func() {
		groupNameEntry := widget.NewEntry()
		groupNameEntry.Validator = createValidator("Group")
		form := dialog.NewForm("Add new group", "Confirm", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", groupNameEntry)}, func(valid bool) {
			if valid {
				newGroup := secretsdb.SecretEntry{Path: pathComponents, Group: path, Title: groupNameEntry.Text, IsGroup: true}
				n.secretsReader.AddSecretEntry(newGroup)
				n.UpdateNavView(path)
			}
		}, n.parent)
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
				}
			}

			label.SetText(secret.Title)
			if secret.IsGroup {
				icon.SetResource(theme.FolderIcon())
				copyPasswordPaddedContainer.Hide()
				// editButton.Hide()

				editButton.OnTapped = func() {
					groupNameEntry := widget.NewEntry()
					groupNameEntry.Text = secret.Title
					groupNameEntry.Validator = createValidator("Group")
					form := dialog.NewForm("Change group name", "Confirm", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", groupNameEntry)}, func(valid bool) {
						if valid {
							// newGroup := secretsdb.SecretEntry{Path: pathComponents, Group: path, Title: groupNameEntry.Text, IsGroup: true}
							// n.secretsReader.AddSecretEntry(newGroup)
							originalTitle := secret.Title
							secret.Title = groupNameEntry.Text
							navView.secretsReader.ModifySecretEntry(originalTitle, secret.Group, secret.IsGroup, secret)
							navView.UpdateNavView(path)
						}
					}, navView.parent)
					form.Show()
				}

				openGroupButton.OnTapped = func() {
					navView.UpdateNavView(strings.Join([]string{secret.Group, secret.Title}, "|"))
				}
			} else {
				openGroupPadddedContainer.Hide()
				icon.SetResource(theme.DocumentIcon())
				copyPasswordButton.OnTapped = func() {
					parent.Clipboard().SetContent(secret.Password)
				}
				editButton.OnTapped = func() {
					navView.addEntryView.ModifyEntry(&secret)
				}

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

func CreateNavView(secretsReader secretsreader.SecretReader, addEntryView EntryUpdater, parent fyne.Window, stageManager *StageManager) NavView {

	breadCrumbs := container.NewHBox()
	generalButtons := container.NewHBox()

	goBackButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {

	})

	breadcrumbsWithBackButton := container.NewBorder(nil, nil, breadCrumbs, goBackButton, nil)

	navTop := container.NewBorder(container.NewPadded(
		generalButtons,
	), breadcrumbsWithBackButton, nil, nil, nil)

	saveButton := widget.NewButtonWithIcon("save", theme.DocumentSaveIcon(), func() {

	})
	groupCreateButton := widget.NewButtonWithIcon("new group", theme.FolderNewIcon(), func() {

	})
	secretEntryCreateButton := widget.NewButtonWithIcon("new secret", theme.DocumentCreateIcon(), func() {

	})
	generalButtons.Add(container.NewPadded(saveButton))
	generalButtons.Add(container.NewPadded(secretEntryCreateButton))
	generalButtons.Add(container.NewPadded(groupCreateButton))

	listPanel := container.NewStack()
	navAndListContainer := container.NewBorder(container.NewVBox(navTop, widget.NewSeparator()), nil, nil, nil, listPanel)

	return NavView{
		stageManager:            stageManager,
		navAndListContainer:     navAndListContainer,
		breadCrumbs:             breadCrumbs,
		listPanel:               listPanel,
		addEntryView:            addEntryView,
		parent:                  parent,
		secretsReader:           secretsReader,
		currentPath:             "",
		GoBackButton:            goBackButton,
		generalButtons:          generalButtons,
		navTop:                  navTop,
		SaveButton:              saveButton,
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
	n.UpdateNavView(n.currentPath)
}
