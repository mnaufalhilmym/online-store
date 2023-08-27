package contract

type Response struct {
	Error      *Error      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Pagination struct {
	Limit *int `json:"limit,omitempty"`
	Count *int `json:"count,omitempty"`
	Page  *int `json:"page,omitempty"`
	Total *int `json:"total,omitempty"`
}
