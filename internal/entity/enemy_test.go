package entity

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		fn           func() Entity
		expectedType Entity
	}{
		{"new thug", newThug, &thug{}},
		{"new acolyte", newAcolyte, &acolyte{}},
		{"new assasin", newAssasin, &assasin{}},
		{"new snakes", newSnakes, &snakes{}},
	}

	for _, tt := range tests {
		got := reflect.TypeOf(tt.fn())
		want := reflect.TypeOf(tt.expectedType)
		if want != got {
			t.Errorf("incorrect type, want %q, got %q", want, got)
		}
	}
}
