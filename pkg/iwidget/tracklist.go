package iwidget

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
)

/*
TrackListSection widget contains a list of tracks via a TrackList widget,
along with filter and sort controls for viewing that list, it also contains
controls to export the tracks to a playlist
*/
type TrackListSection struct {
	widget.BaseWidget

	*TrackList
	*TrackListControls
	*TrackListExportControls
}

func NewTrackListSection(tlb *TrackListBinding, selectedTrack *TrackBinding) *TrackListSection {

	tlb.FilterSortInfo = &FilterSortInfo{}

	trackListSection := &TrackListSection{
		TrackList:               NewTrackList(tlb, selectedTrack),
		TrackListControls:       NewTrackListControls(tlb),
		TrackListExportControls: NewTrackListExportControls(),
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
			t.TrackList,
		),
	)
}

type FilterSortInfo struct {
}

/*
TrackListControls contains controls for filtering and sorting a TrackList
*/
type TrackListControls struct {
	widget.BaseWidget

	trackListBinding *TrackListBinding

	TrackSortControls   *TrackListSortControls
	TrackFilterControls *TrackListFilterControls
}

func NewTrackListControls(tlb *TrackListBinding) *TrackListControls {

	applyFilterSortCallback := func() {
		tlb.ApplyFilterSort()
	}

	tlc := &TrackListControls{
		TrackSortControls:   NewTrackSortControls(tlb.FilterSortInfo, applyFilterSortCallback),
		TrackFilterControls: NewTrackListFilterControls(tlb.FilterSortInfo, applyFilterSortCallback),

		trackListBinding: tlb,
	}

	tlc.ExtendBaseWidget(tlc)

	return tlc
}

func (i *TrackListControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewHBox(
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
	widget.BaseWidget

	SortBy *widget.Select
	Desc   *widget.Check
}

func NewTrackSortControls(fsi *FilterSortInfo, callback func()) *TrackListSortControls {
	tlsc := &TrackListSortControls{
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
	widget.BaseWidget

	ShowLinked *widget.Check
}

func NewTrackListFilterControls(fsi *FilterSortInfo, callback func()) *TrackListFilterControls {
	tlfc := &TrackListFilterControls{
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
type TrackListExportControls struct {
	widget.BaseWidget

	lTemp *widget.Label
}

func NewTrackListExportControls() *TrackListExportControls {
	tlec := &TrackListExportControls{
		lTemp: widget.NewLabel("TrackListExportControls"),
	}

	tlec.ExtendBaseWidget(tlec)

	return tlec
}

func (i *TrackListExportControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		i.lTemp,
	)
}

/*
TrackList widget contains a list of tracks, selecting a track will display the
tracks info in a Track widget to the side
*/
type TrackList struct {
	*widget.List
}

func NewTrackList(tlb *TrackListBinding, selectedTrack *TrackBinding) *TrackList {

	trackList := widget.NewListWithData(
		tlb,
		func() fyne.CanvasObject {
			return NewTrackListItem("FairlyLongTrackNameTemplateIncase")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			// TODO
			trackListItem := o.(*TrackListItem)
			trackBinding := i.(*TrackBinding)

			trackListItem.TrackName.SetText(trackBinding.track.Name)
			if trackBinding.track.LocalPath != "" && !trackBinding.track.LocalPathBroken {
				trackListItem.Linked.SetResource(theme.ConfirmIcon())
			} else {
				trackListItem.Linked.SetResource(theme.CancelIcon())
			}
		},
	)

	return &TrackList{
		trackList,
	}
}

type TrackListItem struct {
	widget.BaseWidget

	OnSelected func()

	TrackName *widget.Label
	Linked    *widget.Icon
}

func NewTrackListItem(name string) *TrackListItem {
	tli := &TrackListItem{
		TrackName: widget.NewLabel(name),
		Linked:    widget.NewIcon(theme.CancelIcon()),
	}

	tli.ExtendBaseWidget(tli)

	return tli
}

func (i *TrackListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewHBox(
			i.Linked,
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

	Tracks        []*database.SoundCloudTrack
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

func (i *TrackListBinding) Set(p []*database.SoundCloudTrack) {
	i.Lock()
	defer i.Unlock()
	i.Tracks = p
}

func (i *TrackListBinding) ApplyFilterSort() {
	i.Lock()
	defer i.Unlock()
	// TODO
	clear(i.VisibleTracks)
	for _, t := range i.Tracks {
		i.VisibleTracks = append(i.VisibleTracks, &TrackBinding{track: t})
	}
}

/*
trackBinding is a binding.DataItem for a SoundCloudTrack
*/
type TrackBinding struct {
	bindBase

	// may want a context in here ?? later problem...
	track *database.SoundCloudTrack
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
