package utils

import (
	"time"

	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
)

func CreateMetaData(meta *domain.Metadata) *domain.Metadata {
	if meta.CreatedAt == "" {
		meta.CreatedAt = time.Now().String()
	}

	return &domain.Metadata{
		CreatedAt: meta.CreatedAt,
		UpdatedAt: time.Now().String(),
	}
}
