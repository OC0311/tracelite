package tracelite

import (
	"encoding/json"
	"sync"
	"time"
)

// Trace represents a trace instance that can track multiple sub-traces
type Trace struct {
	sync.RWMutex
	name      string         // Name of the trace
	tags      map[string]any // Key-value pairs for additional trace information
	traceList map[string]trace // Map of sub-traces
	traceAt   time.Time     // Time when trace was created
	status    int           // Status of the trace
	totalCost int64         // Total time cost in milliseconds
	openTrace bool          // Flag to control if tracing is enabled
}

// trace represents a single sub-trace with its spans
type trace struct {
	name string         // Name of the sub-trace
	tags map[string]any // Key-value pairs for additional trace information
	list []span        // List of spans in this trace
	cost int64         // Total time cost of this trace
}

// span represents a single time point in a trace
type span struct {
	action        string    // Action name or description
	extensionInfo string    // Additional information about the action
	cost          int64     // Time cost in milliseconds
	martAt        time.Time // Timestamp when the span was marked
}

// TraceResult represents the final result of a trace collection
type TraceResult struct {
	TraceName string         `json:"trace_name"` // Name of the main trace
	Tags      map[string]any `json:"tags"`      // Tags associated with the trace
	TotalCost int64                             // Total time cost of all traces
	TraceSet  []map[string]TraceResultItem `json:"trace_set"` // Collection of trace results
}

// ToString converts TraceResult to JSON string
func (t *TraceResult) ToString() string {
	b, err := json.Marshal(&t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// TraceResultItem represents a single trace result
type TraceResultItem struct {
	TraceCost int64                  `json:"trace_cost"` // Total cost of this trace
	Tags      map[string]interface{} `json:"tags"`       // Tags associated with this trace
	List      [][]interface{}        `json:"list"`       // List of span information
}

// NewTrace creates a new Trace instance with the given name
func NewTrace(name string) *Trace {
	return &Trace{
		name:      name,
		tags:      make(map[string]interface{}),
		traceList: make(map[string]trace),
	}
}

// BeginTrace starts a new trace with the given name and tags
func (t *Trace) BeginTrace(traceName string, tags map[string]interface{}) {
	t.Lock()
	defer t.Unlock()
	if !t.openTrace {
		return
	}
	if tags == nil {
		tags = make(map[string]interface{})
	}
	if _, ok := t.traceList[traceName]; !ok {
		t.traceList[traceName] = trace{
			name: traceName,
			tags: tags,
			list: make([]span, 0),
		}
	}

	s := t.traceList[traceName]
	s.list = append(t.traceList[traceName].list, span{
		action:        "Begin",
		extensionInfo: "",
		martAt:        time.Now().UTC(),
	})
	t.traceList[traceName] = s
}

// SetTags sets the tags for the main trace
func (t *Trace) SetTags(tags map[string]interface{}) {
	t.Lock()
	defer t.Unlock()
	t.tags = tags
}

// TraceOn enables tracing
func (t *Trace) TraceOn() {
	t.Lock()
	defer t.Unlock()
	t.openTrace = true
}

// TraceOff disables tracing
func (t *Trace) TraceOff() {
	t.Lock()
	defer t.Unlock()
	t.openTrace = false
}

// Mark adds a new span to the specified trace
func (t *Trace) Mark(traceName, action, ext string) {
	now := time.Now()
	t.Lock()
	defer t.Unlock()
	if !t.openTrace {
		return
	}
	if _, ok := t.traceList[traceName]; !ok {
		return
	}
	s := t.traceList[traceName]
	s.list = append(t.traceList[traceName].list, span{
		action:        action,
		extensionInfo: ext,
		martAt:        now,
	})
	t.traceList[traceName] = s
}

// Collect gathers all trace information and returns a TraceResult
func (t *Trace) Collect() *TraceResult {
	t.Lock()
	defer t.Unlock()
	if !t.openTrace {
		return nil
	}
	result := &TraceResult{
		TraceName: t.name,
		Tags:      t.tags,
		TraceSet:  make([]map[string]TraceResultItem, 0),
	}

	for k, v := range t.traceList {
		var (
			traceCost    int64
			traceItemSet = make(map[string]TraceResultItem)
		)
		traceItem := TraceResultItem{
			Tags: v.tags,
		}
		for i := 1; i < len(v.list); i++ {
			currentSpan := v.list[i]
			preSpan := v.list[i-1]
			cost := currentSpan.martAt.Sub(preSpan.martAt).Milliseconds()
			traceCost += cost
			item := []interface{}{
				currentSpan.action,
				cost,
				currentSpan.extensionInfo,
				currentSpan.martAt,
			}
			traceItem.List = append(traceItem.List, item)
		}
		traceItem.TraceCost = traceCost

		// 将traceItem放入traceItemSet
		traceItemSet[k] = traceItem
		// 将traceItemSet放入result.TraceSet
		result.TraceSet = append(result.TraceSet, traceItemSet)
		result.TotalCost += traceCost
	}
	return result
}

// CollectToString collects trace information and formats it using the provided function
func (t *Trace) CollectToString(fmtFunc func(*TraceResult) string) string {
	result := t.Collect()
	if result == nil {
		return ""
	}
	return fmtFunc(result)
}
