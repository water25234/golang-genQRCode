package cmd

import "testing"

func TestPinCodeInfo(t *testing.T) {
	tests := []struct {
		name         string
		input        []string
		expectedName string
		expectedPin  string
		expectErr    bool
	}{
		{
			name:         "single value",
			input:        []string{"123"},
			expectedName: "123",
			expectedPin:  "123",
			expectErr:    false,
		},
		{
			name:         "two values",
			input:        []string{"abc", "123"},
			expectedName: "abc",
			expectedPin:  "123",
			expectErr:    false,
		},
		{
			name:      "too many values",
			input:     []string{"1", "2", "3"},
			expectErr: true,
		},
	}

	fs := &flags{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name, pin, err := fs.pinCodeInfo(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if name != tc.expectedName || pin != tc.expectedPin {
				t.Errorf("expected name %s pin %s, got name %s pin %s", tc.expectedName, tc.expectedPin, name, pin)
			}
		})
	}
}
