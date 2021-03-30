package gate

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"stathat.com/c/consistent"
)

var (
	cons     *consistent.Consistent = nil
	elements []string
)

func init() {
	cons = consistent.New()
	cons.NumberOfReplicas = 200
	elements = make([]string, 200)

	for n := 0; n < 200; n++ {
		elements[n] = fmt.Sprintf("elem-%d", n+1)
	}
}

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

func BenchmarkConsistentHash(b *testing.B) {
	n := rand.Int31n(200)
	cons.Set(elements[n:])
	_, err := cons.Get(fmt.Sprintf("elem-%d", n+1))
	if err != nil {
		b.Fatal(err)
	}
}
