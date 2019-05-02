package ridata

import (
    "fmt"

    "github.com/dsoprea/go-logging"
)

type lruNode struct {
    before *lruNode
    after  *lruNode
    id     interface{}
}

func (ln *lruNode) String() string {
    var beforePhrase string
    if ln.before != nil {
        beforePhrase = fmt.Sprintf("%v", ln.before.id)
    } else {
        beforePhrase = "<NULL>"
    }

    var afterPhrase string
    if ln.after != nil {
        afterPhrase = fmt.Sprintf("%v", ln.after.id)
    } else {
        afterPhrase = "<NULL>"
    }

    return fmt.Sprintf("[%v] BEFORE=[%s] AFTER=[%s]", ln.id, beforePhrase, afterPhrase)
}

type lruEventFunc func(id interface{}) (err error)

// Lru establises an LRU of IDs of any type.
type Lru struct {
    top     *lruNode
    bottom  *lruNode
    lookup  map[interface{}]*lruNode
    maxSize int
    dropCb  lruEventFunc
}

// NewLru returns a new instance.
func NewLru(maxSize int) *Lru {
    return &Lru{
        lookup:  make(map[interface{}]*lruNode),
        maxSize: maxSize,
    }
}

// SetDropCb sets a callback that will be triggered whenever an item ages out
// or is manually dropped.
func (lru *Lru) SetDropCb(cb lruEventFunc) {
    lru.dropCb = cb
}

// Size returns the number of items in the LRU.
func (lru *Lru) Size() int {
    return len(lru.lookup)
}

// Touch bumps an item to the front of the LRU. It will be added if it doesn't
// already exist. If as a result of adding an item the LRU exceeds the maximum
// size, the least recently used item will be discarded.
func (lru *Lru) Touch(id interface{}) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    node, found := lru.lookup[id]

    if found == true {
        // It's already at the front.
        if node.before == nil {
            return nil
        }

        // Prune.
        if node.before != nil {
            node.before.after = node.after
            node.before = nil
        }

        // If we were at the bottom, the bottom is now whatever was upstream of
        // us.
        if lru.bottom == node {
            lru.bottom = lru.bottom.before
        }

        // Insert at the front.
        node.after = lru.top
    } else {
        node = &lruNode{
            id:    id,
            after: lru.top,
        }

        lru.lookup[id] = node
    }

    // Point the head of the list to us.
    lru.top = node

    // Update the link from the downstream node.
    if node.after != nil {
        node.after.before = node
    }

    if lru.bottom == nil {
        lru.bottom = node
    }

    if len(lru.lookup) > lru.maxSize {
        found, err := lru.Drop(lru.bottom.id)
        log.PanicIf(err)

        if found == false {
            log.Panicf("drop of old item was ineffectual")
        }
    }

    return nil
}

// Drop discards the given item.
func (lru *Lru) Drop(id interface{}) (found bool, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    node, found := lru.lookup[id]
    if found == false {
        return false, nil
    }

    // Keep the `top` node up-to-date.
    if node.before == nil {
        lru.top = node.after
    }

    // Keep the `bottom` node up-to-date.
    if node.after == nil {
        lru.bottom = node.before
    }

    // Detach us from the previous node and link that node to the one after us.
    if node.before != nil {
        node.before.after = node.after
    }

    delete(lru.lookup, id)

    if lru.dropCb != nil {
        err := lru.dropCb(id)
        log.PanicIf(err)
    }

    return true, nil
}

// First returns the most recently used ID.
func (lru *Lru) First() interface{} {
    if lru.top != nil {
        return lru.top.id
    } else {
        return nil
    }
}

// Last returns the least recently used ID.
func (lru *Lru) Last() interface{} {
    if lru.bottom != nil {
        return lru.bottom.id
    } else {
        return nil
    }
}

// All returns a list of all IDs.
func (lru *Lru) All() []interface{} {
    collected := make([]interface{}, len(lru.lookup))
    i := 0
    for value, _ := range lru.lookup {
        collected[i] = value
        i++
    }

    return collected
}

// All returns a list of all IDs.
func (lru *Lru) Dump() {
    fmt.Printf("Count: (%d)\n", len(lru.lookup))
    fmt.Printf("\n")

    fmt.Printf("Top: %v\n", lru.top)
    fmt.Printf("Bottom: %v\n", lru.bottom)
    fmt.Printf("\n")

    i := 0
    for ptr := lru.top; ptr != nil; ptr = ptr.after {
        fmt.Printf("%03d: %s\n", i, ptr)
        i++
    }
}
