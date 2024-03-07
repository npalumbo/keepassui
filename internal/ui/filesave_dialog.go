package ui

import (
	"keepassui/internal/secretsreader"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/dchest/uniuri"
)

//go:generate mockgen -destination=../mocks/ui/mock_filesave_dialog.go -source=./filesave_dialog.go

type FileSaver interface {
	ShowForMasterPassword(masterPassword string)
}

type DefaultFileSaver struct {
	stagerController StagerController
	secretsReader    secretsreader.SecretReader
	parent           fyne.Window
	notify           binding.String
}

// ShowForMasterPassword implements FileSaver.
func (d DefaultFileSaver) ShowForMasterPassword(masterPassword string) {

	fileSaveDialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil {
			handleErrorAndGoToHomeView(err, d.parent, d.stagerController) // Error from closing dialog
			return
		}

		if uc == nil {
			goToHomeView(d.stagerController, d.parent) // This happens when you cancel the save dialog
			return
		}

		emptyDBBytes, err := d.secretsReader.CreateEmptyDBBytes(masterPassword)

		if err != nil {
			handleErrorAndGoToHomeView(err, d.parent, d.stagerController) // Error creating empty DB
			return
		}

		_, err = uc.Write(emptyDBBytes)
		if err != nil {
			handleErrorAndGoToHomeView(err, d.parent, d.stagerController) // Error writing contents to file
			return
		}

		err = uc.Close()

		if err != nil {
			handleErrorAndGoToHomeView(err, d.parent, d.stagerController) // Error from closing writer
			return
		}

		dfsr, ok := d.secretsReader.(*secretsreader.DefaultSecretsReader)

		if !ok {
			goToHomeView(d.stagerController, d.parent)
			return
		}
		dfsr.ContentInBytes = emptyDBBytes
		dfsr.Password = masterPassword
		dfsr.UriID = uc.URI().String()

		err = d.notify.Set(uniuri.New())
		if err != nil {
			handleErrorAndGoToHomeView(err, d.parent, d.stagerController) // Error setting value to notify binding
		}

	}, d.parent)

	// At the moment fyne doesn't allow enforcing the extension of the filename.
	// I tried several ways of manipulating the URI (directly modifying the URI and with storage.Parent / storage.Child)
	// and it doesn't work for android, so for the time being I will leave it for the user to add the extension when
	// creating a new DB. For reference:
	// https://github.com/fyne-io/fyne/issues/4692
	// https://github.com/fyne-io/fyne/issues/1044
	fileSaveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	fileSaveDialog.SetFileName("new_db.kdbx")

	fileSaveDialog.Show()
}

func CreateFileSaver(secretsReader secretsreader.SecretReader, stagerController StagerController, parent fyne.Window) DefaultFileSaver {
	return DefaultFileSaver{
		stagerController: stagerController,
		secretsReader:    secretsReader,
		parent:           parent,
		notify:           binding.NewString(),
	}
}

func (f *DefaultFileSaver) AddListener(l binding.DataListener) {
	f.notify.AddListener(l)
}
