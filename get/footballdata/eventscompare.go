package footballdata

import "footballresult/types"

func CompareEvents(first, second []types.Event) []types.Event {
	// Создаем карту для быстрого поиска элементов первого массива
	eventMap := make(map[int64]struct{})
	for _, event := range first {
		eventMap[event.EventID] = struct{}{}
	}

	var missingEvents []types.Event

	// Проверяем, какие элементы второго массива отсутствуют в первом
	for _, event := range second {
		if _, exists := eventMap[event.EventID]; !exists {
			missingEvents = append(missingEvents, event)
		}
	}

	return missingEvents
}
