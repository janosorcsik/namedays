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

func (cw *CalendarWriter) WriteEvent(
	year int,
	month int,
	day int,
	isLeap bool,
	summary string,
	description string,
	rrule string,
	exdate string,
) {
	fmt.Fprint(cw.w, "BEGIN:VEVENT\r\n")

	var uid string

	if isLeap {
		uid = fmt.Sprintf("nameday-%02d%02d-leap@calendar.local", month, day)
	} else {
		uid = fmt.Sprintf("nameday-%02d%02d@calendar.local", month, day)
	}

	fmt.Fprintf(cw.w, "UID:%s\r\n", uid)
	fmt.Fprintf(cw.w, "DTSTART;VALUE=DATE:%d%02d%02d\r\n", year, month, day)
	fmt.Fprintf(cw.w, "SUMMARY:%s\r\n", summary)
	fmt.Fprintf(cw.w, "DESCRIPTION:%s\r\n", description)

	if rrule != "" {
		fmt.Fprintf(cw.w, "RRULE:%s\r\n", rrule)
	}

	if exdate != "" {
		fmt.Fprintf(cw.w, "EXDATE;VALUE=DATE:%s\r\n", exdate)
	}

	fmt.Fprint(cw.w, "END:VEVENT\r\n\r\n")
}
