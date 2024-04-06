package listener

type IListener interface {
	Listen() bool
}

type listener struct {
}

func NewListener() IListener {
	return &listener{}
}

func (l *listener) Listen() bool {
	return true
}
