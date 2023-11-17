package type_common

type DatabaseError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
