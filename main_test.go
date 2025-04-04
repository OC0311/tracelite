package tracelite

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewTrace(t *testing.T) {
	trace := NewTrace("test_trace")
	if trace.name != "test_trace" {
		t.Errorf("Expected trace name to be 'test_trace', got %s", trace.name)
	}
	if trace.tags == nil {
		t.Error("Expected tags to be initialized")
	}
	if trace.traceList == nil {
		t.Error("Expected traceList to be initialized")
	}
}

func TestTraceOnOff(t *testing.T) {
	trace := NewTrace("test_trace")
	
	trace.TraceOn()
	if !trace.openTrace {
		t.Error("Expected trace to be enabled after TraceOn")
	}
	
	trace.TraceOff()
	if trace.openTrace {
		t.Error("Expected trace to be disabled after TraceOff")
	}
}

func TestSetTags(t *testing.T) {
	trace := NewTrace("test_trace")
	tags := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	
	trace.SetTags(tags)
	
	if len(trace.tags) != len(tags) {
		t.Errorf("Expected tags length to be %d, got %d", len(tags), len(trace.tags))
	}
	
	for k, v := range tags {
		if trace.tags[k] != v {
			t.Errorf("Expected tag %s to be %v, got %v", k, v, trace.tags[k])
		}
	}
}

func TestBeginTraceAndMark(t *testing.T) {
	trace := NewTrace("test_trace")
	trace.TraceOn()
	
	traceName := "subtrace1"
	tags := map[string]interface{}{"tag1": "value1"}
	
	// Test BeginTrace
	trace.BeginTrace(traceName, tags)
	
	if _, ok := trace.traceList[traceName]; !ok {
		t.Error("Expected trace to be created")
	}
	
	// Test Mark
	trace.Mark(traceName, "action1", "ext1")
	trace.Mark(traceName, "action2", "ext2")
	
	subTrace := trace.traceList[traceName]
	if len(subTrace.list) != 3 { // Begin + 2 marks
		t.Errorf("Expected 3 spans, got %d", len(subTrace.list))
	}
}

func TestCollect(t *testing.T) {
	trace := NewTrace("test_trace")
	trace.TraceOn()
	
	// Setup test data
	traceName := "subtrace1"
	tags := map[string]interface{}{"tag1": "value1"}
	trace.BeginTrace(traceName, tags)
	
	// Add some delay to simulate real timing
	time.Sleep(10 * time.Millisecond)
	trace.Mark(traceName, "action1", "ext1")
	time.Sleep(10 * time.Millisecond)
	trace.Mark(traceName, "action2", "ext2")
	
	// Collect results
	result := trace.Collect()
	
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	
	if result.TraceName != "test_trace" {
		t.Errorf("Expected trace name to be 'test_trace', got %s", result.TraceName)
	}
	
	if len(result.TraceSet) != 1 {
		t.Errorf("Expected 1 trace set, got %d", len(result.TraceSet))
	}
}

func TestCollectToString(t *testing.T) {
	trace := NewTrace("test_trace")
	trace.TraceOn()
	
	traceName := "subtrace1"
	trace.BeginTrace(traceName, nil)
	trace.Mark(traceName, "action1", "ext1")
	
	// Test with default ToString function
	result := trace.CollectToString(func(tr *TraceResult) string {
		b, _ := json.Marshal(tr)
		return string(b)
	})
	
	if result == "" {
		t.Error("Expected non-empty string result")
	}
	
	// Test when trace is off
	trace.TraceOff()
	result = trace.CollectToString(func(tr *TraceResult) string {
		return "should not reach here"
	})
	
	if result != "" {
		t.Error("Expected empty string when trace is off")
	}
}