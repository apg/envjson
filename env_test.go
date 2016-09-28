package main

import "testing"

func TestEnvValueJSONSimple(t *testing.T) {
	// simple case. var = "default value"
	var value envValue
	payload := []byte(`"default value"`)

	err := (&value).UnmarshalJSON(payload)
	if err != nil {
		t.Fatalf("Expected %q to Unmarshal cleanly, got err=%q", err)
	}

	if value.Value != "default value" {
		t.Errorf("Expected value=%q, got value=%q", "default value", value.Value)
	}
	if value.isSet != true {
		t.Errorf("Expected isSet=true, got value=%q", value.isSet)
	}
}

func TestEnvValueJSONComplex(t *testing.T) {
	var value envValue
	payload := []byte(`
{
  "value": "default value",
  "required": true,
  "inherit": true,
  "doc": "the rain in Spain"
}`)
	err := (&value).UnmarshalJSON(payload)
	if err != nil {
		t.Fatalf("Expected %q to Unmarshal cleanly, got err=%q", err)
	}

	if value.Value != "default value" {
		t.Errorf("Expected value=%q, got value=%q", "default value", value.Value)
	}
	if value.isSet != true {
		t.Errorf("Expected isSet=true, got value=%q", value.isSet)
	}
	if value.Required != true {
		t.Errorf("Expected required=true, got required=%b", value.Required)
	}
	if value.Inherit != true {
		t.Errorf("Expected inherit=true, got inherit=%b", value.Inherit)
	}
	if value.Doc != "the rain in Spain" {
		t.Errorf("Expected doc=%q, got doc=%b", "the rain in Spain", value.Doc)
	}
}

func TestFromEnv(t *testing.T) {
	environ := []string{"SHELL=/bin/bash", "CWD=/tmp", "AGENT="}
	underTest := make(env)
	err := underTest.FromEnviron(environ)
	if err != nil {
		t.Errorf("Expected nil error, got err=%q", err)
	}

	if v, ok := underTest["SHELL"]; !ok {
		t.Errorf("Expected SHELL to be present, but wasn't")
	} else if v.Value != "/bin/bash" {
		t.Errorf("Expected SHELL=%q, got SHELL=%q instead", "/bin/bash", v.Value)
	} else if v.isSet != true {
		t.Errorf("Expected SHELL to be labeled as provided, but wasn't")
	}

	if v, ok := underTest["AGENT"]; !ok {
		t.Errorf("Expected AGENT to be present, but wasn't")
	} else if v.Value != "" {
		t.Errorf("Expected AGENT=, got AGENT=%q instead", v.Value)
	} else if v.isSet != true {
		t.Errorf("Expected AGENT to be labeled as provided, but wasn't")
	}

	if _, ok := underTest["KWYJIBO"]; ok {
		t.Errorf("Expected KWYJIBO to not be present, but was")
	}
}

func TestMerge(t *testing.T) {
	parent := make(env)
	parent["SHELL"] = envValue{Value: "/bin/sh"}

	local := make(env)
	local["SHELL"] = envValue{Required: true, Inherit: true}

	err := local.Merge(parent)
	if err != nil {
		t.Fatalf("Merge failed with error=%q", err)
	}

	if v, ok := local["SHELL"]; !ok {
		t.Error("Expected SHELL=/bin/sh, but was not even set.")
	} else if v.Value != "/bin/sh" {
		t.Error("Expected SHELL=/bin/sh, but got=%q", v.Value)
	}
}
