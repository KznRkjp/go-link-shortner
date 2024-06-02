package models

// const (
// 	TypeSimpleUtterance = "SimpleUtterance"
// )

// Request описывает запрос пользователя.
// см. https://yandex.ru/dev/dialogs/alice/doc/request.html
type Request struct {
	URL string `json:"url"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
	ShortURL      string
}

// // SimpleUtterance описывает команду, полученную в запросе типа SimpleUtterance.
// type SimpleUtterance struct {
// 	Type    string `json:"type"`
// 	Command string `json:"command"`
// }

// Response описывает ответ сервера.
// см. https://yandex.ru/dev/dialogs/alice/doc/response.html
// type Response struct {
// 	Response ResponsePayload `json:"response"`
// 	Version  string          `json:"version"`
// }

// ResponsePayload описывает ответ, который нужно озвучить.
type Response struct {
	Result string `json:"result"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"short_url"`
}

type URLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
