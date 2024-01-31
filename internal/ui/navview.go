package ui

import (
	"errors"
	"keepassui/internal/keepass"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type NavView struct {
	fullContainer     *fyne.Container
	breadCrumbs       *fyne.Container
	listPanel         *fyne.Container
	detailedView      *DetailedView
	parent            fyne.Window
	dbPathAndPassword binding.Untyped
	createReader      ToSecretReaderFn
	currentPath       string
	secretsDB         keepass.SecretsDB
}

func (n *NavView) DataChanged() {
	o, err := n.dbPathAndPassword.Get()
	if err != nil {
		dialog.ShowError(err, n.parent)
		return
	}
	if o == nil {
		return
	}

	d, ok := o.(DBPathAndPassword)
	if !ok {
		dialog.ShowError(errors.New("Could not cast dbPathAndPassword to DBPathAndPassword"), n.parent)
		return
	}

	secretReader := n.createReader(d)

	secretsDB, err := secretReader.ReadEntriesFromContentGroupedByPath()

	if err != nil {
		dialog.ShowError(errors.New("Error reading secrets: "+err.Error()), n.parent)
		return
	}

	n.secretsDB = secretsDB

	n.UpdateNavView(secretsDB.PathsInOrder[0])
}

func (n *NavView) UpdateNavView(path string) {
	listOfSecretsForPath := n.secretsDB.EntriesByPath[path]

	list, err := createListNav(listOfSecretsForPath, n.detailedView, n.parent, n)
	list.Refresh()

	if err != nil {
		dialog.ShowError(err, n.parent)
		return
	}

	n.breadCrumbs.RemoveAll()
	n.breadCrumbs.Add(widget.NewLabel("Path: "))

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

	n.listPanel.RemoveAll()
	n.listPanel.Add(list)
	n.listPanel.Refresh()
	n.currentPath = path
}

func createListNav(listOfSecretsForPath []keepass.SecretEntry, detailedView *DetailedView, parent fyne.Window, navView *NavView) (*widget.List, error) {
	untypedList := binding.NewUntypedList()
	newList := widget.NewListWithData(untypedList,
		func() fyne.CanvasObject {
			copyButton := widget.NewButtonWithIcon("copy", theme.ContentCopyIcon(), func() {})
			showInfoButton := widget.NewButtonWithIcon("details", theme.InfoIcon(), func() {})
			openGroupButton := widget.NewButtonWithIcon("open", theme.FolderOpenIcon(), func() {})
			buttons := container.NewHBox(copyButton, showInfoButton, openGroupButton)
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

			secret := untyped.(keepass.SecretEntry)
			objects := box.Objects
			iconAndLabel := objects[0].(*fyne.Container)

			icon := iconAndLabel.Objects[0].(*widget.Icon)

			label := iconAndLabel.Objects[1].(*widget.Label)

			buttons := objects[1].(*fyne.Container)
			copyPasswordButton := buttons.Objects[0].(*widget.Button)
			showInfoButton := buttons.Objects[1].(*widget.Button)
			openGroupButton := buttons.Objects[2].(*widget.Button)
			if secret.IsGroup {
				icon.SetResource(theme.FolderIcon())
				label.SetText(secret.Title)
				copyPasswordButton.Hide()
				showInfoButton.Hide()
				openGroupButton.OnTapped = func() {
					navView.UpdateNavView(strings.Join([]string{secret.Group, secret.Title}, "|"))
				}
			} else {
				openGroupButton.Hide()
				icon.SetResource(theme.FileTextIcon())
				label.SetText(secret.Title)
				copyPasswordButton.OnTapped = func() {
					parent.Clipboard().SetContent(secret.Password)
				}
				showInfoButton.OnTapped = func() {
					detailedView.UpdateDetails(secret)
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

func CreateNavView(dbPathAndPassword binding.Untyped, detailedView *DetailedView, parent fyne.Window, createReader ToSecretReaderFn) NavView {

	breadCrumbs := container.NewHBox()
	listPanel := container.NewStack()
	fullContainer := container.NewBorder(container.NewVBox(breadCrumbs, widget.NewSeparator()), nil, nil, nil, listPanel)

	return NavView{
		fullContainer:     fullContainer,
		breadCrumbs:       breadCrumbs,
		listPanel:         listPanel,
		detailedView:      detailedView,
		parent:            parent,
		dbPathAndPassword: dbPathAndPassword,
		createReader:      createReader,
		currentPath:       "",
	}
}
