package chartjs

import "encoding/json"

func New() *Chartjs {
	return &Chartjs{}
}

func (c *Chartjs) String() string {

	configString, err := json.Marshal(c.Config)
	if err != nil {
		return ""
	}

	return string(configString)

}

func (c *Chartjs) NewScatteredLineGraph(labels []string, lineOptions *LineOptions, scatterOptions *ScatterOptions) string {

	c.Config.Data.Labels = labels
	c.Config.Type = String("line")
	c.Config.Data.Datasets = append(
		c.Config.Data.Datasets,
		Dataset{
			Label:       String(lineOptions.Label),
			BorderColor: String(lineOptions.Color),
			Fill:        False(),
			PointRadius: Float(0.0),
			LineTension: Float(0.4),
			Data:        lineOptions.Dataset,
			Order:       Int(1),
		})
	c.Config.Data.Datasets = append(
		c.Config.Data.Datasets,
		Dataset{
			Label:       String(scatterOptions.Label),
			BorderColor: String(scatterOptions.Color),
			PointRadius: Float(3.0),
			BorderWidth: Float(1.0),
			ShowLine:    False(),
			Data:        scatterOptions.Dataset,
			Order:       Int(2),
		})

	return c.String()
}
