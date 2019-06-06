package ridata

import (
	"reflect"
	"sort"
	"testing"

	"github.com/dsoprea/go-logging"
)

type testLruItem struct {
	id LruKey
}

func (tli testLruItem) Id() LruKey {
	return tli.id
}

func CreateTestLru(dropped *int) *Lru {
	droppedCb := func(raw LruKey) (err error) {
		id := raw.(int)
		*dropped = id

		return nil
	}

	lru := NewLru(5)
	lru.SetDropCb(droppedCb)

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	tli3 := testLruItem{
		id: 33,
	}

	_, _, err = lru.Set(tli3)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	tli4 := testLruItem{
		id: 44,
	}

	_, _, err = lru.Set(tli4)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	size := lru.Count()
	if size != 4 {
		log.Panicf("LRU size not correct: (%d)", size)
	}

	return lru
}

func TestLru_Set(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	var dropped int

	lru := CreateTestLru(&dropped)

	node1 := lru.top
	node2 := node1.after
	node3 := node2.after
	node4 := node3.after

	if node1.item.Id() != 44 {
		t.Fatalf("First node not correct: [%v]", node1.item.Id())
	} else if node2.item.Id() != 33 {
		t.Fatalf("Second node not correct: [%v]", node2.item.Id())
	} else if node3.item.Id() != 22 {
		t.Fatalf("Third node not correct: [%v]", node3.item.Id())
	} else if node4.item.Id() != 11 {
		t.Fatalf("Fourth node not correct: [%v]", node4.item.Id())
	}

	if lru.bottom != node4 {
		t.Fatalf("'bottom' node does not event fourth node.")
	}

	if node1.before != nil {
		t.Fatalf("First node has a 'before' reference.")
	} else if node4.after != nil {
		t.Fatalf("Last record has a 'after' reference.")
	}

	if node4.before != node3 {
		t.Fatalf("Fourth node 'before' node not correct.")
	}
	if node3.before != node2 {
		t.Fatalf("Third node 'before' node not correct.")
	}
	if node2.before != node1 {
		t.Fatalf("Second node 'before' node not correct.")
	}

	// Test adding one more. Now we'll be at-capacity.

	tli1 := testLruItem{
		id: 55,
	}

	added, droppedItem, err := lru.Set(tli1)
	log.PanicIf(err)

	if added != true {
		t.Fatalf("Value wasn't added but should've been.")
	} else if droppedItem != nil {
		t.Fatalf("No value should have been dropped but was: %v", droppedItem)
	}

	node55 := lru.top

	if dropped != 0 {
		t.Fatalf("No node should have been discarded.")
	}

	if len(lru.lookup) != 5 {
		t.Fatalf("LRU not the right size.")
	}

	if lru.top.item.Id() != 55 {
		t.Fatalf("Top value not updated.")
	} else if lru.top.before != nil {
		t.Fatalf("Top 'before' node not correct.")
	} else if lru.top.after != node1 {
		t.Fatalf("Top 'after' node not correct.")
	} else if lru.bottom != node4 {
		t.Fatalf("Bottom node not correct.")
	} else if lru.bottom.before != node3 {
		t.Fatalf("Bottom 'before' node not correct.")
	}

	// Cause the oldest to be discarded.

	tli2 := testLruItem{
		id: 66,
	}

	added, droppedItem, err = lru.Set(tli2)
	log.PanicIf(err)

	if added != true {
		t.Fatalf("Value wasn't added but should've been.")
	} else if droppedItem.Id() != 11 {
		t.Fatalf("Dropped value was not correct: %d", droppedItem.Id())
	}

	if dropped != 11 {
		t.Fatalf("The wrong node was discarded: (%d)", dropped)
	}

	node66 := lru.top

	if len(lru.lookup) != 5 {
		t.Fatalf("LRU not the right size.")
	}

	if lru.top.item.Id() != 66 {
		t.Fatalf("Top value not updated.")
	} else if lru.top.before != nil {
		t.Fatalf("Top 'before' node not correct.")
	} else if lru.top.after != node55 {
		t.Fatalf("Top 'after' node not correct.")
	} else if lru.bottom != node3 {
		t.Fatalf("Bottom node not correct.")
	} else if lru.bottom.before != node2 {
		t.Fatalf("Bottom 'before' node not correct.")
	}

	// Push one of the existing nodes to the top.

	dropped = 0

	tli3 := testLruItem{
		id: 33,
	}

	added, droppedItem, err = lru.Set(tli3)
	log.PanicIf(err)

	if added != false {
		t.Fatalf("Value wasn't updated but should've been.")
	} else if droppedItem != nil {
		t.Fatalf("Dropped value was not correct: %d", droppedItem.Id())
	}

	if dropped != 0 {
		t.Fatalf("No node should have been discarded.")
	}

	if len(lru.lookup) != 5 {
		t.Fatalf("LRU not the right size.")
	}

	if lru.top.item.Id() != 33 {
		t.Fatalf("Bottom node not correct")
	} else if lru.top.before != nil {
		t.Fatalf("Top 'before' node not nil.")
	} else if lru.top.after != node66 {
		t.Fatalf("Top 'before' node not nil.")
	}

	if lru.bottom.item.Id() != 22 {
		t.Fatalf("Bottom node not correct")
	}
}

func TestLru_Count(t *testing.T) {
	lru := NewLru(5)

	if lru.Count() != 0 {
		t.Fatalf("Count not correct when empty: (%d)", lru.Count())
	}

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	if lru.Count() != 1 {
		t.Fatalf("Count not correct after adding one: (%d)", lru.Count())
	}

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	if lru.Count() != 2 {
		t.Fatalf("Count not correct after adding two: (%d)", lru.Count())
	}
}

func TestLru_IsFull(t *testing.T) {
	lru := NewLru(2)

	if lru.IsFull() != false {
		t.Fatalf("IsFull not correct when empty.")
	}

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	if lru.IsFull() != false {
		t.Fatalf("IsFull not correct when with one item.")
	}

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	if lru.IsFull() != true {
		t.Fatalf("IsFull not correct with two items.")
	}
}

func TestLru_Exists(t *testing.T) {
	lru := NewLru(2)

	if lru.Exists(22) != false {
		t.Fatalf("Exists not correct when empty.")
	}

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	if lru.Exists(22) != false {
		t.Fatalf("Exists not correct when other items are present.")
	}

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	if lru.Exists(22) != true {
		t.Fatalf("Exists not correct when the key is present.")
	}
}

func TestLru_Get(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	// Check initial state.

	topId := lru.top.item.Id()
	if topId != 22 {
		t.Fatalf("Top-0 item not correct before touch: (%d)", topId)
	}

	secondId := lru.top.after.item.Id()
	if secondId != 11 {
		t.Fatalf("Top-1 item not correct before touch: (%d)", secondId)
	}

	// Try to get an invalid item.

	found, _, err := lru.Get(99)
	log.PanicIf(err)

	if found != false {
		t.Fatalf("Expected miss for unknown key.")
	}

	// Confirm that the order hasn't changed.

	topId = lru.top.item.Id()
	if topId != 22 {
		t.Fatalf("Top-0 item not correct after intentional fault: (%d)", topId)
	}

	secondId = lru.top.after.item.Id()
	if secondId != 11 {
		t.Fatalf("Top-1 item not correct after intentional fault: (%d)", secondId)
	}

	// Try to get a valid item.

	found, item, err := lru.Get(11)
	log.PanicIf(err)

	if found != true {
		t.Fatalf("Known item returned as miss (1)")
	} else if item.Id() != 11 {
		t.Fatalf("Known item return does not have right value (1): (%d)", item.Id())
	}

	// Confirm that the order *has* changed.

	topId = lru.top.item.Id()
	if topId != 11 {
		t.Fatalf("Top-0 item not correct after touch: (%d)", topId)
	}

	secondId = lru.top.after.item.Id()
	if secondId != 22 {
		t.Fatalf("Top-1 item not correct after touch: (%d)", secondId)
	}

	// Try to get another valid item (the original one). Confirm that the order
	// has returned.

	found, item, err = lru.Get(22)
	log.PanicIf(err)

	if found != true {
		t.Fatalf("Known item returned as miss (2)")
	} else if item.Id() != 22 {
		t.Fatalf("Known item return does not have right value (2): (%d)", item.Id())
	}

	// Confirm that the order *has* changed.

	topId = lru.top.item.Id()
	if topId != 22 {
		t.Fatalf("Top-0 item not correct after touch 2: (%d)", topId)
	}

	secondId = lru.top.after.item.Id()
	if secondId != 11 {
		t.Fatalf("Top-1 item not correct after touch 2: (%d)", secondId)
	}
}

func TestLru_Drop(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	if lru.Count() != 2 {
		t.Fatalf("Count before drop not correct: (%d)", lru.Count())
	}

	found, err := lru.Drop(99)
	log.PanicIf(err)

	if found != false {
		t.Fatalf("Dropping non-existent value did not report a miss.")
	}

	found, err = lru.Drop(11)
	log.PanicIf(err)

	if found != true {
		t.Fatalf("Value to drop was reported as not found.")
	}

	if lru.Count() != 1 {
		t.Fatalf("Count after drop not correct: (%d)", lru.Count())
	}
}

func TestLru_Newest(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	// Test.

	if lru.Newest() != 22 {
		t.Fatalf("Newest value not correct (1): (%d)", lru.Newest())
	}

	_, _, err = lru.Set(tli1)
	log.PanicIf(err)

	if lru.Newest() != 11 {
		t.Fatalf("Newest value not correct (2): (%d)", lru.Newest())
	}
}

func TestLru_Oldest(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	// Test.

	if lru.Oldest() != 11 {
		t.Fatalf("Oldest value not correct (1): (%d)", lru.Oldest())
	}

	_, _, err = lru.Set(tli1)
	log.PanicIf(err)

	if lru.Oldest() != 22 {
		t.Fatalf("Oldest value not correct (2): (%d)", lru.Oldest())
	}
}

func TestLru_All(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	actualRaw := lru.All()

	actual := make(sort.IntSlice, 2)
	actual[0] = actualRaw[0].(int)
	actual[1] = actualRaw[1].(int)

	actual.Sort()

	expected := sort.IntSlice{
		11,
		22,
	}

	expected.Sort()

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("All() did not return the right keys: %v != %v", actual, expected)
	}
}

func TestLru_PopOldest(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	item, err := lru.PopOldest()
	log.PanicIf(err)

	if item.Id() != 11 {
		t.Fatalf("Oldest not correct (1)")
	}

	item, err = lru.PopOldest()
	log.PanicIf(err)

	if item.Id() != 22 {
		t.Fatalf("Oldest not correct (2)")
	}

	_, err = lru.PopOldest()
	if err != ErrLruEmpty {
		t.Fatalf("Expected ErrLruEmpty for empty LRU.")
	}
}

func TestLru_FindPosition(t *testing.T) {
	lru := NewLru(2)

	// Load.

	tli1 := testLruItem{
		id: 11,
	}

	_, _, err := lru.Set(tli1)
	log.PanicIf(err)

	p := lru.FindPosition(11)
	if p != 0 {
		t.Fatalf("Position not correct (1).")
	}

	tli2 := testLruItem{
		id: 22,
	}

	_, _, err = lru.Set(tli2)
	log.PanicIf(err)

	// Was the first item pushed back by one?

	p = lru.FindPosition(11)
	if p != 1 {
		t.Fatalf("Position not correct (1).")
	}

	_, _, err = lru.Set(tli1)
	log.PanicIf(err)

	// Is it back?

	p = lru.FindPosition(11)
	if p != 0 {
		t.Fatalf("Position not correct (3).")
	}
}

func TestLru_MaxCount(t *testing.T) {
	lru := NewLru(2)

	if lru.MaxCount() != 2 {
		t.Fatalf("MaxCount not correct (1).")
	}

	lru = NewLru(55)

	if lru.MaxCount() != 55 {
		t.Fatalf("MaxCount not correct (2).")
	}
}
