package cmd

import (
	"reflect"
	"testing"
)

func TestExtractTemplateVars(t *testing.T) {
	tmpl := "Hello {{NAME}}, your port is {{PORT}} and host is {{NAME}}"
	vars := extractTemplateVars(tmpl)
	want := []string{"NAME", "PORT"}
	if !reflect.DeepEqual(vars, want) {
		t.Errorf("got %v, want %v", vars, want)
	}
}

func TestExtractTemplateVarsEmpty(t *testing.T) {
	vars := extractTemplateVars("no placeholders here")
	if len(vars) != 0 {
		t.Errorf("expected empty, got %v", vars)
	}
}

func TestMissingTemplateKeys(t *testing.T) {
	tmpl := "{{HOST}}:{{PORT}}/{{PATH}}"
	envMap := map[string]string{"HOST": "localhost", "PORT": "3000"}
	missing := missingTemplateKeys(tmpl, envMap)
	if len(missing) != 1 || missing[0] != "PATH" {
		t.Errorf("expected [PATH], got %v", missing)
	}
}

func TestMissingTemplateKeysNone(t *testing.T) {
	tmpl := "{{A}} {{B}}"
	envMap := map[string]string{"A": "1", "B": "2"}
	missing := missingTemplateKeys(tmpl, envMap)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestExpandTemplate(t *testing.T) {
	tmpl := "db://{{USER}}:{{PASS}}@{{HOST}}/{{DB}}"
	envMap := map[string]string{
		"USER": "admin",
		"PASS": "secret",
		"HOST": "db.local",
		"DB":   "mydb",
	}
	result := expandTemplate(tmpl, envMap)
	want := "db://admin:secret@db.local/mydb"
	if result != want {
		t.Errorf("got %q, want %q", result, want)
	}
}
