package pubrecdb

import "testing"

func TestDistanceFormula(t *testing.T) {

	// Pass known good values and see if the thing works
	tests := []struct {
		a_lat float64
		a_lon float64
		b_lat float64
		b_lon float64
		dist  int
	}{
		{44, 22, 51.5, -0.12, 1842197},
		{44, 22, 10.02, 125, 10252779},
	}

	for _, test := range tests {
		d := distance(test.a_lat, test.a_lon, test.b_lat, test.b_lon)
		if int(d) != int(test.dist) {
			t.Fatalf("Expected: %d, got: %f", test.dist, d)
		}
	}
}
