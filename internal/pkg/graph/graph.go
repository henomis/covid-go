package graph

import (
	"github.com/henomis/covid-go/internal/pkg/chartjs"
)

type Graph struct{}

type LineOptions struct {
	Dataset interface{}
	Label   string
	Color   string
}

type ScatterOptions struct {
	Dataset interface{}
	Label   string
	Color   string
}

func New() *Graph {
	return &Graph{}
}

func (g *Graph) NewScatteredLineGraph(labels []string, lineOptions *LineOptions, scatterOptions *ScatterOptions) string {

	chart := chartjs.New()
	chart.Config.Data.Labels = labels
	chart.Config.Type = chartjs.String("line")
	chart.Config.Data.Datasets = append(
		chart.Config.Data.Datasets,
		chartjs.Dataset{
			Label:       chartjs.String(lineOptions.Label),
			BorderColor: chartjs.String(lineOptions.Color),
			Fill:        chartjs.False(),
			PointRadius: chartjs.Float(0.0),
			LineTension: chartjs.Float(0.4),
			Data:        lineOptions.Dataset,
			Order:       chartjs.Int(1),
		})
	chart.Config.Data.Datasets = append(
		chart.Config.Data.Datasets,
		chartjs.Dataset{
			Label:       chartjs.String(scatterOptions.Label),
			BorderColor: chartjs.String(scatterOptions.Color),
			PointRadius: chartjs.Float(3.0),
			BorderWidth: chartjs.Float(1.0),
			ShowLine:    chartjs.False(),
			Data:        scatterOptions.Dataset,
			Order:       chartjs.Int(2),
		})

	return chart.String()
}
