package main

import "testing"

func TestThatPathsAreJoinedWithASlash(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected string
	}{
		{
			a:        "/test/",
			b:        "/b/",
			expected: "/test/b/",
		},
		{
			a:        "test",
			b:        "b",
			expected: "test/b",
		},
		{
			a:        "test",
			b:        "/b",
			expected: "test/b",
		},
		{
			a:        "test/",
			b:        "b",
			expected: "test/b",
		},
	}

	for _, test := range tests {
		actual := singleJoiningSlash(test.a, test.b)
		if actual != test.expected {
			t.Errorf("for '%v' and '%v', expected '%v' got '%v'", test.a, test.b, test.expected, actual)
		}
	}
}

func TestGetKeysFromEnvironment(t *testing.T) {
	tests := []struct {
		input         []string
		expected      map[string]string
		expectedError string
	}{
		{
			input: []string{"JWTPROXY_ISSUER_0=example.com", "JWTPROXY_PUBLIC_KEY_0=dsfdsfdsfdsf"},
			expected: map[string]string{
				"example.com": "dsfdsfdsfdsf",
			},
		},
		{
			input:    []string{"unrelated=something"},
			expected: map[string]string{},
		},
		{
			input:         []string{"JWTPROXY_ISSUER_1=example.com", "JWTPROXY_PUBLIC_KEY_0=dsfdsfdsfdsf"},
			expected:      map[string]string{},
			expectedError: "could not find a matching JWTPROXY_PUBLIC_KEY_1 value for JWTPROXY_ISSUER_1",
		},
	}

	for _, test := range tests {
		actual, err := getKeysFromEnvironment(test.input)
		if err != nil && test.expectedError == "" {
			t.Error(err)
		}
		if test.expectedError != "" && err == nil {
			t.Errorf("for input '%v', expected error '%v', got nil", test.input, test.expectedError)
		}
		if !mapsAreEqual(actual, test.expected) {
			t.Errorf("for input '%v', expected '%v', got '%v'", test.input, test.expected, actual)
		}
	}
}

func mapsAreEqual(m, n map[string]string) bool {
	if len(m) != len(n) {
		return false
	}
	for k, v := range m {
		if n[k] != v {
			return false
		}
	}
	return true
}
