package chartjs

type Chartjs struct {
	Config Config
}

type Config struct {
	Type    *string  `json:"type,omitempty"`
	Data    Data     `json:"data,omitempty"`
	Options *Options `json:"options,omitempty"`
	Plugins *Plugins `json:"plugins,omitempty"`
}

type Data struct {
	Labels   []string  `json:"labels,omitempty"`
	Datasets []Dataset `json:"datasets,omitempty"`
}

type Dataset struct {
	Type            *string  `json:"type,omitempty"`
	Label           *string  `json:"label,omitempty"`
	BackgroundColor *string  `json:"backgroundColor,omitempty"`
	BorderColor     *string  `json:"borderColor,omitempty"`
	BorderWidth     *float64 `json:"borderWidth,omitempty"`
	Fill            *bool    `json:"fill,omitempty"`

	SteppedLine            *bool    `json:"steppedLine,omitempty"`
	LineTension            *float64 `json:"lineTension,omitempty"`
	CubicInterpolationMode *string  `json:"cubicInterpolationMode,omitempty"`
	PointBackgroundColor   *string  `json:"pointBackgroundColor,omitempty"`
	PointBorderColor       *string  `json:"pointBorderColor,omitempty"`
	PointBorderWidth       *float64 `json:"pointBorderWidth,omitempty"`
	PointRadius            *float64 `json:"pointRadius,omitempty"`
	PointHitRadius         *float64 `json:"pointHitRadius,omitempty"`
	PointHoverRadius       *float64 `json:"pointHoverRadius,omitempty"`
	PointHoverBorderColor  *string  `json:"pointHoverBorderColor,omitempty"`
	PointHoverBorderWidth  *float64 `json:"pointHoverBorderWidth,omitempty"`
	PointStyle             *int     `json:"pointStyle,omitempty"`
	ShowLine               *bool    `json:"showLine,omitempty"`
	Order                  *int64   `json:"order,omitempty"`

	// Axis ID that matches the ID on the Axis where this dataset is to be drawn.
	XAxisID *string `json:"xAxisID,omitempty"`
	YAxisID *string `json:"yAxisID,omitempty"`

	// set the formatter for the data, e.g. "%.2f"
	// these are not exported in the json, just used to determine the decimals of precision to show
	XFloatFormat string `json:"-"`
	YFloatFormat string `json:"-"`

	Data interface{} `json:"data,omitempty"`
}

type Options struct{}

type Plugins struct {
	Legend *PluginLegend `json:"legend,omitempty"`
	Title  *PluginTitle  `json:"title,omitempty"`
}

type PluginLegend struct {
	Position string `json:"position,omitempty"`
}

type PluginTitle struct {
	Display *bool  `json:"display,omitempty"`
	Text    string `json:"text,omitempty"`
}

func True() *bool {
	t := true
	return &t
}

func False() *bool {
	t := false
	return &t
}

func Float(input float64) *float64 {
	t := input
	return &t
}

func Int(input int64) *int64 {
	t := input
	return &t
}

func String(input string) *string {
	t := input
	return &t
}
