package main

import "fmt"

type LeapYearCalculator struct {
	leapYears []int
}

func NewLeapYearCalculator(startYear int, endYear int) *LeapYearCalculator {
	leapYears := getLeapYears(startYear, endYear)

	return &LeapYearCalculator{
		leapYears: leapYears,
	}
}

func getLeapYears(startYear int, endYear int) []int {
	leapYears := []int{}
	for year := startYear - 1; year <= endYear; year++ {
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			leapYears = append(leapYears, year)
		}
	}

	return leapYears
}

func (l *LeapYearCalculator) GetLeapDays(month int, day int) []string {
	leapDays := make([]string, 0, len(l.leapYears))
	for _, year := range l.leapYears {
		leapDays = append(leapDays, fmt.Sprintf("%04d%02d%02d", year, month, day))
	}

	return leapDays
}

func (l *LeapYearCalculator) GetFirstLeapYear() (int, error) {
	if len(l.leapYears) == 0 {
		return -1, fmt.Errorf("No leap years found")
	}

	return l.leapYears[0], nil
}
