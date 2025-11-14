package utils_test

import (
	"testing"

	"git.riyt.dev/codeuniverse/internal/utils"
)

func TestGenerateNumericCode(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Zero length",
			length:  0,
			wantErr: true,
		},
		{
			name:    "Negative length",
			length:  -5,
			wantErr: true,
		},
		{
			name:    "Length 1",
			length:  1,
			wantErr: false,
		},
		{
			name:    "Length 7",
			length:  7,
			wantErr: false,
		},
		{
			name:    "Length 10",
			length:  10,
			wantErr: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := utils.GenerateNumericCode(c.length)

			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if codeLength := len(got); codeLength != c.length {
				t.Errorf("expected %d digits, got %d", c.length, codeLength)
			}
		})
	}
}
