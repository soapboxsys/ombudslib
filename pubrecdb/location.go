package pubrecdb

import (
	"math"

	"github.com/soapboxsys/ombudslib/ombjson"
)

var (
	selectNearbyBltns string = bltnSql + `
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN endorsements ON bulletins.txid = endorsements.bid
		WHERE bulletins.latitude IS NOT NULL AND bulletins.longitude IS NOT NULL AND
			dist($1, $2, bulletins.latitude, bulletins.longitude) < $3
		GROUP BY bulletins.txid HAVING bulletins.txid NOT null
		ORDER BY blocks.timestamp DESC
	`
)

// GetNearbyBltns returns bulletins that were tagged with a location within r
// meters of lat, lon. The bulletin are ordered by block timestamp and are NOT
// sorted by distance from the point.
func (db *PublicRecord) GetNearbyBltns(lat, lon, r float64) ([]*ombjson.Bulletin, error) {

	rows, err := db.selectNearbyBltns.Query(lat, lon, r)
	defer rows.Close()
	if err != nil {
		return []*ombjson.Bulletin{}, err
	}

	bltns, err := scanBltns(rows)
	if err != nil {
		return []*ombjson.Bulletin{}, err
	}

	return bltns, nil
}

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
