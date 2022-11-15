package utils

import (
	"time"

	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
)

func SetMetaData(meta *setmakerpb.Metadata) {
	if meta.CreatedAt == "" {
		meta.CreatedAt = time.Now().String()
	}

	meta.UpdatedAt = time.Now().String()
}
