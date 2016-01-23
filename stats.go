package main

type Search struct{}

type Filterer interface {
	Filter([]*RequestEvent) []*RequestEvent
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

// rate limiter proxy
// see where your visitors come from
