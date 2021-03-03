package core

type port struct {
}

func (p port) Upload() {
	panic("implement me")
}

func (p port) Stream() {
	panic("implement me")
}

func NewService() Service {

	return &port{}
}
