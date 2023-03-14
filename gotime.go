package time

import (
	"time"
)

type (
	Time     = time.Time
	Duration = time.Duration
	Location = time.Location
)

var (
	LoadLocation    = time.LoadLocation
	Date            = time.Date
	ParseInLocation = time.ParseInLocation
	Unix            = time.Unix
	Sleep           = time.Sleep
	Parse           = func(layout, value string) (Time, error) { return ParseInLocation(layout, value, GetLocation()) }
	Since           = func(tm Time) Duration { return Now().Sub(tm) }
)

var (
	Local = time.Local
	UTC   = time.UTC
)

const (
	Layout      = time.Layout      // "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	ANSIC       = time.ANSIC       //  "Mon Jan _2 15:04:05 2006"
	UnixDate    = time.UnixDate    //  "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = time.RubyDate    //  "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = time.RFC822      //  "02 Jan 06 15:04 MST"
	RFC822Z     = time.RFC822Z     // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = time.RFC850      // "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = time.RFC1123     // "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = time.RFC1123Z    //  "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = time.RFC3339     //  "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = time.Kitchen     //  "3:04PM"
	// Handy time stamps.
	Stamp      = time.Stamp      // "Jan _2 15:04:05"
	StampMilli = time.StampMilli // "Jan _2 15:04:05.000"
	StampMicro = time.StampMicro //  "Jan _2 15:04:05.000000"
	StampNano  = time.StampNano  // "Jan _2 15:04:05.000000000"
)

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
)
