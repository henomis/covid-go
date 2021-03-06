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

		if d.httpClient == nil {
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

func (d *Data) SelectAndCopy(dataframeKeySource, dataframeKeyDestination string, selectIndexes SelectIndexes) {
	d.dataFrame[dataframeKeyDestination] = d.dataFrame[dataframeKeySource].Select(selectIndexes)
}

func (d *Data) Copy(dataframeKeySource, dataframeKeyDestination string) {
	d.dataFrame[dataframeKeyDestination] = d.dataFrame[dataframeKeySource]
}

func (d *Data) Delete(dataframeKey string) {
	delete(d.dataFrame, dataframeKey)
}

func (d *Data) Print(dataframeKey string) string {

	return d.dataFrame[dataframeKey].String()
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

func (d *Data) ColReplaceFloat(dataframeKey string, columnNames []string, replaceFunc func(index int, input float64) float64) {

	d.dataFrame[dataframeKey] = d.dataFrame[dataframeKey].Capply(func(s series.Series) series.Series {
		if !keyIsIn(s.Name, columnNames) {
			return s
		}

		newRecords := []float64{}

		for i, v := range s.Float() {
			newRecords = append(newRecords, replaceFunc(i, v))
		}

		return series.Floats(newRecords)
	})
}

func (d *Data) ColCalculateFloat(dataframeKey string, columnNames []string, replaceFunc func(records []float64) []float64) {

	d.dataFrame[dataframeKey] = d.dataFrame[dataframeKey].Capply(func(s series.Series) series.Series {
		if !keyIsIn(s.Name, columnNames) {
			return s
		}

		newRecords := replaceFunc(s.Float())

		return series.Floats(newRecords)
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

func (d *Data) GroupAndSum(dataframeKeySource, dataframeKeyDestination, groupingColumnName, summingColumnName string) {
	d.dataFrame[dataframeKeyDestination] = d.dataFrame[dataframeKeySource].GroupBy(groupingColumnName).Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{summingColumnName})
}

func keyIsIn(key string, keys []string) bool {
	for _, v := range keys {
		if v == key {
			return true
		}
	}
	return false
}
