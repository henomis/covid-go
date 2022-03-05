package chartjs

import "testing"

func TestChartjs_String(t *testing.T) {
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
						Labels: []string{"a,b,c,d,e"},
						Datasets: []Dataset{
							{
								Type:        "line",
								Label:       "dataset1",
								BorderColor: "#ff0000",
								LineTension: 0.4,
								Data:        []int{2, 5, 7, 3, 4},
							},
						},
					},
				},
			},
			want: "",
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
