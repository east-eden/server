package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/yokaiio/yokai_server/game/player"
)

func init() {

}

func BenchmarkLoadFromMemory(b *testing.B) {
	b.ReportAllocs()

	ctx := context.Background()
	mongoClient, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(ctx.String("db_dsn"))); err != nil {
		logger.Fatal("NewDatastore failed:", err, "with dsn:", ctx.String("db_dsn"))
		return nil
	}

	ds.db = ds.c.Database(ctx.String("database"))
	col := m.g.ds.Database().Collection(m.TableName())
	cl := NewCacheLoader(
		context.Background(),
		m.coll,
		"_id",
		player.NewLitePlayer,
		nil,
	)

	if !ok {
		t.Error("proto assert error")
	}

	fmt.Println("assert success:", retMsg)
}
