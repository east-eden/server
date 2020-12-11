package entries

import (
	"testing"

	"github.com/east-eden/server/utils"
	"github.com/google/go-cmp/cmp"
)

// test cases
var (
	cases = map[string]struct {
		id       int32
		wantName string
	}{
		"hero_entry": {
			id:       1,
			wantName: "梅林",
		},
	}
)

func TestEntries(t *testing.T) {
	// relocate path
	if err := utils.RelocatePath(); err != nil {
		t.Fatalf("TestEntries failed: %s", err.Error())
	}

	InitEntries()

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			entry := GetHeroEntry(cs.id)
			if entry == nil {
				t.Fatalf("get hero entry<%d> failed", cs.id)
			}

			diff := cmp.Diff(entry.Name, cs.wantName)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
