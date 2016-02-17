package ombjson

type BltnPage struct {
	Start     string      `json:"start"`
	Stop      string      `json:"stop"`
	Bulletins []*Bulletin `json:"bulletins"`
}

type Page struct {
	Start        string         `json:"start"`
	Stop         string         `json:"stop"`
	Bulletins    []*Bulletin    `json:"bulletins"`
	Endorsements []*Endorsement `json:"endorsements"`
}
