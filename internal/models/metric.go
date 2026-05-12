package models

type Metric struct {
	Name string 
	Value float64
	Type string
	Labels map[string]string
	Timestamp int64
}