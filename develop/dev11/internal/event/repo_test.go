package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	repo := NewInMemoryRepo()
	event := Event{UserID: 1231}
	_, err := repo.Create(event)
	if err != nil {
		t.Errorf("unexpected error")
	}
	if len(repo.EventStorage) != 1 && len(repo.UserStorage) != 1 {
		t.Errorf("not right")
	}
}

func TestUpdate(t *testing.T) {
	repo := NewInMemoryRepo()

	tests := []struct {
		iEvents   map[string]Event
		iUsers    map[int64][]string
		oEvents   map[string]Event
		oUsers    map[int64][]string
		needError bool
		args      []interface{}
		name      string
	}{
		{
			name:      "no such event",
			iEvents:   nil,
			iUsers:    nil,
			oEvents:   nil,
			oUsers:    nil,
			args:      []interface{}{Event{ID: "111"}},
			needError: true,
		},
		{
			name:    "ok",
			iEvents: map[string]Event{"1111": Event{ID: "1111", Data: "asdasdas"}},
			iUsers:  nil,
			oEvents: map[string]Event{"1111": Event{ID: "1111", Data: "asd"}},
			oUsers:  nil,
			args:    []interface{}{Event{ID: "1111", Data: "asd"}},
		},
	}

	for _, ts := range tests {
		repo.EventStorage = ts.iEvents
		repo.UserStorage = ts.iUsers
		_, err := repo.Update(ts.args[0].(Event))
		if err != nil && ts.needError == true {
			continue
		}
		if err == nil && ts.needError == true {
			t.Errorf("expected error got nil")
		}
		if err != nil && ts.needError == false {
			t.Errorf("expected nil got error")
		}

		assert.Equal(t, ts.oEvents, repo.EventStorage)
	}
}

func TestDelete(t *testing.T) {
	repo := NewInMemoryRepo()

	tests := []struct {
		iEvents   map[string]Event
		iUsers    map[int64][]string
		oEvents   map[string]Event
		oUsers    map[int64][]string
		needError bool
		args      []interface{}
		name      string
	}{
		{
			name:      "no such event",
			iEvents:   nil,
			iUsers:    nil,
			oEvents:   nil,
			oUsers:    nil,
			args:      []interface{}{1, "111"},
			needError: true,
		},
		{
			name:    "ok",
			iEvents: map[string]Event{"1111": Event{ID: "1111", Data: "asdasdas", UserID: 123}},
			iUsers:  map[int64][]string{123: []string{"1111"}},
			oEvents: map[string]Event{},
			oUsers:  map[int64][]string{int64(123): []string{}},
			args:    []interface{}{123, "1111"},
		},
	}

	for _, ts := range tests {
		repo.EventStorage = ts.iEvents
		repo.UserStorage = ts.iUsers
		err := repo.Delete(int64(ts.args[0].(int)), ts.args[1].(string))
		if err != nil && ts.needError == true {
			continue
		}
		if err == nil && ts.needError == true {
			t.Errorf("expected error got nil")
		}
		if err != nil && ts.needError == false {
			t.Errorf("expected nil got error")
		}

		assert.Equal(t, ts.oEvents, repo.EventStorage)
	}
}

func TestEventsForRange(t *testing.T) {
	repo := NewInMemoryRepo()

	var events = []Event{
		{
			ID:      "event_1",
			UserID:  123,
			Data:    "Meeting with client",
			RawDate: "2024-02-10",
			Date:    time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      "event_2",
			UserID:  456,
			Data:    "Team building activity",
			RawDate: "2024-02-20",
			Date:    time.Date(2024, time.February, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      "event_3",
			UserID:  789,
			Data:    "Project kickoff",
			RawDate: "2024-02-15",
			Date:    time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      "event_4",
			UserID:  123,
			Data:    "Client presentation",
			RawDate: "2024-02-15",
			Date:    time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      "event_5",
			UserID:  456,
			Data:    "Team meeting",
			RawDate: "2024-02-05",
			Date:    time.Date(2024, time.February, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      "event_6",
			UserID:  123,
			Data:    "Product demo",
			RawDate: "2024-02-11",
			Date:    time.Date(2024, time.February, 11, 0, 0, 0, 0, time.UTC),
		},
	}

	tmpDate := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		iEvents   map[string]Event
		iUsers    map[int64][]string
		oEvents   []Event
		oUsers    map[int64][]string
		needError bool
		args      []interface{}
		name      string
	}{
		{
			name:    "day",
			iEvents: map[string]Event{"event_1": events[0], "event_4": events[3], "event_6": events[5]},
			iUsers:  map[int64][]string{123: []string{"event_1", "event_4", "event_6"}},
			oEvents: []Event{events[0]},
			args:    []interface{}{123, time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC), time.Date(2024, time.February, 11, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:    "week",
			iEvents: map[string]Event{"event_1": events[0], "event_4": events[3], "event_6": events[5]},
			iUsers:  map[int64][]string{123: []string{"event_1", "event_4", "event_6"}},
			oEvents: []Event{events[3], events[5]},
			args:    []interface{}{123, time.Date(2024, time.February, 11, 0, 0, 0, 0, time.UTC), time.Date(2024, time.February, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:    "month",
			iEvents: map[string]Event{"event_1": events[0], "event_4": events[3], "event_6": events[5]},
			iUsers:  map[int64][]string{123: []string{"event_1", "event_4", "event_6"}},
			oEvents: []Event{events[0], events[3], events[5]},
			args:    []interface{}{123, time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC), tmpDate.AddDate(0, 1, 0)},
		},
	}

	for _, ts := range tests {
		repo.EventStorage = ts.iEvents
		repo.UserStorage = ts.iUsers
		ans, err := repo.GetEventsForRange(int64(ts.args[0].(int)), ts.args[1].(time.Time), ts.args[2].(time.Time))
		if err != nil && ts.needError == true {
			continue
		}
		if err == nil && ts.needError == true {
			t.Errorf("expected error got nil")
		}
		if err != nil && ts.needError == false {
			t.Errorf("expected nil got error")
		}

		assert.Equal(t, ts.oEvents, ans)
	}
}
