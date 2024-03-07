package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/shandler/go-expert-observabilidade/service-two/internal/domain"
	"github.com/shandler/go-expert-observabilidade/service-two/internal/dto"
	"github.com/shandler/go-expert-observabilidade/shared"
)

type ResponseSuccess struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type zipCodeHttp struct {
	client domain.HTTPClient
	config *shared.Config
}

func New(client domain.HTTPClient, config *shared.Config) *zipCodeHttp {
	return &zipCodeHttp{
		client: client,
		config: config,
	}
}

func (z *zipCodeHttp) Search(ctx context.Context, request dto.SearchRequest) dto.SearchResponse {

	regex := regexp.MustCompile("^[0-9]{8}$")
	if !regex.MatchString(request.ZipCode) {
		return z.mountError(http.StatusUnprocessableEntity, "invalid zipCode")
	}

	locale, err := z.FindZipCode(ctx, request.ZipCode)
	if err != nil {
		return z.mountError(http.StatusNotFound, err.Error())
	}

	response, err := z.FindWeather(ctx, locale)
	if err != nil {
		return z.mountError(http.StatusNotFound, err.Error())
	}

	return dto.SearchResponse{
		Status: http.StatusOK,
		Body:   response,
	}
}

func (z *zipCodeHttp) mountError(status int, err string) dto.SearchResponse {
	return dto.SearchResponse{
		Status: status,
		Body: struct {
			Message string `json:"message"`
		}{Message: err},
	}
}

func (z *zipCodeHttp) FindZipCode(ctx context.Context, zipCode string) (string, error) {
	ctx, span := z.config.OTELTracer.Start(ctx, "FindZipCode")
	defer span.End()

	urlZipCode := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode[0:5]+"-"+zipCode[5:])
	log.Println(urlZipCode)

	requestViaCep, _ := http.NewRequest(http.MethodGet, urlZipCode, nil)

	response, err := z.client.Do(requestViaCep)
	if err != nil {
		log.Println(err)
		return "", errors.New("can not found zipCode")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		log.Println(err)
		return "", errors.New("can not found zipCode")
	}

	if _, ok := body["localidade"].(string); !ok {
		log.Println("localidade not found")
		return "", errors.New("can not found zipCode")
	}

	return body["localidade"].(string), nil
}

func (z *zipCodeHttp) FindWeather(ctx context.Context, locale string) (*ResponseSuccess, error) {
	ctx, span := z.config.OTELTracer.Start(ctx, "FindWeather")
	defer span.End()

	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", z.config.WeatherKEY, url.QueryEscape(locale))
	requestWeather, _ := http.NewRequest(http.MethodGet, urlWeather, nil)

	response, err := z.client.Do(requestWeather)
	if err != nil {
		return nil, errors.New("can not found zipCode in weather api")
	}

	var body map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, errors.New("can not found zipCode in weather api")
	}

	current, ok := body["current"]
	if !ok {
		return nil, errors.New("can not found zipCode in weather api")
	}

	tempC, ok := current.(map[string]interface{})["temp_c"].(float64)
	if !ok {
		return nil, errors.New("can not found zipCode in weather api")
	}

	tempF, ok := current.(map[string]interface{})["temp_f"].(float64)
	if !ok {
		return nil, errors.New("can not found zipCode in weather api")
	}

	tempK := tempC + 273

	return &ResponseSuccess{City: locale, TempC: tempC, TempF: tempF, TempK: tempK}, nil
}
