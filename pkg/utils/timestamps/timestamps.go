package timestamps

import (
	"fmt"
	"strings"
	"time"
)

func ParseTimestamp(t string, now time.Time, defaultT ...time.Time) (string, error) {
	defaultTimestamp := now
	if len(defaultT) > 0 {
		defaultTimestamp = defaultT[0]
	}
	ts, err := func() (time.Time, error) {
		fields := strings.Split(t, "-")
		if len(fields) == 0 {
			return defaultTimestamp, nil
		}
		if len(fields) > 1 && fields[0] != "now" {
			return defaultTimestamp, fmt.Errorf("bad timestamp: %s", t)
		}
		if len(fields) == 1 {
			return now, nil
		}
		if len(fields) == 2 {
			d, err := time.ParseDuration(fields[1])
			if err != nil {
				return defaultTimestamp, fmt.Errorf("bad timestamp: %s", t)
			}
			return now.Add(-d), nil
		}
		return defaultTimestamp, fmt.Errorf("bad timestamp: %s", t)
	}()
	if err != nil {
		return "", nil
	}
	return ts.UTC().Format(time.RFC3339), nil
}
