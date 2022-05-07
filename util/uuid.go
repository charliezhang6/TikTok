package util

import (
	"github.com/sony/sonyflake"
	"log"
)

func GenSonyflake() uint64 {

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	return id
}
