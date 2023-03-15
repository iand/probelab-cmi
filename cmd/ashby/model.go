package main

import (
	"context"

	grob "github.com/MetalBlueberry/go-plotly/graph_objects"
)

type PlotSpecJSON struct {
	Datasets []DataSetSpecJSON `json:"datasets,omitempty"`
	Series   []SeriesSpecJSON  `json:"series,omitempty"`
	Layout   grob.Layout       `json:"layout,omitempty"`
}

type DataSetSpecJSON struct {
	Name   string `json:"name,omitempty"`
	Source string `json:"source,omitempty"`
	Query  string `json:"query,omitempty"`
}

type SeriesSpecJSON struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"` // name of the series
	DataSet string `json:"dataset,omitempty"`
	Labels  string `json:"labels,omitempty"` // the name of the field the source should use for labels
	Values  string `json:"values,omitempty"` // the name of the field the source should use for values
}

type DataSource interface {
	GetDataSet(ctx context.Context, query string, params ...any) (DataSet, error)
}

type DataSeries struct {
	Labels []string
	Values []float64
}

type DataSet interface {
	Next() bool
	Err() error
	Field(name string) any
}
