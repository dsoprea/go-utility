package rigeo

import (
	"github.com/dsoprea/go-logging"
	"github.com/golang/geo/s2"
)

// S2CellFromCoordinates returns an `S2` cell for the given real-world
// coordinates.
func S2CellFromCoordinates(latitude, longitude float64) s2.CellID {
	ll := s2.LatLngFromDegrees(latitude, longitude)
	cellId := s2.CellIDFromLatLng(ll)

	if cellId.IsValid() == false {
		log.Panicf("S2TokenFromCoordinates: final cell-ID not valid: (%.6f,%.6f) -> (%d)", latitude, longitude, uint64(cellId))
	}

	return cellId
}
