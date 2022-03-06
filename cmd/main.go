package main

import (
	"fmt"
	"time"

	"github.com/henomis/covid-go/internal/pkg/data"
	"github.com/henomis/covid-go/internal/pkg/graph"
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

	graph1 := graph.New()
	value := graph1.NewScatteredLineGraph(
		dataSet.DatasetAsStrings(data.ImportedCSVDataframeKey, "data"),
		&graph.LineOptions{
			Label:   "positivi AVG7",
			Color:   "#ff0000",
			Dataset: avg7,
		},
		&graph.ScatterOptions{
			Label:   "positivi",
			Color:   "rgba(196,196,196,0.44)",
			Dataset: dataSet.DatasetAsFloats(data.ImportedCSVDataframeKey, "nuovi_positivi"),
		},
	)

	fmt.Println(value)
}
