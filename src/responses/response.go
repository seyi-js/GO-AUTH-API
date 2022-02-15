package responses

type Response struct {
	Status int                    `json:"status"`
	Info   map[string]interface{} `json:"info"`
}
