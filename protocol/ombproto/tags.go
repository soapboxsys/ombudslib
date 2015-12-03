package ombproto

type Tag struct {
	bltn  *Bulletin
	Value string
}

func NewTag(value string, bltn *Bulletin) Tag {
	return Tag{
		bltn:  bltn,
		Value: value,
	}
}
