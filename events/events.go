package events

//Event represents mayboy events between ui elements
type Event int

//ListenFunc represents functions that should be called when an event was fired
type ListenFunc func()

const (
	//ProjectSelected - a project has been selected in projectList
	ProjectSelected Event = iota

	//IssueTableFocusLost - the issue table has lost the focus
	IssueTableFocusLost
	//IssueTableRefreshed - the issueTables data has refreshed its contents
	IssueTableRefreshed
	//IssueTableLineSelectionChanged - the issueTable has selected another row
	IssueTableLineSelectionChanged

	//ExitRequested - user requested to exit the app
	ExitRequested
)

//EventManager handles ui events. It can register and fire events
type EventManager struct {
	subscriber map[Event][]ListenFunc
}

//NewEventManager creates a new EventManager
func NewEventManager() *EventManager {
	return &EventManager{
		subscriber: map[Event][]ListenFunc{},
	}
}

//Listen add a ListenFunc to a specific Event
func (em *EventManager) Listen(e Event, fn ListenFunc) {
	if _, found := em.subscriber[e]; !found {
		em.subscriber[e] = []ListenFunc{}
	}

	em.subscriber[e] = append(em.subscriber[e], fn)
}

//Fire calls all ListenFuncs that listen to Event
func (em *EventManager) Fire(e Event) {
	if _, found := em.subscriber[e]; found {
		for _, fn := range em.subscriber[e] {
			fn()
		}

	}
}
