package rigeo

import (
    "testing"
)

func TestS2CellIdFromCoordinates(t *testing.T) {
    actualCellId := S2CellIdFromCoordinates(26.568629, -80.108965)
    expectedCellId := uint64(9860956140720671831)

    if actualCellId != expectedCellId {
        t.Fatalf("Cell-ID not correct.")
    }
}
