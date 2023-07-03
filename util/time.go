package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration accept in format hh:mm:ss.xxx
func ParseDuration(str string) (time.Duration, error) {
	timeSplits := strings.Split(str, ":")
	if len(timeSplits) != 3 {
		return 0, fmt.Errorf("invalid input")
	}

	var (
		hour, minute, second, ms int64
	)

	hour, err := strconv.ParseInt(timeSplits[0], 10, 64)
	if err != nil {
		return 0, err
	}

	minute, err = strconv.ParseInt(timeSplits[1], 10, 64)
	if err != nil {
		return 0, err
	}

	secondsSplit := strings.Split(timeSplits[2], ".")
	second, err = strconv.ParseInt(secondsSplit[0], 10, 64)
	if err != nil {
		return 0, err
	}

	if len(secondsSplit) > 1 {
		ms, err = strconv.ParseInt(secondsSplit[1], 10, 64)
		if err != nil {
			return 0, err
		}
		if time.Duration(ms) > time.Second {
			return 0, fmt.Errorf("invalid millisecond")
		}
	}

	duration := hour*int64(time.Hour) + minute*int64(time.Minute) + second*int64(time.Second) + ms
	return time.Duration(duration), nil
}
