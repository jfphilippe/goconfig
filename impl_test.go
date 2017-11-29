/*
Copyright (c) 2017 Jean-François PHILIPPE
Package goconfig read config files.
*/

package goconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Check GetString.
func TestGetString0(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\"}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("key")
	if nil != serr {
		t.Error("Key 'key' not found", serr)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

	// Search a Bool key as String
	str, serr = config.GetString("nope")
	if nil != serr {
		t.Error("Key 'nope' not found", serr)
	}
	if "true" != str {
		t.Error("Wrong value found :", str)
	}

	str, serr = config.GetString("missing")
	if nil == serr {
		t.Error("Key 'missing' found")
	}

	// Missing value with a default value
	str, serr = config.GetString("missing", "deflt")
	if "deflt" != str {
		t.Error("Wrong value found :", str)
	}

	// Existing value with a default value
	str, serr = config.GetString("key", "deflt")
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check GetString. for nested string
func TestGetString1(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"value\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

	// Search a sub key in a non existant sub item
	str, serr = config.GetString("key.sub")
	if nil == serr {
		t.Error("Key 'key.sub' found")
	}

	// Missing value with a default value
	str, serr = config.GetString("nope.sub", "deflt")
	if "deflt" != str {
		t.Error("Wrong value found :", str)
	}

	// Existing value with a default value
	str, serr = config.GetString("sub.nope.key", "deflt")
	if "deflt" != str {
		t.Error("Wrong value found :", str)
	}

}

func TestBool00(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"false\", \"sub\": { \"key\":\"TRUE\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	val, serr := config.GetBool("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if !val {
		t.Error("Wrong value found :", val)
	}

	// getBool, stored a bool
	val, serr = config.GetBool("nope")
	if nil != serr {
		t.Error("Key 'nope' not found", serr)
	}
	if !val {
		t.Error("Wrong value found :", val)
	}

	// Search a sub key in a non existant sub item
	val, serr = config.GetBool("key.sub")
	if nil == serr {
		t.Error("Key 'key.sub' found")
	}

	// Missing value with a default value as string
	val, serr = config.GetBool("nope.sub", "true")
	if !val {
		t.Error("Wrong value found :", val)
	}
	if nil != serr {
		t.Error("Error found with default value", serr)
	}

	// Missing value with a default value as bool
	val, serr = config.GetBool("nope.sub", true)
	if !val {
		t.Error("Wrong value found :", val)
	}

	// Missing value with a default value
	val, serr = config.GetBool("sub.nope.key", "false")
	if val {
		t.Error("Wrong value found :", val)
	}

}

// Test GetValue from default
func TestDefault0(t *testing.T) {
	// Create configDefault with nil default
	def := &ConfigDefault{prefix: "Ctx_", values: nil, maxRecursion: 5}
	val, found := def.GetValue("nope")

	if found {
		t.Error("Missig key found in Default")
	}
	if nil != val {
		t.Error("Missing value found in default", val)
	}

	// Test from Env.
	os.Setenv("CTX_TEST0", "test")

	val, found = def.GetValue("test0")
	if !found {
		t.Error("key not found in Default")
	}
	if "test" != val {
		t.Error("Wrong value found in default", val)
	}

	// Defaults should be prioritary
	def.AddDefault("test0", "test1")
	val, found = def.GetValue("test0")
	if !found {
		t.Error("key not found in Default")
	}
	if "test1" != val {
		t.Error("Wrong value found in default", val)
	}

	found = def.AddDefault("sub.test", "somevalue")
	if !found {
		t.Error("Could not add default value")
	}
	val, found = def.GetValue("sub.test")
	if !found {
		t.Error("key not found in Default")
	}
	if "somevalue" != val {
		t.Error("Wrong value found in default", val)
	}

}

// Check GetString. for nested string
func TestFind0(t *testing.T) {
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"value\" }}"

	jsonBytes, err := ioutil.ReadAll(strings.NewReader(str))
	if err != nil {
		t.Error("ReadAll failed", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		t.Error("Unmarshall failed", err)
	}
	def := &ConfigDefault{prefix: "Ctx_", values: nil}

	config := ConfigImpl{values: obj, parent: nil, def: def}

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

	// Search a sub key in a non existant sub item
	str, serr = config.GetString("key.sub")
	if nil == serr {
		t.Error("Key 'key.sub' found")
	}

	// Missing value with a default value
	str, serr = config.GetString("nope.sub", "deflt")
	if "deflt" != str {
		t.Error("Wrong value found :", str)
	}

	// Existing value with a default value
	str, serr = config.GetString("sub.nope.key", "deflt")
	if "deflt" != str {
		t.Error("Wrong value found :", str)
	}

}

// vi:set fileencoding=utf-8 tabstop=4 ai
