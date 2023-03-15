package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	grob "github.com/MetalBlueberry/go-plotly/graph_objects"
	"github.com/MetalBlueberry/go-plotly/offline"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var plotCommand = &cli.Command{
	Name:   "plot",
	Usage:  "Interactive command to generate a single plot",
	Action: Plot,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:        "preview",
			Required:    false,
			Usage:       "Preview the plot in a browser window.",
			Destination: &plotOpts.preview,
		},
		&cli.BoolFlag{
			Name:        "compact",
			Required:    false,
			Usage:       "Emit compact json instead of pretty-printed.",
			Destination: &plotOpts.compact,
		},
		&cli.StringSliceFlag{
			Name:        "source",
			Aliases:     []string{"s"},
			Required:    false,
			Usage:       "Specify the url of a data source, in the format name=url. May be repeated to specify multiple sources. Postgres urls take the form 'postgres://username:password@hostname:5432/database_name'",
			Destination: &plotOpts.sources,
		},
	}, loggingFlags...),
}

var plotOpts struct {
	preview bool
	compact bool
	sources cli.StringSlice
}

func Plot(cc *cli.Context) error {
	ctx := cc.Context

	sources := map[string]DataSource{
		"demo": &DemoDataSource{},
	}

	for _, sopt := range plotOpts.sources.Value() {
		name, url, ok := strings.Cut(sopt, "=")
		if !ok {
			return fmt.Errorf("source option not valid, use format 'name=url'")
		}

		if _, exists := sources[name]; exists {
			return fmt.Errorf("duplicate source %q specified", name)
		}

		if strings.HasPrefix(url, "postgres:") {
			sources[name] = NewPgDataSource(url)
		} else {
			return fmt.Errorf("unsupported source url: %s", url)
		}

	}

	if cc.NArg() != 1 {
		return fmt.Errorf("plot specification must be supplied (examples in plots folder)")
	}

	fname := cc.Args().Get(0)

	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("failed to open plot specification: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var plotSpec PlotSpecJSON
	if err := dec.Decode(&plotSpec); err != nil {
		return fmt.Errorf("failed to unmarshal plot definition: %w", err)
	}

	fig := &grob.Fig{
		Layout: &plotSpec.Layout,
	}

	dataSets := make(map[string]DataSet)
	for _, ds := range plotSpec.Datasets {
		src, exists := sources[ds.Source]
		if !exists {
			return fmt.Errorf("unknown dataset source: %s", ds.Source)
		}
		var err error
		dataSets[ds.Name], err = src.GetDataSet(ctx, ds.Query)
		if err != nil {
			return fmt.Errorf("failed to get dataset from source %q: %w", ds.Source, err)
		}
	}

	seriesByDataSet := make(map[string][]SeriesSpecJSON)
	for i, s := range plotSpec.Series {
		if _, ok := dataSets[s.DataSet]; !ok {
			slog.Error(fmt.Sprintf("unknown dataset name %q in series %d", s.DataSet, i))
			continue
		}
		seriesByDataSet[s.DataSet] = append(seriesByDataSet[s.DataSet], s)
	}

	for dsname, series := range seriesByDataSet {
		ds := dataSets[dsname]

		data := make([][2][]any, len(series))

		for ds.Next() {
			for i, s := range series {
				data[i][0] = append(data[i][0], ds.Field(s.Labels))
				data[i][1] = append(data[i][1], ds.Field(s.Values))
			}
		}

		if ds.Err() != nil {
			return fmt.Errorf("dataset iteration ended with an error: %w", ds.Err())
		}

		fig.Data = grob.Traces{}

		for i, s := range series {
			switch s.Type {
			case "bar":
				trace := &grob.Bar{
					Type: grob.TraceTypeBar,
					Name: s.Name,
					X:    data[i][0],
					Y:    data[i][1],
				}

				if trace.Name == "" {
					trace.Name = fmt.Sprintf("series %d", i)
				}
				fig.Data = append(fig.Data, trace)
			}
		}

	}

	var data []byte
	if plotOpts.compact {
		data, err = json.Marshal(fig)
	} else {
		data, err = json.MarshalIndent(fig, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("failed to marshal to json: %w", err)
	}

	fmt.Println(string(data))

	if plotOpts.preview {
		offline.Show(fig)
	}
	return nil
}

type DemoDataSource struct{}

func (s *DemoDataSource) GetDataSet(_ context.Context, query string, params ...any) (DataSet, error) {
	switch query {
	case "populations":
		return &StaticDataSet{Data: map[string][]any{
			"creature": {"giraffes", "orangutans", "monkeys"},
			"month1":   {20, 14, 23},
			"month2":   {2, 18, 29},
		}}, nil
	default:
		return nil, fmt.Errorf("unknown demo dataset: %s", query)
	}
}
