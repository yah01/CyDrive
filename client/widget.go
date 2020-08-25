package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/yah01/CyDrive/model"
)

type FileInfoLabel struct {
	widget.Label
	fileInfo model.FileInfo
	onTap func(fileInfo model.FileInfo)
}

func NewTappableLabel(fileInfo model.FileInfo, tap func(fileInfo model.FileInfo)) *FileInfoLabel {
	label := FileInfoLabel{
		Label: widget.Label{
			Text: fileInfo.FilePath,
		},
		fileInfo: fileInfo,
		onTap: tap,
	}
	label.ExtendBaseWidget(&label)
	return &label
}

func (label *FileInfoLabel) Tapped(e *fyne.PointEvent) {
	label.onTap(label.fileInfo)
}
