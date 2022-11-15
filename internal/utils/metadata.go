package utils

import (
	"time"

	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
)

func CreateMetaData(meta *setmakerpb.Metadata) *setmakerpb.Metadata {
	if meta.CreatedAt == "" {
		meta.CreatedAt = time.Now().String()
	}

	return &setmakerpb.Metadata{
		CreatedAt: meta.CreatedAt,
		UpdatedAt: time.Now().String(),
	}
}
