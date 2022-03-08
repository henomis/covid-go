package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/henomis/covid-go/internal/pkg/chartjs"
	"github.com/henomis/covid-go/internal/pkg/data"
	"github.com/henomis/covid-go/internal/pkg/httpclient"
)

func main() {

	httpClient := httpclient.New(10 * time.Second)

	dataSet := data.New(httpClient)

	err := dataSet.ImportCSV("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv")
	if err != nil {
		panic(err)
	}

	dataSet.Select(data.ImportedCSVDataframeKey, []string{"data", "nuovi_positivi"})

	replaceData := func(input string) string {
		return input[0:10]
	}

	dataSet.ColReplaceString(data.ImportedCSVDataframeKey, []string{"data"}, replaceData)

	avg7 := dataSet.NewAVG7Dataset(data.ImportedCSVDataframeKey, "nuovi_positivi")

	chart1 := chartjs.New()
	value := chart1.NewScatteredLineGraph(
		dataSet.DatasetAsStrings(data.ImportedCSVDataframeKey, "data"),
		&chartjs.LineOptions{
			Label:   "positivi AVG7",
			Color:   "#ff0000",
			Dataset: avg7,
		},
		&chartjs.ScatterOptions{
			Label:   "positivi",
			Color:   "rgba(196,196,196,0.44)",
			Dataset: dataSet.DatasetAsFloats(data.ImportedCSVDataframeKey, "nuovi_positivi"),
		},
	)

	os.Mkdir("data", 0776)
	ioutil.WriteFile("./data/italia.json", []byte(value), 0666)

	fmt.Println(value)
}
