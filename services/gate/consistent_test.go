package gate

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"stathat.com/c/consistent"
)

func TestConsistent(t *testing.T) {
	cons := consistent.New()
	cons.NumberOfReplicas = 200

	var nodeNum int = 10
	var userNum int = 10000

	nodeNames := []string{}
	nodeHitsInit := make(map[string]int)
	for n := 0; n < nodeNum; n++ {
		name := fmt.Sprintf("game-%d", n+1)
		nodeNames = append(nodeNames, name)
	}

	cons.Set(nodeNames)

	for n := 0; n < userNum; n++ {
		userId := fmt.Sprintf("user-%d", n+1)
		nodeName, err := cons.Get(userId)
		if err != nil {
			t.Fatalf("get node name failed: %s", err.Error())
		}

		nodeHitsInit[nodeName]++
	}

	for k, v := range nodeHitsInit {
		fmt.Println("node_name: ", k, ", hit number: ", v)
	}

	fmt.Println("remove node game-5...")
	cons.Remove("game-5")

	nodeHitsRemove := make(map[string]int)
	for n := 0; n < userNum; n++ {
		userId := fmt.Sprintf("user-%d", n+1)
		nodeName, err := cons.Get(userId)
		if err != nil {
			t.Fatalf("get node name failed: %s", err.Error())
		}

		nodeHitsRemove[nodeName]++
	}

	for k, v := range nodeHitsRemove {
		fmt.Println("node_name: ", k, ", hit number: ", v)
	}

	fmt.Println("add node game-5...")
	cons.Add("game-5")

	nodeHitsAdd := make(map[string]int)
	for n := 0; n < userNum; n++ {
		userId := fmt.Sprintf("user-%d", n+1)
		nodeName, err := cons.Get(userId)
		if err != nil {
			t.Fatalf("get node name failed: %s", err.Error())
		}

		nodeHitsAdd[nodeName]++
	}

	for k, v := range nodeHitsAdd {
		fmt.Println("node_name: ", k, ", hit number: ", v)
	}

	diff := cmp.Diff(nodeHitsInit, nodeHitsAdd)
	if diff != "" {
		t.Fatalf("consistent hash result different: %s", diff)
	}
}
