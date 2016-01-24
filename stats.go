package main

import "time"

type Search struct{}

type Filterer interface {
	Filter([]*RequestEvent) []*RequestEvent
}

type TimeFilter struct {
	from time.Time
	to   time.Time
}

func (f *TimeFilter) Filter(events []*RequestEvent) []*RequestEvent {
	filtered := make([]*RequestEvent, 0, len(events))
	for _, e := range events {
		if f.from.Before(e.RecordedAt) && f.to.After(e.RecordedAt) {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

type Grouper interface {
	Group([]*RequestEvent) map[string][]*RequestEvent
}

type RouteGrouper struct{}

func (g *RouteGrouper) Group(events []*RequestEvent) map[string][]*RequestEvent {
	groupings := map[string][]*RequestEvent{}
	for _, e := range events {
		path := e.Req.URL.Path
		if _, ok := groupings[path]; ok {
			groupings[path] = append(groupings[path], e)
		} else {
			groupings[path] = []*RequestEvent{e}
		}
	}

	return groupings
}

// all zereod
type ResponseTimeAnalysis struct {
	Min  time.Duration
	Max  time.Duration
	Mean time.Duration
}

func analyzeResponseTimes(events []*RequestEvent) *ResponseTimeAnalysis {
	if len(events) <= 0 {
		return nil
	}

	ra := &ResponseTimeAnalysis{Min: events[0].UpstreamDuration, Max: events[0].UpstreamDuration}
	var sum time.Duration
	for _, e := range events {
		if ra.Min > e.UpstreamDuration {
			ra.Min = e.UpstreamDuration
		}

		if ra.Max < e.UpstreamDuration {
			ra.Max = e.UpstreamDuration
		}

		sum += e.UpstreamDuration
	}

	ra.Mean = sum / time.Duration(len(events))
	return ra
}

// rate limiter proxy
// see where your visitors come from
