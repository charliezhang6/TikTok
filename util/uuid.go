package util

import (
	"github.com/sony/sonyflake"
	"log"
)

func GenSonyflake() int64 {

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Printf("flake.NextID() failed with %s\n", err)
	}
	return int64(id)
}
