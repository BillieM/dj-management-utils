package iwidget

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
TrackListSection widget contains a list of tracks via a TrackList widget,
along with filter and sort controls for viewing that list, it also contains
controls to export the tracks to a playlist
*/
type TrackListSection struct {
	*Base
	widget.BaseWidget

	List                    *widget.List
	TrackListControls       *TrackListControls
	TrackListExportControls *TrackListImportExportControls
}

func NewTrackListSection(widgetBase *Base, tlb *TrackListBinding, selectedTrack *SelectedTrackBinding, trackListFuncs TrackListFuncs) *TrackListSection {

	tlb.FilterSortInfo = &FilterSortInfo{}

	trackListSection := &TrackListSection{
		Base:                    widgetBase,
		List:                    NewTrackList(widgetBase, tlb, selectedTrack),
		TrackListControls:       NewTrackListControls(widgetBase, tlb),
		TrackListExportControls: NewTrackListImportExportControls(widgetBase, trackListFuncs),
	}

	trackListSection.ExtendBaseWidget(trackListSection)

	return trackListSection

}

func (t *TrackListSection) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(
			t.TrackListControls,
			t.TrackListExportControls,
			nil, nil,
			t.List,
		),
	)
}

type FilterSortInfo struct {
}

/*
TrackListControls contains controls for filtering and sorting a TrackList
*/
type TrackListControls struct {
	*Base
	widget.BaseWidget

	trackListBinding *TrackListBinding

	TrackSortControls   *TrackListSortControls
	TrackFilterControls *TrackListFilterControls
}

func NewTrackListControls(widgetBase *Base, tlb *TrackListBinding) *TrackListControls {

	applyFilterSortCallback := func() {
		tlb.ApplyFilterSort()
	}

	tlc := &TrackListControls{
		Base:                widgetBase,
		TrackSortControls:   NewTrackSortControls(widgetBase, tlb.FilterSortInfo, applyFilterSortCallback),
		TrackFilterControls: NewTrackListFilterControls(widgetBase, tlb.FilterSortInfo, applyFilterSortCallback),

		trackListBinding: tlb,
	}

	tlc.ExtendBaseWidget(tlc)

	return tlc
}

func (i *TrackListControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewGridWithColumns(
			2,
			container.NewVBox(
				widget.NewLabel("Sort"),
				widget.NewSeparator(),
				i.TrackSortControls,
			),
			container.NewVBox(
				widget.NewLabel("Filter"),
				widget.NewSeparator(),
				i.TrackFilterControls,
			),
		),
	)
}

type TrackListSortControls struct {
	*Base
	widget.BaseWidget

	SortBy *widget.Select
	Desc   *widget.Check
}

func NewTrackSortControls(widgetBase *Base, fsi *FilterSortInfo, callback func()) *TrackListSortControls {
	tlsc := &TrackListSortControls{
		Base:   widgetBase,
		SortBy: widget.NewSelect([]string{"Default", "Name", "Genre", "Tags", "Publisher", "SoundCloud User"}, nil),
		Desc:   widget.NewCheck("Descending", nil),
	}

	tlsc.ExtendBaseWidget(tlsc)

	return tlsc
}

func (i *TrackListSortControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.SortBy,
			i.Desc,
		),
	)
}

type TrackListFilterControls struct {
	*Base
	widget.BaseWidget

	ShowLinked *widget.Check
}

func NewTrackListFilterControls(widgetBase *Base, fsi *FilterSortInfo, callback func()) *TrackListFilterControls {
	tlfc := &TrackListFilterControls{
		Base:       widgetBase,
		ShowLinked: widget.NewCheck("Show Linked", nil),
	}

	tlfc.ExtendBaseWidget(tlfc)

	return tlfc
}

func (i *TrackListFilterControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.ShowLinked,
		),
	)
}

/*
TrackListExportControls contains controls for exporting a TrackList to a playlist
within DJ software
*/
type TrackListImportExportControls struct {
	*Base
	widget.BaseWidget

	refreshButton *widget.Button

	lTemp *widget.Label
}

func NewTrackListImportExportControls(widgetBase *Base, trackListFuncs TrackListFuncs) *TrackListImportExportControls {
	tliec := &TrackListImportExportControls{
		Base:  widgetBase,
		lTemp: widget.NewLabel("TrackListImportExportControls"),
		refreshButton: widget.NewButtonWithIcon("Refresh playlist", theme.ViewRefreshIcon(), func() {
			trackListFuncs.RefreshSoundCloudPlaylist()
		}),
	}

	tliec.ExtendBaseWidget(tliec)

	return tliec
}

func (i *TrackListImportExportControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.lTemp,
			i.refreshButton,
		),
	)
}

func NewTrackList(widgetBase *Base, tlb *TrackListBinding, selectedTrack *SelectedTrackBinding) *widget.List {

	trackList := widget.NewListWithData(
		tlb,
		func() fyne.CanvasObject {
			return NewTrackListItem("FairlyLongTrackNameTemplateIncase")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			trackListItem := o.(*TrackListItem)
			trackBinding := i.(*TrackBinding)

			trackListItem.TrackName.SetText(trackBinding.Track.Name)
			if trackBinding.Track.LocalPath != "" && !trackBinding.Track.LocalPathBroken {
				trackListItem.Linked.SetResource(theme.ConfirmIcon())
			} else {
				trackListItem.Linked.SetResource(theme.ContentRemoveIcon())
			}
		},
	)

	trackList.OnSelected = func(id widget.ListItemID) {
		if selectedTrack.Locked {
			trackList.Select(selectedTrack.ListID)
			dialog.ShowError(helpers.ErrPleaseWaitForDownload, widgetBase.MainWindow)
			return
		}

		tli, err := tlb.GetItem(id)
		if err != nil {
			fmt.Println("error getting track from list", err)
			return
		}
		selectedTrackBind := tli.(*TrackBinding)

		selectedTrack.TrackBinding = selectedTrackBind
		selectedTrack.ListID = id
		selectedTrack.Trigger()
	}

	return trackList
}

type TrackListItem struct {
	widget.BaseWidget

	TrackName *widget.Label
	Linked    *widget.Icon
}

func NewTrackListItem(name string) *TrackListItem {
	tli := &TrackListItem{
		TrackName: widget.NewLabel(name),
		Linked:    widget.NewIcon(theme.ContentRemoveIcon()),
	}

	tli.ExtendBaseWidget(tli)

	return tli
}

func (i *TrackListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(
			nil, nil,
			i.Linked,
			nil,
			i.TrackName,
		),
	)
}

type bindBase struct {
	sync.RWMutex
	listeners sync.Map // map[DataListener]bool
}

type TrackListBinding struct {
	bindBase

	FilterSortInfo *FilterSortInfo

	Tracks        []*streaming.SoundCloudTrack
	VisibleTracks []*TrackBinding
}

func (i *TrackListBinding) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *TrackListBinding) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *TrackListBinding) GetItem(index int) (binding.DataItem, error) {
	i.Lock()
	defer i.Unlock()
	if index < 0 || index >= len(i.VisibleTracks) {
		return nil, helpers.ErrIndexOutOfBounds
	}

	return i.VisibleTracks[index], nil
}

func (i *TrackListBinding) Length() int {
	i.Lock()
	defer i.Unlock()
	return len(i.VisibleTracks)
}

func (i *TrackListBinding) Append(p *TrackBinding) {
	i.Lock()
	defer i.Unlock()
	i.VisibleTracks = append(i.VisibleTracks, p)
}

func (i *TrackListBinding) Set(p []*streaming.SoundCloudTrack) {
	i.Lock()
	defer i.Unlock()
	i.Tracks = p
}

/*
ApplyFilterSort applies the current filter and sort settings to the list of tracks

This uses the list of tracks (i.Tracks) attached to the widget, placing the filtered and sorted
tracks into i.VisibleTracks
*/
func (i *TrackListBinding) ApplyFilterSort() {
	i.Lock()
	defer i.Unlock()
	// TODO
	i.VisibleTracks = []*TrackBinding{}
	for _, t := range i.Tracks {
		i.VisibleTracks = append(i.VisibleTracks, &TrackBinding{Track: t})
	}
}

/*
trackBinding is a binding.DataItem for a SoundCloudTrack
*/
type TrackBinding struct {
	bindBase

	// may want a context in here ?? later problem...
	Track *streaming.SoundCloudTrack
}

func (i *TrackBinding) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *TrackBinding) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

type TrackListFuncs struct {
	RefreshSoundCloudPlaylist func()
}
