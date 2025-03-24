package utils

import (
	"testing"
)

func TestCustomError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      CustomError
		expected string
	}{
		{
			name: "Error with message",
			err: CustomError{
				Message: "An error occurred",
				Code:    500,
			},
			expected: "message: An error occurred\n",
		},
		{
			name: "Error with empty message",
			err: CustomError{
				Message: "",
				Code:    400,
			},
			expected: "message: \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("CustomError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
