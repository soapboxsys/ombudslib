package ombjson

type Page struct {
	Start        string         `json:"start"`
	Stop         string         `json:"stop"`
	Bulletins    []*Bulletin    `json:"bulletins", omitempty`
	Endorsements []*Endorsement `json:"endorsements", omitempty`
}
