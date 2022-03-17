package render

import (
	"io/ioutil"
	"math"

	"github.com/henomis/covid-go/internal/pkg/chartjs"
	"github.com/henomis/covid-go/internal/pkg/data"
)

const (
	datiItaliaEndpoint        = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv"
	popolazioneItaliaEndpoint = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-statistici-riferimento/popolazione-istat-regione-range.csv"

	italiaDataframeKey               = "italia"
	popolazioneItaliaDataframeKey    = "italia_popolazione"
	italiaNuoviPositiviDataframeKey  = "italia_nuovi_positivi"
	italiaPositiviTotaliDataframeKey = "italia_positivi_totali"
)

func (r *Render) Italia() {

	r.importBasicData()
	r.nuoviPositivi()
	r.RT()

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

func (r *Render) RT() {

	r.dataset.SelectAndCopy(italiaDataframeKey, italiaPositiviTotaliDataframeKey, []string{"data", "totale_positivi"})

	calculateRT := func(records []float64) []float64 {

		newRecords := make([]float64, len(records))

		for i := range records {

			if i == 0 {
				newRecords[0] = 0.0
				continue
			}

			Pn := records[i]
			Pn1 := records[i-1]

			En := (Pn - Pn1) / Pn1

			Rt := math.Exp(5 * En)
			if math.IsInf(Rt, 0) || math.IsNaN(Rt) {
				Rt = 1
			}

			newRecords[i] = Rt

		}

		return newRecords

	}

	r.dataset.ColCalculateFloat(italiaPositiviTotaliDataframeKey, []string{"totale_positivi"}, calculateRT)

	avg7 := r.dataset.NewAVG7Dataset(italiaPositiviTotaliDataframeKey, "totale_positivi")

	chart1 := chartjs.New()
	value := chart1.NewScatteredLineGraph(
		r.dataset.DatasetAsStrings(italiaPositiviTotaliDataframeKey, "data"),
		&chartjs.LineOptions{
			Label:   "RT AVG7",
			Color:   "#ff0000",
			Dataset: avg7,
		},
		&chartjs.ScatterOptions{
			Label:   "RT",
			Color:   "rgba(196,196,196,0.44)",
			Dataset: r.dataset.DatasetAsFloats(italiaPositiviTotaliDataframeKey, "totale_positivi"),
		},
	)

	r.dataset.Delete(italiaPositiviTotaliDataframeKey)

	ioutil.WriteFile(r.outputPath+"italia_rt.json", []byte(value), 0666)

}
