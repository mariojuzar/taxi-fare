package fare

import (
	"github.com/mariojuzar/taxi-fare/util"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TaxiFare interface {
	Calculate(input string) (int64, error)
	ShowRecords()
}

type fare struct {
	currentElapsedTime time.Duration
	currentFare        int64
	currentDistance    float64
	records            []fareRecord
}

type fareRecord struct {
	Timestamp  string
	Mileage    float64
	MileageStr string
	Diff       float64
}

func NewTaxiFare() TaxiFare {
	return &fare{}
}

// Calculate
// input receive with format hh:mm:ss.fff<SPACE>xxxxxxxx.f
// for example: 00:02:00.125 1141.2
func (f *fare) Calculate(input string) (int64, error) {
	logrus.Debugf("receive input: %s", input)

	if input == "" {
		return f.currentFare, ErrBlankLine
	}

	if input == "end" {
		if len(f.records) < 2 {
			return 0, ErrNotEnoughData
		}
		if f.currentDistance == 0.0 {
			return 0, ErrZeroMileage
		}
		return f.currentFare, nil
	}

	err := f.parseInput(input)
	if err != nil {
		logrus.WithError(err).Error("failed to parse input")
		return 0, err
	}

	f.calculateFare()

	return f.currentFare, nil
}

func (f *fare) ShowRecords() {
	recordToShow := f.records
	sort.Slice(recordToShow, func(i, j int) bool {
		return recordToShow[i].Diff > recordToShow[j].Diff
	})

	for _, record := range recordToShow {
		log.Printf("%s %s %s\n", record.Timestamp, record.MileageStr, strconv.FormatFloat(record.Diff, 'f', -1, 64))
	}
}

func (f *fare) parseInput(input string) error {
	inputs := strings.Split(input, " ")
	if len(inputs) != 2 {
		return ErrImproperFormat
	}

	duration, err := util.ParseDuration(inputs[0])
	if err != nil {
		return ErrImproperFormat
	}

	if duration < f.currentElapsedTime {
		return ErrPastTime
	}

	if (duration - f.currentElapsedTime) > (5 * time.Minute) {
		return ErrIntervalTimeTooLong
	}

	distance, err := strconv.ParseFloat(inputs[1], 64)
	if err != nil {
		return ErrImproperFormat
	}

	if len(f.records) == 0 && distance > 0 {
		return ErrInvalidStart
	}

	distanceDec := decimal.NewFromFloat(distance)
	diff := distanceDec.Sub(decimal.NewFromFloat(f.currentDistance))
	exactDiff, _ := diff.Float64()

	f.records = append(f.records, fareRecord{
		Timestamp:  inputs[0],
		Mileage:    distance,
		MileageStr: inputs[1],
		Diff:       exactDiff,
	})
	f.currentElapsedTime = duration
	f.currentDistance = distance

	return nil
}

func (f *fare) calculateFare() {
	var currentFare int64 = 0

	// applying fare rule
	for _, rule := range fareRules {
		if f.currentDistance > rule.minDistance {
			// switch based on valid rule to apply
			switch rule.ruleType {
			case flatFareRule:
				currentFare += rule.price
			case progressiveFareRule:
				distanceCharge := 0.0
				if f.currentDistance >= rule.maxDistance && rule.maxDistance != maxDistance {
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
