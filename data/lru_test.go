package ridata

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func CreateTestLru(dropped *int) *Lru {
	droppedCb := func(raw interface{}) (err error) {
		id := raw.(int)
		*dropped = id

		return nil
	}

	lru := NewLru(5)
	lru.SetDropCb(droppedCb)

	err := lru.Touch(11)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	err = lru.Touch(22)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	err = lru.Touch(33)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	err = lru.Touch(44)
	log.PanicIf(err)

	if *dropped != 0 {
		log.Panicf("No node should have been discarded.")
	}

	size := lru.Size()
	if size != 4 {
		log.Panicf("LRU size not correct: (%d)", size)
	}

	return lru
}

func TestLru_Touch(t *testing.T) {
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

	if node1.id != 44 {
		t.Fatalf("First node not correct: [%v]", node1.id)
	} else if node2.id != 33 {
		t.Fatalf("Second node not correct: [%v]", node2.id)
	} else if node3.id != 22 {
		t.Fatalf("Third node not correct: [%v]", node3.id)
	} else if node4.id != 11 {
		t.Fatalf("Fourth node not correct: [%v]", node4.id)
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

	err := lru.Touch(55)
	log.PanicIf(err)

	node55 := lru.top

	if dropped != 0 {
		t.Fatalf("No node should have been discarded.")
	}

	if len(lru.lookup) != 5 {
		t.Fatalf("LRU not the right size.")
	}

	if lru.top.id != 55 {
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

	err = lru.Touch(66)
	log.PanicIf(err)

	if dropped != 11 {
		t.Fatalf("The wrong node was discarded: (%d)", dropped)
	}

	node66 := lru.top

	if len(lru.lookup) != 5 {
		t.Fatalf("LRU not the right size.")
	}

	if lru.top.id != 66 {
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

	err = lru.Touch(33)
	log.PanicIf(err)

	if dropped != 0 {
		t.Fatalf("No node should have been discarded.")
	}

	if len(lru.lookup) != 5 {
		t.Fatalf("LRU not the right size.")
	}

	if lru.top.id != 33 {
		t.Fatalf("Bottom node not correct")
	} else if lru.top.before != nil {
		t.Fatalf("Top 'before' node not nil.")
	} else if lru.top.after != node66 {
		t.Fatalf("Top 'before' node not nil.")
	}

	if lru.bottom.id != 22 {
		t.Fatalf("Bottom node not correct")
	}
}
