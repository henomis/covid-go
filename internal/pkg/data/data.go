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

		if d.httpClient != nil {
			return fmt.Errorf("invalid httpclient")
		}

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

func (d *Data) DatasetAsInt(dataframeKey string, datasetKey string) []int {
	integers, _ := d.dataFrame[dataframeKey].Col(datasetKey).Int()
	return integers
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
