package ui

import (
	"errors"
	"keepassui/internal/keepass"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type KeyList struct {
	dbPathAndPassword binding.Untyped
	content           binding.Bytes
	elements          binding.StringList
	listWidget        *widget.List
	parent            fyne.Window
}

func (k *KeyList) DataChanged() {
	err := k.elements.Set([]string{})
	if err != nil {
		slog.Error("Error initialising elements", err)
	}
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
					for _, secretEntry := range secrets {
						err = k.elements.Append(secretEntry.Path + " | " + secretEntry.Title + ": " + secretEntry.Password)
						if err != nil {
							slog.Error("Error appending entries to list", err)
						}
					}
				}
			}
		}
	}
	k.listWidget.Resize(k.listWidget.Size().AddWidthHeight(0, float32(k.elements.Length()*20)))
}

func CreatekeyList(dbPathAndPassword binding.Untyped, content binding.Bytes, parent fyne.Window) KeyList {
	elements := binding.NewStringList()
	listWidget := widget.NewListWithData(elements,
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			return label
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			label := item.(*widget.Label)
			strData := data.(binding.String)
			u, _ := strData.Get()
			if u != "" {
				label.SetText(u)
			}
		})

	return KeyList{
		dbPathAndPassword: dbPathAndPassword,
		content:           content,
		elements:          elements,
		listWidget:        listWidget,
		parent:            parent,
	}
}
