package listener

type IListener interface {
}

type listener struct {
}

func NewListener() IListener {
	return &listener{}
}
