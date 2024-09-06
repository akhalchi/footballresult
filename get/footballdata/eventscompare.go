package footballdata

import "footballresult/types"

func CompareEvents(first, second []types.Event) []types.Event {

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
