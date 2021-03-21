package crushftp

import (
	"testing"
	"time"
)

func TestUnixDateTimeParse(t *testing.T) {
	tests := map[string]time.Time{
		"Tue Jan 02 12:14:16 CET 2021": time.Date(2021, time.January, 02, 12, 14, 16, 0, time.FixedZone("CET", int(time.Hour*1/time.Second))),
	}

	for dateString, expectedTime := range tests {
		u := unixDateTime{}
		u.Parse(dateString)

		if expectedTime.Year() != u.Year() {
			t.Errorf("unix date time mismatch Year on %q (expected %q, got %q)", dateString, expectedTime.Year(), u.Year())
			continue
		}
		if expectedTime.Month() != u.Month() {
			t.Errorf("unix date time mismatch Month on %q (expected %q, got %q)", dateString, expectedTime.Month(), u.Month())
			continue
		}
		if expectedTime.Day() != u.Day() {
			t.Errorf("unix date time mismatch Day on %q (expected %q, got %q)", dateString, expectedTime.Day(), u.Day())
			continue
		}
		if expectedTime.Hour() != u.Hour() {
			t.Errorf("unix date time mismatch Hour on %q (expected %q, got %q)", dateString, expectedTime.Hour(), u.Hour())
			continue
		}
		if expectedTime.Minute() != u.Minute() {
			t.Errorf("unix date time mismatch Minute on %q (expected %q, got %q)", dateString, expectedTime.Minute(), u.Minute())
			continue
		}
		if expectedTime.Second() != u.Second() {
			t.Errorf("unix date time mismatch Second on %q (expected %q, got %q)", dateString, expectedTime.Second(), u.Second())
			continue
		}
		if expectedTime.Nanosecond() != u.Nanosecond() {
			t.Errorf("unix date time mismatch Nanosecond on %q (expected %q, got %q) ", dateString, expectedTime.Nanosecond(), u.Nanosecond())
			continue
		}
	}
}
