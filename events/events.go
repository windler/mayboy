package events

type Event int

type ListenFunc func()

const (
	ProjectSelected Event = iota

	IssueTableFocusLost
	IssueTableRefreshed
	IssueTableLineSelectionChanged

	ExitRequested
)

type EventManager struct {
	subscriber map[Event][]ListenFunc
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscriber: map[Event][]ListenFunc{},
	}
}

func (em *EventManager) Listen(e Event, fn ListenFunc) {
	if _, found := em.subscriber[e]; !found {
		em.subscriber[e] = []ListenFunc{}
	}

	em.subscriber[e] = append(em.subscriber[e], fn)
}

func (em *EventManager) Fire(e Event) {
	if _, found := em.subscriber[e]; found {
		for _, fn := range em.subscriber[e] {
			fn()
		}

	}
}
