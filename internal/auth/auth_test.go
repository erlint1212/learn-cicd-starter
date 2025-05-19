package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Define the expected malformed error for comparison by message later
	expectedMalformedError := errors.New("malformed authorization header")

	testCases := []struct {
		name          string
		headers       http.Header
		expectedKey   string
		expectedError error
	}{
		{
			name:          "No Authorization Header",
			headers:       http.Header{},
			expectedKey:   "",
			expectedError: ErrNoAuthHeaderIncluded, // This is a named error, use errors.Is
		},
		{
			name: "Malformed Header - Wrong Prefix",
			headers: http.Header{
				"Authorization": []string{"Bearer someapikey"},
			},
			expectedKey:   "",
			expectedError: expectedMalformedError, // Expect the malformed error message
		},
		{
			name: "Malformed Header - Missing Key",
			headers: http.Header{
				"Authorization": []string{"ApiKey"},
			},
			expectedKey:   "",
			expectedError: expectedMalformedError, // Expect the malformed error message
		},
		{
			name: "Valid API Key Header",
			headers: http.Header{
				"Authorization": []string{"ApiKey validkey123"},
			},
			expectedKey:   "validkey123",
			expectedError: nil,
		},
		{
			// Based on the original auth.go code's behavior,
			// this header with leading/trailing spaces is considered malformed.
			name: "Header with leading/trailing spaces (Malformed in Original Code)",
			headers: http.Header{
				"Authorization": []string{"  ApiKey   anotherkey456  "},
			},
			expectedKey:   "",                     // Original code returns empty key
			expectedError: expectedMalformedError, // Original code returns malformed error
		},
		{
			// Based on the original auth.go code's behavior,
			// this header with extra spaces between prefix and key
			// results in an empty key with no error.
			name: "Header with extra spaces (Empty Key, No Error in Original Code)",
			headers: http.Header{
				"Authorization": []string{"ApiKey    spacedkey789"},
			},
			expectedKey:   "",  // Original code returns empty key
			expectedError: nil, // Original code returns no error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiKey, err := GetAPIKey(tc.headers)

			// --- Error checking ---
			if tc.expectedError != nil {
				// We expect an error
				if err == nil {
					t.Fatalf("Expected error: %v, Got nil", tc.expectedError)
				}

				// Check if the received error matches the expected one
				if tc.expectedError == ErrNoAuthHeaderIncluded {
					// For the named ErrNoAuthHeaderIncluded, use errors.Is
					if !errors.Is(err, tc.expectedError) {
						t.Errorf("Expected error: %v, Got error: %v", tc.expectedError, err)
					}
				} else if tc.expectedError.Error() == expectedMalformedError.Error() {
					// For the malformed error (created by errors.New), compare error messages
					if err.Error() != tc.expectedError.Error() {
						t.Errorf("Expected error message: %q, Got error message: %q", tc.expectedError.Error(), err.Error())
					}
				} else {
					// Handle unexpected expectedError values in the test case definition
					t.Fatalf("Test case error: Unhandled expected error type: %v", tc.expectedError)
				}

			} else {
				// We expect no error
				if err != nil {
					t.Fatalf("Expected no error, Got error: %v", err)
				}
			}

			// --- Key checking ---
			// Check the returned key only if the expected error was nil (meaning the function should have successfully parsed a key)
			// or if the test specifically expects a non-empty key.
			// Given the original function's behavior, a non-empty key is only returned on success (err == nil).
			if apiKey != tc.expectedKey {
				t.Errorf("Expected key: %q, Got key: %q", tc.expectedKey, apiKey)
			}

		})
	}
}
