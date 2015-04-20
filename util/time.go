package util

import "time"

const ISOTimeFormat = "2006-01-02T15:04:05Z"

func ISOTime(t time.Time) string {
	return t.UTC().Format(ISOTimeFormat)
}

func ParseISOTime(str string) (time.Time, error) {
	return time.Parse(ISOTimeFormat, str)
}
