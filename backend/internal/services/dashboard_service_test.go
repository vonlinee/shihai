package services

import "testing"

func TestCalculateTrendPercent(t *testing.T) {
	tests := []struct {
		name     string
		current  int64
		previous int64
		want     int64
	}{
		{name: "both zero", current: 0, previous: 0, want: 0},
		{name: "new growth from zero", current: 3, previous: 0, want: 100},
		{name: "increase", current: 15, previous: 10, want: 50},
		{name: "decrease", current: 7, previous: 10, want: -30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTrendPercent(tt.current, tt.previous); got != tt.want {
				t.Fatalf("calculateTrendPercent(%d, %d) = %d, want %d", tt.current, tt.previous, got, tt.want)
			}
		})
	}
}
