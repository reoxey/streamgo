package core

type Service interface {
	Upload(fileName string) error
	Stream()
}
