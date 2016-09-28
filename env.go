package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type env map[string]envValue

type envValue struct {
	Value    string
	Required bool
	Inherit  bool
	Doc      string

	isSet bool
}

func (c *envValue) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && len(b) >= 2 {
		c.Value = string(b[1 : len(b)-1])
		c.isSet = true
		return nil
	}

	var target struct {
		Value    string `json:"value"`
		Required bool   `json:"required"`
		Inherit  bool   `json:"inherit"`
		Doc      string `json:"doc"`
	}

	err := json.Unmarshal(b, &target)
	if err != nil {
		return err
	}

	if target.Value != "" {
		c.isSet = true
	}

	c.Value = target.Value
	c.Required = target.Required
	c.Inherit = target.Inherit
	c.Doc = target.Doc

	return nil
}

// Read populates a env from a reader.
func (c env) Read(in io.Reader) error {
	dec := json.NewDecoder(in)
	err := dec.Decode(&c)

	return err
}

// FromFile reads in envValues from a file
func (c env) FromFile(f string) error {
	in, err := os.Open(f)
	if err != nil {
		return err
	}

	return c.Read(in)
}

// FromEnviron reads in envValues from the process environment
func (c env) FromEnviron(environ []string) error {
	for _, kv := range environ {
		bits := strings.SplitN(kv, "=", 2)
		if len(bits) == 2 {
			c[bits[0]] = envValue{Value: bits[1], isSet: true}
		} else {
			c[bits[0]] = envValue{Value: "", isSet: true}
		}
	}
	return nil
}

// Merge attempts to inherit any values from `p` that would satisfy required fields.
func (c env) Merge(p env) error {
	for k, v := range p { // Merge parent into c
		l, ok := c[k]
		if !ok {
			c[k] = v // just copy it over.
			continue
		}

		// Overwrite our local value with the parent's value.
		if l.Inherit {
			l.Value = v.Value
			l.isSet = true
			c[k] = l
		}
	}

	for k, v := range c {
		if v.Required && v.Value == "" {
			return fmt.Errorf("Required variable, %q, has no value", k)
		}
	}

	return nil
}

// Environ returns a []string suitable for calling os.Exec with.
func (c env) Environ() []string {
	kv := make([]string, 0, len(c))
	for k, v := range c {
		kv = append(kv, k+"="+v.Value)
	}
	return kv
}

func (c env) Dump(w io.Writer) error {
	values := make(map[string]string)
	for k, ev := range c {
		values[k] = ev.Value
	}

	enc := json.NewEncoder(w)
	return enc.Encode(values)
}
