/*
Copyright (c) 2017 Jean-François PHILIPPE
Package goconfig read config files.
*/

package goconfig

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

// Check GetString.
func TestExpand0(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\", \"subst\": \"${key}\"}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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

// Check recursive expand, with  doted key
func TestExpand12(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.AddDefault("test.none", "${key}")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ test.none \" }}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found", err)
	}
	if "${ test.none " != str {
		t.Error("Wrong value found :", str)
	}

}

// Check recursive expand, With Max recursion to 0
func TestExpand13(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.SetMaxRecursion(0)
	builder.AddDefault("test.none", "${key}")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"${ test.none }\" }}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Error found", err)
	}
	// With max recursion to 0 , should return the value
	if "${ test.none }" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check recursive expand, with  dmex recursion to 1
func TestExpand14(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.SetMaxRecursion(1)
	builder.AddDefault("env", "dev")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"dev\": {\"db\": {\"pwd\": \"azerty\"}}, \"int\":{\"db\":{\"pwd\":\"qwerty\"}}, \"database\": { \"pwd\":\"${ ${env}.db.pwd }\" }}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("database.pwd")
	if nil == serr {
		t.Error("Error should be found")
	}
	// Max recursion reached , should return the value !
	if "${ ${env}.db.pwd }" != str {
		t.Error("Wrong value found :", str)
	}
}

// Check recursive expand, with  dmex recursion to 1
func TestExpand15(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	builder.SetMaxRecursion(5)
	builder.AddDefault("env", "dev")
	builder.AddDefault("nope.none", "-**-")
	str := "{ \"dev\": {\"db\": {\"pwd\": \"a ${int.db.pwd}\"}}, \"int\":{\"db\":{\"pwd\":\"b ${ ${env}.db.pwd}\"}}, \"database\": { \"pwd\":\"${ ${env}.db.pwd }\" }}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	// Search a key as string
	str, serr := config.GetString("database.pwd")
	if nil == serr {
		t.Error("Error should be found")
	}
	// Max recursion reached , should return the value !
	if "${ ${env}.db.pwd }" != str {
		t.Error("Wrong value found :", str)
	}
}

// Check GetFloat
func TestTranslate16(t *testing.T) {
	str := "{ \"string2\": \"${key}\", \"key\":\"value\", \"sub\": { \"key\":\"value2\" }}"
	builder := NewBuilder("Ctx_", nil)
	builder.SetMaxRecursion(5)
	builder.AddDefault("env", "dev")
	builder.AddDefault("nope.none", "-**-")
	_, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	jsonBytes, err := ioutil.ReadAll(strings.NewReader(str))
	if err != nil {
		t.Error("ReadAll failed", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		t.Error("Unmarshall failed", err)
	}
	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	config := builder.conf

	str, err = config.GetString("string2")
	if nil != err {
		t.Error("GetString Failed", err)
	}
	if "value" != str {
		t.Error("Wrong Value returned, expecting value :", str)
	}

	m := make(map[string]interface{})
	m["key0"] = "test"
	m["key1"] = "${nope}"

	m2 := make(map[string]interface{})
	m2["sub2"] = m
	m2["root"] = 12
	m2["string"] = "${key}"
	m2["array"] = []string{"${sub.key}", "array"}

	m3 := config.Translate(m2)
	if nil == m3 {
		t.Error("Translate should not have returned nil")
	}

	switch tsrc := m3.(type) {
	case map[string]interface{}:
		val := tsrc["string"]
		if nil == val {
			t.Error("m3[string] should not be nil")
		}
		switch tsrc2 := val.(type) {
		case string:
			if "value" != tsrc2 {
				t.Error("Wrong value found :", val)
			}
		default:
			t.Error("tsrc[string] should be a string")
		}
		// Test for array
		val = tsrc["array"]
		if nil == val {
			t.Error("m3[array] should not be nil")
		}
		switch tsrc2 := val.(type) {
		case []string:
			if "value2" != tsrc2[0] {
				t.Error("Wrong value found :", val)
			}
		default:
			t.Error("tsrc[string] should be a string")
		}
	default:
		t.Error("m3 should be a map")
	}
}

// vi:set fileencoding=utf-8 tabstop=4 ai
