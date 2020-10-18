package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/yah01/CyDrive/model"
)

type FileInfoItem struct {
	widget.Box
	fileInfo model.FileInfo
	onTap func()
}

func NewFileInfoItem(fileInfo model.FileInfo, tap func()) *FileInfoItem {
	res := fyne.StaticResource{
		StaticName:    theme.FileIcon().Name(),
		StaticContent: theme.FileIcon().Content(),
	}
	icon := canvas.NewImageFromResource(&res)
	icon.FillMode = canvas.ImageFillOriginal

	label:=widget.NewLabel(fileInfo.FilePath)
	label.Alignment = fyne.TextAlignCenter

	layout.NewVBoxLayout()
	item := FileInfoItem{
		Box: widget.Box{
			Children:   []fyne.CanvasObject{
				icon,
				label,
			},
		},
		fileInfo: fileInfo,
		onTap: tap,
	}
	item.ExtendBaseWidget(&item)
	return &item
}

func (label *FileInfoItem) Tapped(e *fyne.PointEvent) {
	label.onTap()
}
