package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	INPUT     = "namedays.json"
	OUTPUT    = "namedays.ics"
	STARTYEAR = 2025
	ENDYEAR   = 2060
)

type Root map[string]Month

type Month map[string]Day

type Day struct {
	Names []string `json:"names"`
}

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
		idx := strconv.Itoa(m)
		month := root[idx]

		for d := 1; d <= 31; d++ {
			idx := strconv.Itoa(d)
			day, ok := month[idx]

			if !ok {
				continue
			}

			summary := strings.Join(day.Names, ", ")

			if m == 2 && d >= 24 {
				leapDays := calculator.GetLeapDays(m, d)
				excludedDates := strings.Join(leapDays, ",")

				rrule := fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", m, d)
				writer.WriteEvent(STARTYEAR, m, d, false, summary, "Hungarian name day (non-leap years)", rrule, excludedDates)

				if d == 24 {
					summary = "Szökőnap"
				} else if d != 28 {
					nextIdx := strconv.Itoa(d - 1)
					nextDay, ok := month[nextIdx]

					if !ok {
						continue
					}

					summary = strings.Join(nextDay.Names, ", ")
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
