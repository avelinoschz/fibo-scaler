package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	tests := []struct {
		it     string
		setup  func() error
		assert func(*testing.T, error)
	}{
		{
			it: `given
				when 
				then`,
			setup: func() error {
				return errors.New("something went wrong")
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			err := tt.setup()

			tt.assert(t, err)
		})
	}
}
