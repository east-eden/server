package hero

import (
	"os"
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/logger"
	"github.com/east-eden/server/utils"
)

var (
	h *Hero
)

func init() {
	// snow flake init
	utils.InitMachineID(101)

	// reload to project root path
	if err := utils.RelocatePath("/server", "\\server"); err != nil {
		os.Exit(0)
	}

	// logger init
	logger.InitLogger("hero_att_test")

	excel.ReadAllEntries("config/excel/")

	heroEntry, ok := auto.GetHeroEntry(1)
	if !ok {
		os.Exit(0)
	}

	h = NewHero()
	h.Init(
		Id(1001),
		TypeId(1),
		Level(10),
		PromoteLevel(2),
		Star(2),
		Entry(heroEntry),
	)
}

func BenchmarkCalcHeroAtt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		h.GetAttManager().CalcAtt()
	}
}
