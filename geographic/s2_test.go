package rigeo

import (
	"testing"
)

func TestS2CellFromCoordinates(t *testing.T) {
	actualCell := S2CellFromCoordinates(26.568629, -80.108965)
	expectedCellId := uint64(9860956140720671831)

	if uint64(actualCell) != expectedCellId {
		t.Fatalf("Cell-ID not correct: (%d) != (%d)", uint64(actualCell), expectedCellId)
	}
}
