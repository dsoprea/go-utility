package rigeo

import (
    "github.com/golang/geo/s2"
)

func S2CellIdFromCoordinates(latitude, longitude float64) uint64 {
    ll := s2.LatLngFromDegrees(latitude, longitude)
    cellId := s2.CellIDFromLatLng(ll)

    return uint64(cellId)
}
