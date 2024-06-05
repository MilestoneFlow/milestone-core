package flows

import (
	"milestone_core/tours/tracker"
)

type Analytics struct {
	Tracker tracker.Tracker
}

func (s Analytics) GetFlowAnalytics(flow *Flow) (FlowAnalytics, error) {
	if flow == nil {
		return FlowAnalytics{
			FlowID:       "",
			Views:        0,
			AvgTotalTime: 0,
			AvgStepTime:  nil,
		}, nil
	}

	analytics := FlowAnalytics{
		FlowID:       flow.ID.Hex(),
		Views:        0,
		AvgTotalTime: 0,
		AvgStepTime:  make(map[string]int64),
	}

	events, err := s.Tracker.FetchTrackDataForFlow(flow.ID.Hex())
	if err != nil {
		return analytics, err
	}

	analytics.Views = s.getUniqueViews(events)
	analytics.AvgStepTime = s.getStepAvgTime(events)

	avgTotalTime := int64(0)
	for _, stepTime := range analytics.AvgStepTime {
		avgTotalTime += stepTime
	}
	analytics.AvgTotalTime = avgTotalTime
	analytics.NoOfFinished = s.getNoOfFinished(events)
	analytics.NoOfSkipped = s.getNoOfSkipped(events)

	return analytics, nil
}

func (s Analytics) getUniqueViews(events []tracker.EventTrack) int {
	seenUserIds := make(map[string]bool)

	for _, event := range events {
		seenUserIds[event.ExternalUserID] = true
	}

	return len(seenUserIds)
}

func (s Analytics) getNoOfFinished(events []tracker.EventTrack) int {
	cnt := 0

	for _, event := range events {
		if event.EventType == tracker.EventTypeFlowFinished {
			cnt++
		}
	}

	return cnt
}

func (s Analytics) getNoOfSkipped(events []tracker.EventTrack) int {
	cnt := 0

	for _, event := range events {
		if event.EventType == tracker.EventTypeFlowSkipped {
			cnt++
		}
	}

	return cnt
}

func (s Analytics) getStepAvgTime(events []tracker.EventTrack) map[string]int64 {
	stepTimesPerUserId := make(map[string]map[string]int64)
	startStepSeen := make(map[string]map[string]bool)
	finishStepSeen := make(map[string]map[string]bool)

	for _, event := range events {
		if event.EventType != tracker.EventTypeFlowStepStart && event.EventType != tracker.EventTypeFlowStepFinish {
			continue
		}

		if _, ok := event.Metadata["stepId"]; !ok {
			continue
		}
		stepId := event.Metadata["stepId"]

		if _, ok := stepTimesPerUserId[stepId]; !ok {
			stepTimesPerUserId[stepId] = make(map[string]int64)
			startStepSeen[stepId] = make(map[string]bool)
			finishStepSeen[stepId] = make(map[string]bool)
		}

		if _, ok := stepTimesPerUserId[stepId][event.ExternalUserID]; !ok {
			stepTimesPerUserId[stepId][event.ExternalUserID] = 0
			startStepSeen[stepId][event.ExternalUserID] = false
			finishStepSeen[stepId][event.ExternalUserID] = false
		}

		if event.EventType == tracker.EventTypeFlowStepStart && !startStepSeen[stepId][event.ExternalUserID] {
			stepTimesPerUserId[stepId][event.ExternalUserID] -= event.Timestamp
			startStepSeen[stepId][event.ExternalUserID] = true
		}
		if event.EventType == tracker.EventTypeFlowStepFinish && !finishStepSeen[stepId][event.ExternalUserID] {
			stepTimesPerUserId[stepId][event.ExternalUserID] += event.Timestamp
			finishStepSeen[stepId][event.ExternalUserID] = true
		}
	}

	stepAvgTime := make(map[string]int64)
	for stepId, times := range stepTimesPerUserId {
		var total int64
		for _, time := range times {
			total += time
		}
		stepAvgTime[stepId] = total / int64(len(times))
	}

	return stepAvgTime
}
