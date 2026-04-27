package observe

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixed }

func newTestObserver() (*Observer, *bytes.Buffer) {
	var buf bytes.Buffer
	o := newWithClock(&buf, fixedClock)
	return o, &buf
}

func TestEmit_RecordsEvent(t *testing.T) {
	o, _ := newTestObserver()
	o.Emit("INFO", "secret/app", "synced successfully")
	events := o.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != "INFO" {
		t.Errorf("expected level INFO, got %s", events[0].Level)
	}
	if events[0].Path != "secret/app" {
		t.Errorf("unexpected path: %s", events[0].Path)
	}
}

func TestEmit_WritesToWriter(t *testing.T) {
	o, buf := newTestObserver()
	o.Info("secret/db", "no changes detected")
	if !strings.Contains(buf.String(), "no changes detected") {
		t.Errorf("expected output to contain message, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("expected output to contain level INFO")
	}
}

func TestWarn_RecordsWarnLevel(t *testing.T) {
	o, _ := newTestObserver()
	o.Warn("secret/app", "TTL nearing expiry")
	if o.CountByLevel("WARN") != 1 {
		t.Errorf("expected 1 WARN event")
	}
}

func TestError_RecordsErrorLevel(t *testing.T) {
	o, _ := newTestObserver()
	o.Error("secret/app", "failed to fetch")
	if o.CountByLevel("ERROR") != 1 {
		t.Errorf("expected 1 ERROR event")
	}
}

func TestCountByLevel_MultipleEvents(t *testing.T) {
	o, _ := newTestObserver()
	o.Info("a", "msg")
	o.Info("b", "msg")
	o.Warn("c", "msg")
	o.Error("d", "msg")
	if o.CountByLevel("INFO") != 2 {
		t.Errorf("expected 2 INFO events")
	}
	if o.CountByLevel("WARN") != 1 {
		t.Errorf("expected 1 WARN event")
	}
	if o.CountByLevel("ERROR") != 1 {
		t.Errorf("expected 1 ERROR event")
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	o, _ := newTestObserver()
	o.Info("x", "hello")
	events := o.Events()
	events[0].Level = "MUTATED"
	original := o.Events()
	if original[0].Level != "INFO" {
		t.Errorf("Events() should return a copy, original was mutated")
	}
}
