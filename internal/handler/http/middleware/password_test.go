package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePass(t *testing.T) {
	password := "test"

	hash, err := GeneratePass(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	match, err := ComparePass(password, hash)
	require.NoError(t, err)
	assert.True(t, match)
}

func TestComparePass(t *testing.T) {
	password := "test"

	hash, err := GeneratePass(password)
	require.NoError(t, err)

	tests := []struct {
		name      string
		pass      string
		hash      string
		wantMatch bool
		wantErr   bool
	}{
		{
			name:      "ok password",
			pass:      password,
			hash:      hash,
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:      "bad pass",
			pass:      "bad pass",
			hash:      hash,
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:      "bad hash",
			pass:      password,
			hash:      "bad hash",
			wantMatch: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComparePass(tt.pass, tt.hash)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantMatch, got)
		})
	}
}
