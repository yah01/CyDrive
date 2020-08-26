package page

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

type DrivePage struct {
	*fyne.Container
	FileList *fyne.Container
}

func NewDrivePage(fileList *fyne.Container) *DrivePage {
	return &DrivePage{Container: fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(4), fileList), FileList: fileList}
}
