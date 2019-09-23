package global

import "time"

// GetZoneTime returns the time in teh given zone NOW
func GetZoneNow(zone string) time.Time {
	if dur, found := ZoneOffset[zone]; found {
		return time.Now().Add(dur)
	}
	return time.Now()
}

// GetZoneTime converts the given time into the given timezone time
func GetZoneTime(t time.Time, zone string) time.Time {
	if dur, found := ZoneOffset[zone]; found {
		return t.Add(dur)
	}
	return t
}
