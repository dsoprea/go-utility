package rigeo

import (
	"github.com/dsoprea/go-logging"
	"github.com/golang/geo/s2"
)

func S2CellIdFromCoordinates(latitude, longitude float64) uint64 {
	ll := s2.LatLngFromDegrees(latitude, longitude)
	cellId := s2.CellIDFromLatLng(ll)

	if cellId.IsValid() == false {
		log.Panicf("S2CellIdFromCoordinates: final cell-ID not valid: (%.6f,%.6f) -> (%d)", latitude, longitude, uint64(cellId))
	}

	return uint64(cellId)
}

func S2TokenFromCoordinates(latitude, longitude float64) string {
	ll := s2.LatLngFromDegrees(latitude, longitude)
	cellId := s2.CellIDFromLatLng(ll)

	if cellId.IsValid() == false {
		log.Panicf("S2TokenFromCoordinates: final cell-ID not valid: (%.6f,%.6f) -> (%d)", latitude, longitude, uint64(cellId))
	}

	token := cellId.ToToken()
	return token
}
