package model

type Router interface {
	Initialize(p string) error
	GetPort() string
	Run()
	Close()
}
