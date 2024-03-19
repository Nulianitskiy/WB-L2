package event

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EventInMemoRepo struct {
	EventStorage map[string]Event
	UserStorage  map[int64][]string
	mu           *sync.Mutex
}

func NewInMemoryRepo() *EventInMemoRepo {
	newRepo := &EventInMemoRepo{
		EventStorage: make(map[string]Event, 10),
		UserStorage:  make(map[int64][]string, 10),
		mu:           &sync.Mutex{},
	}
	return newRepo
}

func (e *EventInMemoRepo) Create(newEvent Event) (Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	newEvent.ID = uuid.New().String()

	if _, ok := e.EventStorage[newEvent.ID]; ok {
		return newEvent, fmt.Errorf("event already exists")
	}
	e.EventStorage[newEvent.ID] = newEvent
	e.UserStorage[newEvent.UserID] = append(e.UserStorage[newEvent.UserID], newEvent.ID)
	return newEvent, nil
}

func (e *EventInMemoRepo) Update(newEvent Event) (Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.EventStorage[newEvent.ID]; !ok {
		return newEvent, fmt.Errorf("no such event")
	}
	e.EventStorage[newEvent.ID] = newEvent
	return newEvent, nil
}

func (e *EventInMemoRepo) getUserEvents(userID int64) ([]string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var ar []string
	var ok bool

	if ar, ok = e.UserStorage[userID]; !ok {
		return nil, fmt.Errorf("no such user")
	}
	return ar, nil
}

func (e *EventInMemoRepo) deleteFromUserMap(userID int64, eventID string) error {
	ar, err := e.getUserEvents(userID)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	deleteInd := -1
	for i, v := range ar {
		if v == eventID {
			deleteInd = i
			break
		}
	}
	if deleteInd == -1 {
		return fmt.Errorf("no such event")
	}
	ar[deleteInd], ar[len(ar)-1] = ar[len(ar)-1], ar[deleteInd]
	e.UserStorage[userID] = ar[:len(ar)-1]
	return nil
}

func (e *EventInMemoRepo) Delete(userID int64, eventID string) error {
	err := e.deleteFromUserMap(userID, eventID)
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.EventStorage[eventID]; !ok {
		return fmt.Errorf("no such event")
	}
	delete(e.EventStorage, eventID)
	return nil
}

func (e *EventInMemoRepo) GetEventsForRange(userID int64, start time.Time, end time.Time) ([]Event, error) {
	eventsID, err := e.getUserEvents(userID)
	if err != nil {
		return nil, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	ans := make([]Event, 0, len(e.EventStorage))
	for _, v := range eventsID {
		if e.EventStorage[v].Date == start || e.EventStorage[v].Date.After(start) && e.EventStorage[v].Date.Before(end) {
			ans = append(ans, e.EventStorage[v])
		}
	}

	return ans, nil
}
