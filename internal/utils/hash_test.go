package utils_test

import (
	"git.riyt.dev/codeuniverse/internal/utils"
	"testing"
)

func TestCreateToken(t *testing.T) {
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
			length:  -1,
			wantErr: true,
		},
		{
			name:    "Length 1",
			length:  1,
			wantErr: false,
		},
		{
			name:    "Length 32",
			length:  32,
			wantErr: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := utils.CreateToken(c.length)

			if c.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedLen := c.length * 2
			if len(got) != expectedLen {
				t.Fatalf("invalid token length: got %d, want %d", len(got), expectedLen)
			}
		})
	}
}
