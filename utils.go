package lrq

import "time"

func timeToInt64(t *time.Time) *int64 {
	if t == nil {
		return nil
	} else {
		n := t.Unix()
		return &n
	}
}

