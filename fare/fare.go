package fare

import (
	"errors"
	"github.com/mariojuzar/taxi-fare/util"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	UNLIMITED = -999
)

var (
	ErrImproperFormat      = errors.New("improper input format")
	ErrBlankLine           = errors.New("receive a blank input")
	ErrPastTime            = errors.New("receive past time")
	ErrIntervalTimeTooLong = errors.New("time interval is more than 5 minutes from latest data")
	ErrNotEnoughData       = errors.New("records is less than 2")
	ErrZeroMileage         = errors.New("invalid movement mileage")
)

type TaxiFare interface {
	CalculateFareMeter(input string) (int64, error)
}

type fare struct {
	currentElapsedTime time.Duration
	currentFare        int64
	currentDistance    float64
	countHistory       uint
}

func NewTaxiFare() TaxiFare {
	return &fare{}
}

func (f *fare) CalculateFareMeter(input string) (int64, error) {
	if input == "" {
		return 0, ErrBlankLine
	}

	if input == "finished" {
		if f.countHistory < 2 {
			return 0, ErrNotEnoughData
		}
		if f.currentDistance == 0.0 {
			return 0, ErrZeroMileage
		}
		return f.currentFare, nil
	}

	inputs := strings.Split(input, " ")
	if len(inputs) != 2 {
		return 0, ErrImproperFormat
	}

	duration, err := util.ParseDuration(inputs[0])
	if err != nil {
		return 0, ErrImproperFormat
	}

	if duration < f.currentElapsedTime {
		return 0, ErrPastTime
	}

	if (duration - f.currentElapsedTime) > (5 * time.Minute) {
		return 0, ErrIntervalTimeTooLong
	}

	distance, err := strconv.ParseFloat(inputs[1], 64)
	if err != nil {
		return 0, ErrImproperFormat
	}

	f.currentElapsedTime = duration
	f.currentDistance = distance
	f.countHistory++
	f.calculateFare()

	return f.currentFare, nil
}

func (f *fare) calculateFare() {
	var currentFare int64 = 0

	// applying fare rule
	for _, rule := range fareRules {
		if f.currentDistance > rule.minDistance {
			// valid rule to apply
			switch rule.ruleType {
			case flatFareRule:
				currentFare += rule.price
			case progressiveFareRule:
				distanceCharge := 0.0
				if f.currentDistance >= rule.maxDistance && rule.maxDistance != UNLIMITED {
					distanceCharge = rule.maxDistance - rule.minDistance
				} else {
					distanceCharge = f.currentDistance - rule.minDistance
				}
				cnt := int64(math.Ceil(distanceCharge / rule.distanceCharge))
				currentFare += cnt * rule.price
			}
		}
	}

	f.currentFare = currentFare
}
