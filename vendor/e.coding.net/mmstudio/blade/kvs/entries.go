package kvs

import "fmt"

// Entries is a list of *Entry with some helpers.
type Entries []*Entry

// Equals returns true if cmp contains the same data.
func (e Entries) Equals(cmp Entries) bool {
	// Check if the file has really changed.
	if len(e) != len(cmp) {
		return false
	}
	for i := range e {
		if !e[i].Equals(cmp[i]) {
			return false
		}
	}
	return true
}

func (e Entries) Each(handler func(stringer fmt.Stringer))  {
	for _, curr := range e {
		handler(curr)
	}
}

// Contains returns true if the Entries contain a given Entry.
func (e Entries) Contains(entry *Entry) bool {
	for _, curr := range e {
		if curr.Equals(entry) {
			return true
		}
	}
	return false
}

// Diff compares two entries and returns the added and removed entries.
func (e Entries) Diff(cmp Entries) (Entries, Entries) {
	added := Entries{}
	for _, entry := range cmp {
		if !e.Contains(entry) {
			added = append(added, entry)
		}
	}

	removed := Entries{}
	for _, entry := range e {
		if !cmp.Contains(entry) {
			removed = append(removed, entry)
		}
	}

	return added, removed
}

func CreateEntriesFromJson(addrs [][]byte) (Entries, error) {
	entries := Entries{}
	if addrs == nil {
		return entries, nil
	}

	for _, jsonb := range addrs {
		if len(jsonb) == 0 {
			continue
		}
		entry, err := NewEntryFromJson(jsonb)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
