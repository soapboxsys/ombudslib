package ombwire

import "time"

// This file contains helpful functions that extend the generated protobuf
// types.

// Passed location can be nil.
func NewBulletin(msg string, ts uint64, loc *Location) *Bulletin {
	bltn := &Bulletin{
		Message:   &msg,
		Timestamp: &ts,
	}
	if loc != nil {
		bltn.Location = loc
	}
	return bltn
}

func NewLocation(lat, lon, h float64) *Location {
	loc := &Location{
		Lat: &lat,
		Lon: &lon,
		H:   &h,
	}
	return loc
}

func NewBulletinFromStr(msg string) *Bulletin {
	return NewBulletin(msg, uint64(time.Now().Unix()), nil)
}
