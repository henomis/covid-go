package data

import (
	"fmt"
	"io"
	"net/url"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/henomis/covid-go/internal/pkg/httpclient"
)

const (
	ImportedCSVDataframeKey = "import-csv"
)

type ReplaceStringFunction func(string) string

type Data struct {
	httpClient *httpclient.HttpClient
	dataFrame  map[string]dataframe.DataFrame
}

type SelectIndexes interface{}

func New(httpClient *httpclient.HttpClient) *Data {
	return &Data{
		httpClient: httpClient,
		dataFrame:  make(map[string]dataframe.DataFrame),
	}
}

func (d *Data) ImportCSV(endpoint string) error {

	//TODO insert parse and distinction
	// between url and file://

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	var bodyStream io.ReadCloser

	switch endpointURL.Scheme {
	case "http", "https":
		bodyStream, err = d.httpClient.Download(endpoint)
		defer bodyStream.Close()
		if err != nil {
			return err
		}

		d.dataFrame[ImportedCSVDataframeKey] = dataframe.ReadCSV(bodyStream)

	}

	return nil

}

func (d *Data) Select(dataframeKey string, selectIndexes SelectIndexes) {
	d.dataFrame[dataframeKey] = d.dataFrame[dataframeKey].Select(selectIndexes)
}

func (d *Data) String() string {

	output := ""

	for name, dataframe := range d.dataFrame {
		output += fmt.Sprintf("Key: '%s'\n%s", name, dataframe.String())
	}

	return output
}

func (d *Data) ColReplaceString(dataframeKey string, columnNames []string, replaceFunc func(input string) string) {

	d.dataFrame[dataframeKey] = d.dataFrame[dataframeKey].Capply(func(s series.Series) series.Series {
		if !keyIsIn(s.Name, columnNames) {
			return s
		}

		newRecords := []string{}

		for _, v := range s.Records() {
			newRecords = append(newRecords, replaceFunc(v))
		}

		return series.Strings(newRecords)
	})
}

func (d *Data) DatasetAsStrings(dataframeKey string, datasetKey string) []string {
	return d.dataFrame[dataframeKey].Col(datasetKey).Records()
}

func (d *Data) DatasetAsFloats(dataframeKey string, datasetKey string) []float64 {
	return d.dataFrame[dataframeKey].Col(datasetKey).Float()
}

func (d *Data) NewAVG7Dataset(dataframeKey, datasetKey string) []float64 {
	p := d.dataFrame[dataframeKey].Col(datasetKey).Rolling(7).Mean()

	nans := p.IsNaN()
	for i, v := range nans {
		if v {
			p.Elem(i).Set(0.0)
		}
	}

	return p.Float()

}

func keyIsIn(key string, keys []string) bool {
	for _, v := range keys {
		if v == key {
			return true
		}
	}
	return false
}

// func main() {

// 	httpStream, err := downloadFile("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv")
// 	if err != nil {
// 		panic(err)
// 	}

// 	df := dataframe.ReadCSV(httpStream)
// 	df = df.Select([]string{"data", "nuovi_positivi"})

// 	removeData := func(s series.Series) series.Series {

// 		if s.Name != "data" {
// 			return s
// 		}

// 		dates := s.Records()
// 		newDates := []string{}

// 		for _, v := range dates {
// 			newDates = append(newDates, v[0:10])
// 		}
// 		return series.Strings(newDates)
// 	}

// 	df = df.Capply(removeData)

// 	p := df.Col("nuovi_positivi").Rolling(7).Mean()

// 	nans := p.IsNaN()
// 	for i, v := range nans {
// 		if v {
// 			p.Elem(i).Set(0.0)
// 		}
// 	}

// 	np, _ := df.Col("nuovi_positivi").Int()

// 	chart := charthtmlgo.New()
// 	chart.Config.Type = "line"
// 	chart.Config.Data.Labels = df.Col("data").Records()
// 	chart.Config.Data.Datasets = append(
// 		chart.Config.Data.Datasets,
// 		charthtmlgo.Dataset{
// 			Label:                "dataset1",
// 			FillColor:            "rgba(151,187,205,0.2)",
// 			StrokeColor:          "rgba(151,187,205,1)",
// 			PointColor:           "rgba(151,187,205,1)",
// 			PointStrokeColor:     "#fff",
// 			PointHighlightFill:   "#fff",
// 			PointHighlightStroke: "rgba(151,187,205,1)",
// 			Data:                 np,
// 		})
// 	chart.Config.Data.Datasets = append(
// 		chart.Config.Data.Datasets,
// 		charthtmlgo.Dataset{
// 			Label:                "dataset2",
// 			FillColor:            "",
// 			StrokeColor:          "rgba(245, 15, 15, 0.5)",
// 			PointColor:           "rgba(245, 15, 15, 0.5)",
// 			PointStrokeColor:     "#fff",
// 			PointHighlightFill:   "#fff",
// 			PointHighlightStroke: "rgba(220,220,220,1)",

// 			Data: p.Float(),
// 		})

// 	fmt.Println(chart)

// }
