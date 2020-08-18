package cache

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli/v2"
)

// test object
type Object struct {
	Id      int64 `json:"_id"`
	OwnerId int64 `json:"owner_id"`
	TypeId  int32 `json:"type_id"`
	Exp     int64 `json:"exp"`
	Level   int32 `json:"level"`
}

func (o *Object) GetObjID() int64 {
	return o.Id
}

func (o *Object) GetStoreIndex() int64 {
	return o.OwnerId
}

func TestCache(t *testing.T) {
	set := flag.NewFlagSet("cache", flag.ContinueOnError)
	set.String("redis_addr", "localhost:6379", "redis address")
	ctx := cli.NewContext(nil, set, nil)
	cc := NewCache(ctx)

	o := &Object{
		Id:      1001100,
		OwnerId: 1,
		TypeId:  1001,
		Exp:     2000,
		Level:   99,
	}

	err := cc.SaveObject("test_obj", o)
	if err != nil {
		t.Fatalf("TestCache SaveObject failed: %s", err.Error())
	}

	var newObj Object
	err = cc.LoadObject("test_obj", 1001100, &newObj)
	if err != nil {
		t.Fatalf("TestCache LoadObject hit failed: %s", err.Error())
	}

	var newObj2 Object
	err = cc.LoadObject("test_obj", 20002, &newObj2)
	if err != nil && !errors.Is(err, ErrNoResult) && !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("TestCache LoadObject not hit failed: %s", err.Error())
	}

	diff := cmp.Diff(o, &newObj)
	if diff != "" {
		t.Fatalf("TestCache Compare failed: %s", diff)
	}

	var wg sync.WaitGroup
	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go performCacheLoad(t, cc, &wg)
		}
	})

	wg.Wait()
	fmt.Println("cache benchmark result: ", result.String(), result.MemString())
}

func performCacheLoad(t *testing.T, c Cache, wg *sync.WaitGroup) {
	var obj Object
	err := c.LoadObject("test_obj", 1001100, &obj)
	if err != nil && !errors.Is(err, ErrNoResult) && !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("performCacheLoad not hit: %s", err.Error())
	}

	defer wg.Done()
}
