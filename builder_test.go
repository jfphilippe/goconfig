/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

import (
	"strings"
	"testing"
)

// test de valeurs globales avec deux sections
func TestBuilder0(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)

	if nil == builder {
		t.Error("NewBuilder failed")
	}
	// Check prefix is to Upper !
	if "CTX_" != builder.GetPrefix() {
		t.Error("Bad prefix '", builder.GetPrefix(), "' CTX_ expected")
	}
}

// Check Json parsing.
func TestBuilder1(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\"}"
	_, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}
}

// Check Json parsing.
func TestBuilder2(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\"}"
	str2 := "{ \"nope\": false, \"key2\":\"value2\"}"
	_, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}
	config, err2 := builder.LoadJson(strings.NewReader(str2))

	if nil != err2 {
		t.Error("LoadJson Failed", err)
	}

	// New key should not overide existing one
	val, serr := config.GetBool("nope")
	if nil != serr {
		t.Error("Key 'nope' not found", serr)
	}
	if !val {
		t.Error("Wrong value found :", val)
	}

	// previous key should still exists
	str, serr = config.GetString("key")
	if nil != serr {
		t.Error("Key 'key' not found", serr)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

	// previous key should still exists
	str, serr = config.GetString("key2")
	if nil != serr {
		t.Error("Key 'key2' not found", serr)
	}
	if "value2" != str {
		t.Error("Wrong value found :", str)
	}
}

// Check Json parsing.
// multiple parsing with sub-maps
func TestBuilder3(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\",  \"sub\": { \"bool\": false }}"
	str2 := "{ \"nope\": false, \"key2\":\"value2\",  \"sub\": { \"bool\": true, \"string\": \"test\" }}"
	_, err := builder.LoadJson(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJson Failed", err)
	}
	config, err2 := builder.LoadJson(strings.NewReader(str2))

	if nil != err2 {
		t.Error("LoadJson Failed", err)
	}

	// New key should not overide existing one
	val, serr := config.GetBool("sub.bool")
	if nil != serr {
		t.Error("Key 'sub.bool' not found", serr)
	}
	if val {
		t.Error("Wrong value found :", val)
	}

	// New key in sub map
	str, serr = config.GetString("sub.string")
	if nil != serr {
		t.Error("Key 'sub.string' not found", serr)
	}
	if "test" != str {
		t.Error("Wrong value found :", str)
	}

}

// vi:set fileencoding=utf-8 tabstop=4 ai
