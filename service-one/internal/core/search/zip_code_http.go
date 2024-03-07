package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/shandler/go-expert-observabilidade/service-one/internal/domain"
	"github.com/shandler/go-expert-observabilidade/service-one/internal/dto"
	"github.com/shandler/go-expert-observabilidade/shared"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ZipCodeHttp struct {
	client domain.HTTPClient
	config *shared.Config
}

func New(client domain.HTTPClient, config *shared.Config) *ZipCodeHttp {
	return &ZipCodeHttp{
		client: client,
		config: config,
	}
}

func (z *ZipCodeHttp) Search(ctx context.Context, request dto.SearchRequest) dto.SearchResponse {
	regex := regexp.MustCompile("^[0-9]{8}$")
	if !regex.MatchString(request.ZipCode) {
		return z.mountError(http.StatusUnprocessableEntity, "invalid zipCode")
	}

	url := fmt.Sprintf("%s?zipCode=%s", z.config.ExternalCallURL, request.ZipCode)
	req, err := http.NewRequest(z.config.ExternalCallMethod, url, nil)
	if err != nil {
		return z.mountError(http.StatusUnprocessableEntity, err.Error())
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return z.mountError(http.StatusUnprocessableEntity, err.Error())
	}

	var body map[string]interface{}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return z.mountError(res.StatusCode, "Error on decode response service two "+err.Error())
	}

	log.Println("Success response service two response ", body)

	return dto.SearchResponse{
		Status: res.StatusCode,
		Body:   body,
	}

}

func (z *ZipCodeHttp) mountError(status int, err string) dto.SearchResponse {
	return dto.SearchResponse{Status: status, Body: struct {
		Message string `json:"message"`
	}{Message: err},
	}
}
