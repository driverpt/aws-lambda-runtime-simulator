package event

type RuntimeError struct {
	ErrorMessage string   `json:"errorMessage"`
	ErrorType    string   `json:"errorType"`
	StackTrace   []string `json:"stackTrace,omitempty"`
}
