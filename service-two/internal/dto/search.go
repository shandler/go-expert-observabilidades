package dto

type SearchRequest struct {
	ZipCode string
}

type SearchResponse struct {
	Status int
	Body   interface{}
}
