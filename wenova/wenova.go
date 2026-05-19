package wenova

import "strings"

const DefaultBaseURL = "https://apimicroservices.wenova.fun"

type Wenova struct {
	BaseUrl string
	Token   string
}

func NewWenovaAPI(token string) Wenova {
	return Wenova{
		BaseUrl: DefaultBaseURL,
		Token:   strings.TrimSpace(token),
	}
}

func (w *Wenova) SysChangeBaseUrl(baseURL string) {
	w.BaseUrl = strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
}
