package counter

type Counter interface {
	Inc(url string)
	Count() (count int)
	Shutdown() (err error)
}

type JSONElement struct {
	URL    string `json:"url"`
	Expire int64  `json:"time"`
}
