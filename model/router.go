package model

import "github.com/prometheus/prometheus/discovery/targetgroup"

type Router interface {
	Initialize(p string) error
	GetPort() string
	GetTarget() *targetgroup.Group
	Run()
	Close()
	SetMetricEndpoint()
}
