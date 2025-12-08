package main

import (
	"fmt"
	"io"
)

type CalendarWriter struct {
	w io.Writer
}

func NewCalendarWriter(w io.Writer) *CalendarWriter {
	return &CalendarWriter{w: w}
}

func (cw *CalendarWriter) WriteHeader() {
	fmt.Fprint(cw.w, "BEGIN:VCALENDAR\r\n")
	fmt.Fprint(cw.w, "VERSION:2.0\r\n")
	fmt.Fprint(cw.w, "PRODID:-//Name Days//Hungarian Calendar//EN\r\n")
	fmt.Fprint(cw.w, "CALSCALE:GREGORIAN\r\n")
	fmt.Fprint(cw.w, "METHOD:PUBLISH\r\n\r\n")
}

func (cw *CalendarWriter) WriteFooter() {
	fmt.Fprint(cw.w, "END:VCALENDAR\r\n")
}

type CalendarEvent struct {
	Year          int
	Month         int
	Day           int
	IsLeapYear    bool
	Summary       string
	Description   string
	Rule          *string
	ExcludedDates *string
}

func (cw *CalendarWriter) WriteEvent(event CalendarEvent) {
	fmt.Fprint(cw.w, "BEGIN:VEVENT\r\n")

	var uid string

	if event.IsLeapYear {
		uid = fmt.Sprintf("nameday-%02d%02d-leap@calendar.local", event.Month, event.Day)
	} else {
		uid = fmt.Sprintf("nameday-%02d%02d@calendar.local", event.Month, event.Day)
	}

	fmt.Fprintf(cw.w, "UID:%s\r\n", uid)
	fmt.Fprintf(cw.w, "DTSTART;VALUE=DATE:%d%02d%02d\r\n", event.Year, event.Month, event.Day)
	fmt.Fprintf(cw.w, "SUMMARY:%s\r\n", event.Summary)
	fmt.Fprintf(cw.w, "DESCRIPTION:%s\r\n", event.Description)

	if event.Rule != nil {
		fmt.Fprintf(cw.w, "RRULE:%s\r\n", *event.Rule)
	}

	if event.ExcludedDates != nil {
		fmt.Fprintf(cw.w, "EXDATE;VALUE=DATE:%s\r\n", *event.ExcludedDates)
	}

	fmt.Fprint(cw.w, "END:VEVENT\r\n\r\n")
}
