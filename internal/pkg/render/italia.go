package render

import (
	"io/ioutil"

	"github.com/henomis/covid-go/internal/pkg/chartjs"
	"github.com/henomis/covid-go/internal/pkg/data"
)

const (
	datiItaliaEndpoint = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv"
)

func (r *Render) Italia() {

	err := r.dataset.ImportCSV(datiItaliaEndpoint)
	if err != nil {
		panic(err)
	}

	r.replaceData()

	r.nuoviPositivi()

}

func (r *Render) replaceData() {

	replaceData := func(input string) string {
		return input[0:10]
	}

	r.dataset.ColReplaceString(data.ImportedCSVDataframeKey, []string{"data"}, replaceData)

}

func (r *Render) nuoviPositivi() {

	nuoviPositiviDataframeKey := "nuoviPositivi"

	r.dataset.SelectAndCopy(data.ImportedCSVDataframeKey, nuoviPositiviDataframeKey, []string{"data", "nuovi_positivi"})

	avg7 := r.dataset.NewAVG7Dataset(nuoviPositiviDataframeKey, "nuovi_positivi")

	chart1 := chartjs.New()
	value := chart1.NewScatteredLineGraph(
		r.dataset.DatasetAsStrings(nuoviPositiviDataframeKey, "data"),
		&chartjs.LineOptions{
			Label:   "positivi AVG7",
			Color:   "#ff0000",
			Dataset: avg7,
		},
		&chartjs.ScatterOptions{
			Label:   "positivi",
			Color:   "rgba(196,196,196,0.44)",
			Dataset: r.dataset.DatasetAsFloats(nuoviPositiviDataframeKey, "nuovi_positivi"),
		},
	)

	r.dataset.Delete(nuoviPositiviDataframeKey)

	ioutil.WriteFile(r.outputPath+"italia_nuovi_positivi.json", []byte(value), 0666)
}
