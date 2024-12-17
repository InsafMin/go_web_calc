package calculator

import "testing"

func TestCalc(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
		wantErr    bool
	}{
		{
			name:       "Res = 4",
			expression: "1.2 + 1 * (2 + 1)",
			want:       4.2,
			wantErr:    false,
		},
		{
			name:       "Res = 10",
			expression: "2 * 3 + 4 * (2 - 1)",
			want:       10.0,
			wantErr:    false,
		},
		{
			name:       "Res = 52.8",
			expression: "2 * 3 * 4 * (2 / 1)",
			want:       48,
			wantErr:    false,
		},
		{
			name:       "Res = 3",
			expression: "2 - 3 + 4 * (2 - 1)",
			want:       3.0,
			wantErr:    false,
		},
		{
			name:       "Res = 0",
			expression: "1 - 1",
			want:       0,
			wantErr:    false,
		},
		{
			name:       "Full empty input",
			expression: "",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Empty input",
			expression: "   ",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Extra open bracket",
			expression: "1 + (3 * ()",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Extra close bracket",
			expression: "1 + 1 * (2 + 1))",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Division by zero",
			expression: "2 / 0",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Unacceptable symbol",
			expression: "& j",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Extra operator 1",
			expression: "2 / + 0",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Extra operator 2",
			expression: "1 + *",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Invalid expression",
			expression: "2 ( 8)",
			want:       0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calc(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
