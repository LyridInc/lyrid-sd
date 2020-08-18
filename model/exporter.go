package model

type ExporterEndpoint struct {
	ID           string
	Gateway      string
	URL          string
	ExporterType string
	AdditionalLabels map[string]string `json:"additional_labels"`
}
