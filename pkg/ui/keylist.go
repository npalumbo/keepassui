package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"keepassui/pkg/keepass"
	"log"
)

type KeyList struct {
	dbPathAndPassword binding.Untyped
	elements          binding.StringList
	listWidget        *widget.List
}

func (k *KeyList) DataChanged() {
	err := k.elements.Set([]string{})
	if err != nil {
		log.Printf("Error initialising elements %v", err)
	}
	o, err := k.dbPathAndPassword.Get()
	if err == nil && o != nil {
		d, ok := o.(Data)
		if ok {
			if d.Password != "" {
				secrets, err := keepass.ReadEntries(d.Path, d.Password)

				if err != nil {
					log.Printf("Error reading secret entries: %v", err)
				} else {
					for _, secretEntry := range secrets {
						err = k.elements.Append(secretEntry.Path + " | " + secretEntry.Title)
						if err != nil {
							log.Printf("Error appending entries to list %v", err)
						}
					}
				}
			}
		}
	}
}

func CreatekeyList(dbPathAndPassword binding.Untyped) KeyList {
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
	}
}
