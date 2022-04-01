package utils

import (
	"math/rand"
	"strings"
	"time"
)

func IsSameDay(nanoSec1, nanoSec2 int64) bool {
	year1, month1, day1 := time.Unix(0, nanoSec1).Date()
	year2, month2, day2 := time.Unix(0, nanoSec2).Date()

	if year1 == year2 && month1 == month2 && day1 == day2 {
		return true
	}

	return false
}

func GetRandInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randVal := min + rand.Intn(max-min+1)
	return randVal
}

func IsStrEmpty(str *string) bool {
	if str == nil || *str == "" {
		return true
	}

	strTrim := strings.Trim(*str, " ")
	if len(strTrim) == 0 {
		return true
	}

	return false
}

func GetEpochFromDay(day int) int {
	return day * 24 * 60 * 2
}

func GetCurrentEpoch() int {
	currentNanoSec := time.Now().UnixNano()
	currentEpoch := (currentNanoSec/1e9 - 1598306471) / 30
	return int(currentEpoch)
}
