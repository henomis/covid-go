package chartjs

import (
	"testing"
)

func TestChartjs_String(t *testing.T) {

	line := "line"
	dataset1 := "dataset1"
	red := "#ff0000"

	type fields struct {
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test #1",
			fields: fields{
				Config: Config{
					Data: Data{
						Labels: []string{"a", "b", "c", "d", "e"},
						Datasets: []Dataset{
							{
								Type:        &line,
								Label:       &dataset1,
								BorderColor: &red,
								Data:        []int{2, 5, 7, 3, 4},
							},
						},
					},
				},
			},
			want: "{\"data\":{\"labels\":[\"a\",\"b\",\"c\",\"d\",\"e\"],\"datasets\":[{\"type\":\"line\",\"label\":\"dataset1\",\"borderColor\":\"#ff0000\",\"data\":[2,5,7,3,4]}]}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chartjs{
				Config: tt.fields.Config,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Chartjs.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
