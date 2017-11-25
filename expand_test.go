/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

import (
	"strings"
	"testing"
)

// Check GetString.
func TestExpand0(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"subst\": \"${key}\"}"
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

	str, serr = config.GetString("subst")
	if "value" != str {
		t.Error("Key 'subst' not expanded, found", str)
	}

	// Existing value with a default value
	str, serr = config.GetString("subst", "deflt")
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check GetString. for nested string
func TestExpand1(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${nope}\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "true" != str {
		t.Error("Wrong value found :", str)
	}

	// Existing value with a default value
	str, serr = config.GetString("sub.nope.key", "deflt")
	if "deflt" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check expand, with space in key
func TestExpand2(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ nope}\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "true" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check expand, with space in key
func TestExpand3(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ nope }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "true" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check expand, with missing subst
func TestExpand4(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ none }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil == serr {
		t.Error("No error for missing subst")
	}
	if "${ none }" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check expand, with subst in defaults
func TestExpand5(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("none", "--")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ none }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found for value in defaults", serr)
	}
	if "--" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check expand, with  doted key in defaults
func TestExpand6(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("none", "--")
	builder.AddDefault("test.none", "-**-")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ test.none }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found for value in defaults", serr)
	}
	if "-**-" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check expand, with  doted key, full name is prioritary
// if full name is not found try "last" name (last part of dot)
func TestExpand7(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("test.none", "--")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ test.none }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found", err)
	}
	if "--" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check recursive expand, with  doted key
func TestExpand8(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("test.none", "${key}")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ test.none }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found", err)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check recursive expand, with  doted key
func TestExpand9(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("idx", "1")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"nope\": true, \"key1\":\"value\", \"sub\": { \"key\":\"${ key${idx} }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found", err)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}
}

// Check recursive expand, with  doted key
func TestExpand10(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("idx", "key")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ ${idx} }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found", err)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}
}

// Check recursive expand, with  doted key
func TestExpand11(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("env", "dev")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"dev\": {\"db\": {\"pwd\": \"azerty\"}}, \"int\":{\"db\":{\"pwd\":\"qwerty\"}}, \"database\": { \"pwd\":\"${ ${env}.db.pwd }\" }}"
	config, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("database.pwd")
	if nil != serr {
		t.Error("Error found", err)
	}
	if "azerty" != str {
		t.Error("Wrong value found :", str)
	}
}

// vi:set fileencoding=utf-8 tabstop=4 ai
