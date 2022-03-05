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
