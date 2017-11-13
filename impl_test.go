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


// vi:set fileencoding=utf-8 tabstop=4 ai