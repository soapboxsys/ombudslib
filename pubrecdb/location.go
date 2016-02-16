package pubrecdb

import "math"

// distance uses the distance formula derived using the spherical law of cosines
// to reasonablely accurate approximations of the distances between 'a' and
// 'b' on the earth's surface. The reference implementation lives at this site:
// http://www.movable-type.co.uk/scripts/latlong.html
func distance(a_lat, a_lon, b_lat, b_lon float64) float64 {
	ToRad := func(x float64) float64 {
		return x * (math.Pi / 180)
	}
	φ1 := ToRad(a_lat)
	φ2 := ToRad(b_lat)
	λ := ToRad(b_lon - a_lon)

	// Earth's radius in meters.
	R := 6371000.0
	z := math.Sin(φ1)*math.Sin(φ2) + math.Cos(φ1)*math.Cos(φ2)*math.Cos(λ)
	d := math.Acos(z) * R

	return d
}

// Exists to help us test our implementation. Computes x^y
func pow(x, y int64) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}
