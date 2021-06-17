package utils

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
)

func TestSonyflakeDecompose(t *testing.T) {
	m := sonyflake.Decompose(uint64(441123353153))
	highM := MachineIDHigh(int64(m["machine-id"]))
	lowM := MachineIDLow(int64(m["machine-id"]))
	log.Info().Interface("decompose_data", m).Interface("high_machine_id", highM).Interface("low_machine_id", lowM).Msg("decompose id")
}
