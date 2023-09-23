package ui

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"keepassui/pkg/keepass"
	"log/slog"
)

type KeyList struct {
	dbPathAndPassword binding.Untyped
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
		if ok {
			if d.Password != "" {
				secrets, err := keepass.ReadEntries(d.Path, d.Password)

				if err != nil {
					slog.Error("Error reading secret entries", err)
					dialog.ShowError(errors.New("Error reading secrets: "+err.Error()), k.parent)
				} else {
					for _, secretEntry := range secrets {
						err = k.elements.Append(secretEntry.Path + " | " + secretEntry.Title)
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

func CreatekeyList(dbPathAndPassword binding.Untyped, parent fyne.Window) KeyList {
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
		elements:          elements,
		listWidget:        listWidget,
		parent:            parent,
	}
}
