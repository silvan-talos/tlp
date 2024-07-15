package tlp

type Attr struct {
	Key   string
	Value any
}

func NewAttr(name string, value any) Attr {
	return Attr{Key: name, Value: value}
}
