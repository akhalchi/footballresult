package footballdata

import "footballresult/types"

func CheckExistEvents(first, second []types.Event) []types.Event {

	eventMap := make(map[int64]struct{})
	for _, event := range first {
		eventMap[event.EventID] = struct{}{}
	}

	var missingEvents []types.Event

	for _, event := range second {
		if _, exists := eventMap[event.EventID]; !exists {
			missingEvents = append(missingEvents, event)
		}
	}

	return missingEvents
}

func CompareEvents(event1, event2 types.Event) bool {

	if event1.EventID != event2.EventID ||
		!event1.EventDate.Equal(event2.EventDate) ||
		event1.Tournament != event2.Tournament ||
		event1.TeamHome != event2.TeamHome ||
		event1.TeamAway != event2.TeamAway ||
		event1.GoalsHome != event2.GoalsHome ||
		event1.GoalsAway != event2.GoalsAway ||
		event1.PenHome != event2.PenHome ||
		event1.PenAway != event2.PenAway ||
		event1.RcHome != event2.RcHome ||
		event1.RcAway != event2.RcAway ||
		event1.Importance != event2.Importance ||
		event1.EventStatus != event2.EventStatus ||
		event1.PublishedStatus != event2.PublishedStatus {
		return false
	}
	return true
}
