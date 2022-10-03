package config

import "testing"

func Test_LoadConfig_InvalidPath(t *testing.T) {
	_, err := Load(LoadParam{Path: "invalid file"})
	expected := "missing safebox config file invalid file"

	if err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v", err, expected)
	}
}
