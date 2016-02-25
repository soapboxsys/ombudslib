package ombjson

type Tag struct {
	Value   string `json:"val"`
	FirstTs int64  `json:"ts"`
	Count   int64  `json:"num"`
	Score   int64  `json:"score"`
}

func NewTag(value string, cnt, ts int64) Tag {
	s := tagScore(cnt, ts)
	return Tag{
		Value:   value,
		FirstTs: ts,
		Count:   cnt,
		Score:   s,
	}
}

// tagScore uses the ratio of time passed since the peg block along with the
// passed timestamp to compute the tags 'score'
func tagScore(cnt, ts int64) int64 {
	pegStart := 1451606601.0
	var r = int64((float64(ts) - pegStart) / 86000.0)
	return cnt + r
}

type ByScore []*Tag

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
