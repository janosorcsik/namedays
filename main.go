package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	INPUT     = "namedays.json"
	OUTPUT    = "namedays.ics"
	STARTYEAR = 2025
	ENDYEAR   = 2060
)

type Root map[int]Month

type Month map[int]Day

type Day []string

func main() {
	jsonFile, err := os.Open(INPUT)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer jsonFile.Close()

	var root Root
	err = json.NewDecoder(jsonFile).Decode(&root)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	calendar, err := os.Create(OUTPUT)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer calendar.Close()

	calculator := NewLeapYearCalculator(STARTYEAR, ENDYEAR)
	firstLeapYear, err := calculator.GetFirstLeapYear()

	if err != nil {
		fmt.Println("Error getting first leap year:", err)
		return
	}

	writer := NewCalendarWriter(calendar)
	writer.WriteHeader()

	for m := 1; m <= 12; m++ {
		month := root[m]

		for d := 1; d <= 31; d++ {
			day, ok := month[d]

			if !ok {
				continue
			}

			summary := strings.Join(day, ", ")

			if m == 2 && d >= 24 {
				leapDays := calculator.GetLeapDays(m, d)
				excludedDates := strings.Join(leapDays, ",")

				rrule := fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", m, d)
				writer.WriteEvent(STARTYEAR, m, d, false, summary, "Hungarian name day (non-leap years)", rrule, excludedDates)

				if d == 24 {
					summary = "Szökőnap"
				} else if d != 28 {
					prevDay, ok := month[d-1]

					if !ok {
						continue
					}

					summary = strings.Join(prevDay, ", ")
				}

				rrule = fmt.Sprintf("FREQ=YEARLY;INTERVAL=4;BYMONTH=%d;BYMONTHDAY=%d", m, d)
				writer.WriteEvent(firstLeapYear, m, d, true, summary, "Hungarian name day (leap years)", rrule, "")

				if d == 28 {
					rrule = fmt.Sprintf("FREQ=YEARLY;INTERVAL=4;BYMONTH=%d;BYMONTHDAY=%d", m, d+1)
					writer.WriteEvent(firstLeapYear, m, d+1, true, summary, "Hungarian name day (leap years)", rrule, "")
				}
			} else {
				writer.WriteEvent(STARTYEAR, m, d, false, summary, "Hungarian name day", "FREQ=YEARLY", "")
			}
		}
	}

	writer.WriteFooter()
}
