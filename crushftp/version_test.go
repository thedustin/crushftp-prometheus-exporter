package crushftp

import "testing"

func TestVersionExtraction(t *testing.T) {
	tests := map[string]string{
		"Version 1.2.3": "1.2.3",
	}

	for fullStr, expectedStr := range tests {
		cleanStr := extractVersionFromString(fullStr)

		if expectedStr != cleanStr {
			t.Errorf("extraction failed on %q", fullStr)
			continue
		}
	}
}

func TestVersionParse(t *testing.T) {
	tests := map[string]*version{
		"1.2.3": {Major: 1, Minor: 2, Patch: 3},
	}

	for versionStr, expected := range tests {
		v := version{}
		v.Parse(versionStr)

		if expected.Major != v.Major ||
			expected.Minor != v.Minor ||
			expected.Patch != v.Patch {
			t.Errorf("version mismatch on %q", versionStr)
			continue
		}
	}
}

func TestVersionStringify(t *testing.T) {
	tests := map[*version]string{
		{Major: 1, Minor: 2, Patch: 3}: "1.2.3",
	}

	for v, expectedStr := range tests {
		if v.String() != expectedStr {
			t.Errorf("version %q was stringified to %q", expectedStr, v)
			continue
		}
	}
}
