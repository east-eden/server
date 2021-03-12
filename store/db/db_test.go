package db

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

// test object
type Object struct {
	Id      int64 `json:"_id" bson:"_id"`
	OwnerId int64 `json:"owner_id" bson:"owner_id"`
	TypeId  int32 `json:"type_id" bson:"type_id"`
	Exp     int64 `json:"exp" bson:"exp"`
	Level   int32 `json:"level" bson:"level"`
}

func (o *Object) GetStoreIndex() int64 {
	return o.OwnerId
}

func TestDB(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := NewMockDB(ctl)

	mockDB.EXPECT().MigrateTable("test_obj", "owner_id").Return(nil).AnyTimes()

	if err := mockDB.MigrateTable("test_obj", "owner_id"); err != nil {
		t.Fatal("migrate collection test_obj failed:", err)
	}

	o := &Object{
		Id:      1001100,
		OwnerId: 1,
		TypeId:  1001,
		Exp:     2000,
		Level:   99,
	}

	mockDB.EXPECT().SaveObject("test_obj", o.Id, o).Return(nil).AnyTimes()
	err := mockDB.SaveObject("test_obj", o.Id, o)
	if err != nil {
		t.Fatalf("TestDB SaveObject failed: %s", err.Error())
	}

	var newObj Object
	newObj.Id = 1001100
	newObj.OwnerId = 1
	newObj.TypeId = 1001
	newObj.Exp = 2000
	newObj.Level = 99
	mockDB.EXPECT().LoadObject("test_obj", "_id", 1001100, &newObj).Return(nil).AnyTimes()
	err = mockDB.LoadObject("test_obj", "_id", 1001100, &newObj)
	if err != nil {
		t.Fatalf("TestDB LoadObject hit failed: %s", err.Error())
	}

	var newObj2 Object
	mockDB.EXPECT().LoadObject("test_obj", "_id", 20002, &newObj2).Return(nil).AnyTimes()
	err = mockDB.LoadObject("test_obj", "_id", 20002, &newObj2)
	if err != nil && !errors.Is(err, ErrNoResult) {
		t.Fatalf("TestDB LoadObject not hit failed: %s", err.Error())
	}

	diff := cmp.Diff(o, &newObj)
	if diff != "" {
		t.Fatalf("TestDB Compare failed: %s", diff)
	}
}

// func BenchmarkDB(b *testing.B) {
// 	set := flag.NewFlagSet("db", flag.ContinueOnError)
// 	set.String("db_dsn", "mongodb://localhost:27017", "db address")
// 	set.String("database", "unit_test", "db database")
// 	ctx := cli.NewContext(nil, set, nil)
// 	cc := NewDB(ctx)

// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			performDBLoad(b, cc)
// 		}
// 	})
// }

// func performDBLoad(b *testing.B, db DB) {
// 	var obj Object
// 	err := db.LoadObject("test_obj", "_id", 1001100, &obj)
// 	if err != nil && !errors.Is(err, ErrNoResult) {
// 		b.Fatalf("performCacheLoad not hit: %s", err.Error())
// 	}
// }
