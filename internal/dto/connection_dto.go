package dto

type ConnectionRequest struct {
	First  *int32  `json:"first,omitempty"`
	Last   *int32  `json:"last,omitempty"`
	Before *string `json:"before,omitempty"`
	After  *string `json:"after,omitempty"`
	Search *string `json:"search,omitempty"`
}
