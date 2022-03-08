package render

import (
	"os"

	"github.com/henomis/covid-go/internal/pkg/data"
	"github.com/henomis/covid-go/internal/pkg/httpclient"
)

type Render struct {
	httpClient *httpclient.HttpClient
	dataset    *data.Data
	outputPath string
}

func New(
	httpClient *httpclient.HttpClient,
	dataset *data.Data,
	outputPath string,
) *Render {

	os.Mkdir(outputPath, 0776)

	return &Render{
		httpClient: httpClient,
		dataset:    dataset,
		outputPath: outputPath,
	}
}

func (r *Render) All() {
	r.Italia()
}
