package main

import (
	"time"

	"github.com/henomis/covid-go/internal/pkg/data"
	"github.com/henomis/covid-go/internal/pkg/httpclient"
	"github.com/henomis/covid-go/internal/pkg/render"
)

func main() {

	httpClient := httpclient.New(10 * time.Second)

	dataSet := data.New(httpClient)

	jsonRender := render.New(httpClient, dataSet, "./data/")

	jsonRender.All()

}
