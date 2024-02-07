package ui

import (
	"errors"
	"keepassui/internal/keepass"
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
	fullContainer     *fyne.Container
	navTop            *fyne.Container
	saveButton        *widget.Button
	breadCrumbs       *fyne.Container
	generalButtons    *fyne.Container
	listPanel         *fyne.Container
	detailedView      *DetailedView
	parent            fyne.Window
	dbPathAndPassword *DBPathAndPassword
	createReader      ToSecretReaderFn
	currentPath       string
	secretsDB         *keepass.SecretsDB
}

func (n *NavView) DataChanged() {
	if n.dbPathAndPassword.UriID == "" {
		return
	}

	secretReader := n.createReader(*n.dbPathAndPassword)

	secretsDB, err := secretReader.ReadEntriesFromContentGroupedByPath()

	if err != nil {
		dialog.ShowError(errors.New("Error reading secrets: "+err.Error()), n.parent)
		return
	}

	n.secretsDB = &secretsDB

	n.UpdateNavView(secretsDB.PathsInOrder[0])

	n.saveButton.OnTapped = func() {
		bytes, err := secretsDB.WriteDBBytes(n.dbPathAndPassword.Password)

		if err != nil {
			dialog.ShowError(err, n.parent)
			return
		}
		fileSaveDialog := createFileSaveDialog(bytes, n.dbPathAndPassword.UriID, n.parent)

		if fileSaveDialog != nil {
			fileSaveDialog.Show()
		}
	}

	n.fullContainer.Show()
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
			deleteButton := widget.NewButtonWithIcon("delete", theme.DeleteIcon(), func() {})
			buttons := container.NewHBox(copyButton, showInfoButton, openGroupButton, deleteButton)
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
			deleteButton := buttons.Objects[3].(*widget.Button)
			deleteButton.OnTapped = func() {
				deleted := navView.secretsDB.DeleteSecretEntry(secret)
				if deleted {
					navView.UpdateNavView(secret.Group)
				}
			}

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

func CreateNavView(dbPathAndPassword *DBPathAndPassword, detailedView *DetailedView, parent fyne.Window, createReader ToSecretReaderFn) NavView {

	breadCrumbs := container.NewHBox()
	generalButtons := container.NewHBox()

	navTop := container.NewBorder(nil, nil, breadCrumbs, generalButtons, nil)

	saveButton := widget.NewButtonWithIcon("save", theme.DocumentSaveIcon(), func() {

	})
	generalButtons.Add(saveButton)

	listPanel := container.NewStack()
	fullContainer := container.NewBorder(container.NewVBox(navTop, widget.NewSeparator()), nil, nil, nil, listPanel)
	fullContainer.Hide()

	return NavView{
		fullContainer:     fullContainer,
		breadCrumbs:       breadCrumbs,
		listPanel:         listPanel,
		detailedView:      detailedView,
		parent:            parent,
		dbPathAndPassword: dbPathAndPassword,
		createReader:      createReader,
		currentPath:       "",
		generalButtons:    generalButtons,
		navTop:            navTop,
		saveButton:        saveButton,
	}
}
