package render

import (
	"io/ioutil"

	"github.com/henomis/covid-go/internal/pkg/chartjs"
	"github.com/henomis/covid-go/internal/pkg/data"
)

const (
	datiItaliaEndpoint        = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv"
	popolazioneItaliaEndpoint = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-statistici-riferimento/popolazione-istat-regione-range.csv"

	italiaDataframeKey              = "italia"
	popolazioneItaliaDataframeKey   = "italia_popolazione"
	italiaNuoviPositiviDataframeKey = "italia_nuovi_positivi"
)

func (r *Render) Italia() {

	r.importBasicData()
	r.nuoviPositivi()

}

func (r *Render) importBasicData() {

	err := r.dataset.ImportCSV(datiItaliaEndpoint)
	if err != nil {
		panic(err)
	}

	r.replaceData()

	r.dataset.Copy(data.ImportedCSVDataframeKey, italiaDataframeKey)

	err = r.dataset.ImportCSV(popolazioneItaliaEndpoint)
	if err != nil {
		panic(err)
	}

	r.dataset.GroupAndSum(data.ImportedCSVDataframeKey, popolazioneItaliaDataframeKey, "denominazione_regione", "totale_generale")

}

func (r *Render) replaceData() {

	replaceData := func(input string) string {
		return input[0:10]
	}

	r.dataset.ColReplaceString(data.ImportedCSVDataframeKey, []string{"data"}, replaceData)

}

func (r *Render) nuoviPositivi() {

	r.dataset.SelectAndCopy(italiaDataframeKey, italiaNuoviPositiviDataframeKey, []string{"data", "nuovi_positivi"})

	avg7 := r.dataset.NewAVG7Dataset(italiaNuoviPositiviDataframeKey, "nuovi_positivi")

	chart1 := chartjs.New()
	value := chart1.NewScatteredLineGraph(
		r.dataset.DatasetAsStrings(italiaNuoviPositiviDataframeKey, "data"),
		&chartjs.LineOptions{
			Label:   "positivi AVG7",
			Color:   "#ff0000",
			Dataset: avg7,
		},
		&chartjs.ScatterOptions{
			Label:   "positivi",
			Color:   "rgba(196,196,196,0.44)",
			Dataset: r.dataset.DatasetAsFloats(italiaNuoviPositiviDataframeKey, "nuovi_positivi"),
		},
	)

	r.dataset.Delete(italiaNuoviPositiviDataframeKey)

	ioutil.WriteFile(r.outputPath+"italia_nuovi_positivi.json", []byte(value), 0666)
}
