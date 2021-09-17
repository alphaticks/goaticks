package alphaticks

import (
	"fmt"
	"gitlab.com/tachikoma.ai/tickstore-go-client/query"
	"time"
)

type QuerySettings struct {
	settings *query.Settings
	sel      string
	tags     map[string]string
}

func format(sel string, tags map[string]string) string {
	sel = fmt.Sprintf("SELECT %s", sel)
	if tags != nil {
		sel += " WHERE "
		for k, v := range tags {
			sel += k + `="` + v + `" `
		}
	}
	return sel
}

func NewQuery() *QuerySettings {
	return &QuerySettings{settings: query.NewQuerySettings()}
}

func (q *QuerySettings) WithFrom(time time.Time) {
	q.settings.WithFrom(uint64(time.UnixNano() / 1000000))
}

func (q *QuerySettings) WithTo(time time.Time) {
	q.settings.WithTo(uint64(time.UnixNano() / 1000000))
}

func (q *QuerySettings) WithSelector(sel string) {
	q.sel = sel
	q.settings.WithSelector(format(q.sel, q.tags))
}

func (q *QuerySettings) WithTags(tags map[string]string) {
	q.tags = tags
	q.settings.WithSelector(format(q.sel, q.tags))
}

func (q *QuerySettings) WithSamplingFrequency(dur time.Duration) {
	q.settings.Sampler = query.TickSampler{Interval: uint64(dur.Milliseconds())}
}

func (q *QuerySettings) WithStreaming(stream bool) {
	q.settings.WithStreaming(stream)
}

func (q *QuerySettings) WithTimeout(duration time.Duration) {
	q.settings.Timeout = uint64(duration.Nanoseconds())
}
