package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Cursorable = (*PasswordRevealerNotDisabled)(nil)
var _ fyne.Tappable = (*PasswordRevealerNotDisabled)(nil)
var _ fyne.Widget = (*PasswordRevealerNotDisabled)(nil)

type PasswordRevealerNotDisabled struct {
	widget.BaseWidget

	icon  *canvas.Image
	entry *widget.Entry
}

func NewPasswordRevealerNotDisabled(e *widget.Entry) *PasswordRevealerNotDisabled {
	pr := &PasswordRevealerNotDisabled{
		icon:  canvas.NewImageFromResource(theme.VisibilityOffIcon()),
		entry: e,
	}
	pr.ExtendBaseWidget(pr)
	return pr
}

func (r *PasswordRevealerNotDisabled) CreateRenderer() fyne.WidgetRenderer {
	return &passwordRevealerRenderer{
		WidgetRenderer: widget.NewSimpleRenderer(r.icon),
		icon:           r.icon,
		entry:          r.entry,
	}
}

func (r *PasswordRevealerNotDisabled) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (r *PasswordRevealerNotDisabled) Tapped(*fyne.PointEvent) {
	r.entry.Password = !r.entry.Password
	r.entry.Refresh()
	// fyne.CurrentApp().Driver().CanvasForObject(r).Focus(r.entry.super().(fyne.Focusable))
}

var _ fyne.WidgetRenderer = (*passwordRevealerRenderer)(nil)

type passwordRevealerRenderer struct {
	fyne.WidgetRenderer
	entry *widget.Entry
	icon  *canvas.Image
}

func (r *passwordRevealerRenderer) Layout(size fyne.Size) {
	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

func (r *passwordRevealerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *passwordRevealerRenderer) Refresh() {
	if !r.entry.Password {
		r.icon.Resource = theme.VisibilityIcon()
	} else {
		r.icon.Resource = theme.VisibilityOffIcon()
	}

	canvas.Refresh(r.icon)
}
