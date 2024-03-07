package dto

type SearchRequest struct {
	ZipCode string `json:"zipCode"`
}

type SearchResponse struct {
	Status int
	Body   interface{}
}
