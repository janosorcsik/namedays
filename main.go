package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	inputFile          = "namedays.json"
	outputFile         = "namedays.ics"
	startYear          = 2025
	endYear            = 2060
	leapMonth          = 2
	leapDay            = 24
	lastDayInLeapMonth = 28
	monthsInYear       = 12
	daysInMonth        = 31
)

type Root map[int]Month

type Month map[int]Day

type Day []string

func main() {
	jsonFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v", err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	var root Root
	err = json.NewDecoder(jsonFile).Decode(&root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing file: %v", err)
		os.Exit(1)
	}

	calendar, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %v", err)
		os.Exit(1)
	}

	defer calendar.Close()

	calculator := NewLeapYearCalculator(startYear, endYear)
	firstLeapYear, err := calculator.GetFirstLeapYear()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting first leap year: %v", err)
		os.Exit(1)
	}

	writer := NewCalendarWriter(calendar)
	writer.WriteHeader()

	for m := 1; m <= monthsInYear; m++ {
		month := root[m]

		for d := 1; d <= daysInMonth; d++ {
			day, ok := month[d]

			if !ok {
				continue
			}

			summary := strings.Join(day, ", ")

			if m == leapMonth && d >= leapDay {
				leapDays := calculator.GetLeapDays(m, d)
				excludedDates := strings.Join(leapDays, ",")

				rrule := fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", m, d)
				description := "Hungarian name day (non-leap years)"

				event := CalendarEvent{
					Year:          startYear,
					Month:         m,
					Day:           d,
					IsLeapYear:    false,
					Summary:       summary,
					Description:   description,
					Rule:          &rrule,
					ExcludedDates: &excludedDates,
				}

				writer.WriteEvent(event)

				if d == leapDay {
					summary = "Szökőnap"
				} else if d != lastDayInLeapMonth {
					prevDay, ok := month[d-1]

					if !ok {
						continue
					}

					summary = strings.Join(prevDay, ", ")
				}

				rrule = fmt.Sprintf("FREQ=YEARLY;INTERVAL=4;BYMONTH=%d;BYMONTHDAY=%d", m, d)
				description = "Hungarian name day (leap years)"
				event = CalendarEvent{
					Year:          firstLeapYear,
					Month:         m,
					Day:           d,
					IsLeapYear:    true,
					Summary:       summary,
					Description:   description,
					Rule:          &rrule,
					ExcludedDates: nil,
				}

				writer.WriteEvent(event)

				if d == lastDayInLeapMonth {
					rrule = fmt.Sprintf("FREQ=YEARLY;INTERVAL=4;BYMONTH=%d;BYMONTHDAY=%d", m, d+1)
					event := CalendarEvent{
						Year:          firstLeapYear,
						Month:         m,
						Day:           d + 1,
						IsLeapYear:    true,
						Summary:       summary,
						Description:   description,
						Rule:          &rrule,
						ExcludedDates: nil,
					}

					writer.WriteEvent(event)
				}
			} else {
				rrule := "FREQ=YEARLY"
				event := CalendarEvent{
					Year:          startYear,
					Month:         m,
					Day:           d,
					IsLeapYear:    false,
					Summary:       summary,
					Description:   "Hungarian name day",
					Rule:          &rrule,
					ExcludedDates: nil,
				}

				writer.WriteEvent(event)
			}
		}
	}

	writer.WriteFooter()
}
