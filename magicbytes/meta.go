package magicbytes

type Meta struct {
	Type   string
	Bytes  []byte
	Offset int64
}

type OnMatchFunc func(path, metaType string) bool
