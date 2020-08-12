package model

import (
	"time"
)

type LyFnInputParams struct {
	Command string

	Exporter    ExporterEndpoint
	ScapeResult ScrapesEndpointResult
	Payload     RequestParam
	Gateway 	Gateway
}

// LyFnOutputParams a struct that will be returned
// The struct name need to be static, but the internal composition of the struct can be changed to fit your usage
type LyFnOutputParams struct {
	ReturnPayload interface{}
}

type ScrapesEndpointResult struct {
	ExporterID   string
	ScrapeResult string
	IsCompress bool
	ScrapeTime     time.Time
	LastUpdateTime time.Time
}

type RequestParam struct {
	ID string
}

type Gateway struct {
	ID string
	Status string
	Hostname string
}
