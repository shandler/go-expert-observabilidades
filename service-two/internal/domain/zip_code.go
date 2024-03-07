package domain

import (
	"context"

	"github.com/shandler/go-expert-observabilidade/service-two/internal/dto"
)

type ZipCode interface {
	Search(ctx context.Context, request dto.SearchRequest) dto.SearchResponse
}
