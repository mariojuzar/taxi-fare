package fare

const (
	flatFareRule        = "FLAT"
	progressiveFareRule = "PROGRESSIVE"
)

var (
	fareRules = []fareRule{
		{
			minDistance:    0,
			maxDistance:    1000,
			ruleType:       flatFareRule,
			price:          400,
			distanceCharge: 0,
		},
		{
			minDistance:    1000,
			maxDistance:    10000,
			ruleType:       progressiveFareRule,
			price:          40,
			distanceCharge: 400,
		},
		{
			minDistance:    10000,
			maxDistance:    UNLIMITED,
			ruleType:       progressiveFareRule,
			price:          40,
			distanceCharge: 350,
		},
	}
)

type fareRule struct {
	minDistance    float64
	maxDistance    float64
	ruleType       string
	price          int64
	distanceCharge float64
}
