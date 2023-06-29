package mylayout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type BlankLayout struct {
}

type FullWidthLayout struct {
}

func NewBlankLayout() fyne.Layout {
	return &BlankLayout{}
}

func NewFullWidthLayout() fyne.Layout {
	return &FullWidthLayout{}
}

func (d *BlankLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

func (d *BlankLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		childsize := o.MinSize()
		o.Resize(childsize)
	}
}

func (d *FullWidthLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

func (d *FullWidthLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		childsize := o.Size()
		o.Resize(fyne.NewSize(size.Width-theme.InnerPadding(), childsize.Height))
	}
}
