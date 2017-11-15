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

// vi:set fileencoding=utf-8 tabstop=4 ai
