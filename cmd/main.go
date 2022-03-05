package main

import (
	"fmt"
	"time"

	"github.com/henomis/covid-go/internal/pkg/chartjs"
	"github.com/henomis/covid-go/internal/pkg/data"
	"github.com/henomis/covid-go/internal/pkg/httpclient"
)

func main() {

	httpClient := httpclient.New(10 * time.Second)

	dataSet := data.New(httpClient)

	dataSet.ImportCSV("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv")

	dataSet.Select(data.ImportedCSVDataframeKey, []string{"data", "nuovi_positivi"})

	//fmt.Println(dataSet)

	replaceData := func(input string) string {
		return input[0:10]
	}

	dataSet.ColReplaceString(data.ImportedCSVDataframeKey, []string{"data"}, replaceData)

	//fmt.Println(dataSet)

	avg7 := dataSet.NewAVG7Dataset(data.ImportedCSVDataframeKey, "nuovi_positivi")
	//fmt.Println(avg7)

	chart := chartjs.New()
	chart.Config.Data.Labels = dataSet.DatasetAsStrings(data.ImportedCSVDataframeKey, "data")
	chart.Config.Type = chartjs.String("line")
	chart.Config.Data.Datasets = append(
		chart.Config.Data.Datasets,
		chartjs.Dataset{
			Label:       chartjs.String("positivi-avg7"),
			BorderColor: chartjs.String("#ff00ff"),
			Fill:        chartjs.False(),
			PointRadius: chartjs.Float(0.0),
			LineTension: chartjs.Float(0.4),
			Data:        avg7,
		})
	chart.Config.Data.Datasets = append(
		chart.Config.Data.Datasets,
		chartjs.Dataset{
			Label:       chartjs.String("positivi"),
			BorderColor: chartjs.String("#0000ff"),
			Fill:        chartjs.False(),
			BorderWidth: chartjs.Float(0.0),

			PointRadius: chartjs.Float(3.0),
			LineTension: chartjs.Float(0.4),
			Data:        dataSet.DatasetAsFloats(data.ImportedCSVDataframeKey, "nuovi_positivi"),
		})
	// chart.Config.Data.Datasets = append(
	// 	chart.Config.Data.Datasets,
	// 	charthtmlgo.Dataset{
	// 		Label:                "dataset2",
	// 		FillColor:            "",
	// 		StrokeColor:          "rgba(245, 15, 15, 0.5)",
	// 		PointColor:           "rgba(245, 15, 15, 0.5)",
	// 		PointStrokeColor:     "#fff",
	// 		PointHighlightFill:   "#fff",
	// 		PointHighlightStroke: "rgba(220,220,220,1)",

	// 		Data: p.Float(),
	// 	})

	fmt.Println(chart)
}
