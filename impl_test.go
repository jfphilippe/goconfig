/*
Copyright (c) Jean-François PHILIPPE 2017-2018
Package goconfig read config files.
*/

package goconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

// Check GetString.
func TestGetString0(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"nope\": true, \"key\":\"value\"}"
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
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
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
	str := "{ \"nope\": true, \"key\":\"false\", \"sub\": { \"key\":\"TRUE\", \"int\":0 }}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	// Search a key as string
	val, serr := config.GetBool("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if !val {
		t.Error("Wrong value found :", val)
	}
	// Search a key as int
	val, serr = config.GetBool("sub.int")
	if nil != serr {
		t.Error("Key 'sub.int' not found", serr)
	}
	if val {
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

// Test Duration
func TestDuration0(t *testing.T) {
	builder := NewBuilder("Ctx_", nil)
	str := "{ \"t0\":\"2h\" , \"key\":\"false\", \"sub\": { \"key\":\"TRUE\", \"int\":260 }}"
	config, err := builder.LoadJSON(strings.NewReader(str))

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	// Search a key as string
	val, serr := config.GetDuration("t0")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if val != (time.Hour * 2) {
		t.Error("Wrong value found :", val)
	}

	// Value as int
	val, serr = config.GetDuration("sub.int")
	if nil == serr {
		t.Error("GetDuration sub.int should fail")
	}

	// Missing value with a default value
	val, serr = config.GetDuration("sub.nope.key", time.Minute*3)
	if (time.Minute * 3) != val {
		t.Error("Wrong value found :", val)
	}

	// Missing value without a default value
	val, serr = config.GetDuration("sub.nope.key")
	if nil == serr {
		t.Error("Should be error")
	}
	if (time.Second * 0) != val {
		t.Error("Wrong value found :", val)
	}

}

// Test GetValue from default
func TestDefault0(t *testing.T) {
	// Create configDefault with nil default
	def := &ConfigDefault{prefix: "Ctx_", values: nil, maxRecursion: 5}

	if 5 != def.GetMaxRecursion() {
		t.Error("Wrong MaxRecursion value found", def.GetMaxRecursion())
	}

	if "Ctx_" != def.GetPrefix() {
		t.Error("Wrong Prefix Found '", def.GetPrefix(), "' Ctx_ expected")
	}
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

	found = def.AddDefault("nope", nil)
	if found {
		t.Error("Should not be able to add a nil value")
	}

}

// Test GetValue from default
func TestDefault1(t *testing.T) {
	// Create configDefault with nil default
	def := &ConfigDefault{prefix: "Ctx_", values: nil, maxRecursion: 5}
	// Defaults should be prioritary
	def.AddDefault("test0", "test1")
	val, found := def.GetValue("test0")
	if !found {
		t.Error("key not found in Default")
	}

	config := ConfigImpl{values: nil, parent: nil, def: def}
	val, err := config.GetString("test0")
	if nil != err {
		t.Error("key not found in Default", err)
	}
	if "test1" != val {
		t.Error("Wrong value found :", val)
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
		t.Error("LoadJSON Failed", err)
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

// Check GetString. for nested string
func TestSetValue0(t *testing.T) {
	str := "{ \"nope\": true, \"key\":\"value\", \"sub\": { \"key\":\"value\" }}"

	jsonBytes, err := ioutil.ReadAll(strings.NewReader(str))
	if err != nil {
		t.Error("ReadAll failed", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		t.Error("Unmarshall failed", err)
	}
	var defaults = make(map[string]interface{})
	defaults["def"] = "default"
	defaults["some"] = 0
	def := &ConfigDefault{prefix: "Ctx_", values: defaults}

	config := ConfigImpl{values: obj, parent: nil, def: def}

	if nil != err {
		t.Error("LoadJSON Failed", err)
	}

	var value = make(map[string]interface{})
	value["key2"] = false
	value["key3"] = "nope"
	config.SetValue("sub", value)

	// Search a key as string
	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "value" != str {
		t.Error("Wrong value found :", str)
	}

	// Search a sub key in a non existant sub item
	str, serr = config.GetString("sub.key3")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}
	if "nope" != str {
		t.Error("Wrong value found :", str)
	}

}

// Check GetString. for nested string
func TestGetConfig0(t *testing.T) {
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
		t.Error("LoadJSON Failed", err)
	}

	str, serr := config.GetString("sub.key")
	if nil != serr {
		t.Error("Key 'sub.key' not found", serr)
	}

	conf, serr := config.GetConfig("missing")
	if nil != conf {
		t.Error("Should return nil")
	}
	if nil == serr {
		t.Error("GetConfig from missing key should return error")
	}

	// Get Existing SubConfig
	conf, serr = config.GetConfig("sub")
	if nil == conf {
		t.Error("GetConfig(sub) Should not return nil")
	}
	if nil != serr {
		t.Error("GetConfig from existing key should not return error", serr)
	}
	str, serr = conf.GetString("key")
	if nil != serr {
		t.Error("Sub Key 'key' not found", serr)
	}
	if "value" != str {
		t.Error("Wrong value found for sub config", str, " expecting value")
	}

}

// Check GetString. for nested string
func TestGetUInt00(t *testing.T) {
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
		t.Error("LoadJSON Failed", err)
	}

	// GetUint from a string not an int
	_, serr := config.GetUint("nope", 5)
	if nil == serr {
		t.Error("Should return error")
	}

	val, serr := config.GetUint("missing", int(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetUint("missing", uint(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetUint("missing", int8(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetUint("missing", uint8(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetUint("missing", int16(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetUint("missing", uint16(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetUint("missing", int32(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetUint("missing", uint32(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetUint("missing", int64(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetUint("missing", uint64(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetUint("missing", float32(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetUint("missing", float64(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetUint("missing", "126")
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 126 != val {
		t.Error("Wrong value returned ", val, " expecting 126")
	}

}

// Check GetInt
func TestGetInt00(t *testing.T) {
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
		t.Error("LoadJSON Failed", err)
	}

	// GetInt from a string not an int
	_, serr := config.GetInt("nope", 5)
	if nil == serr {
		t.Error("Should return error")
	}

	val, serr := config.GetInt("missing", int(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetInt("missing", uint(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetInt("missing", int8(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetInt("missing", uint8(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetInt("missing", int16(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetInt("missing", uint16(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetInt("missing", int32(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetInt("missing", uint32(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetInt("missing", int64(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetInt("missing", uint64(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetInt("missing", float32(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetInt("missing", float64(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetInt("missing", "126")
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 126 != val {
		t.Error("Wrong value returned ", val, " expecting 126")
	}

}

// Check GetFloat
func TestGetFloat00(t *testing.T) {
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
		t.Error("LoadJSON Failed", err)
	}

	// GetFloat from a string not an int
	_, serr := config.GetFloat("nope", 5)
	if nil == serr {
		t.Error("Should return error")
	}

	val, serr := config.GetFloat("missing", int(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetFloat("missing", uint(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetFloat("missing", int8(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetFloat("missing", uint8(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetFloat("missing", int16(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetFloat("missing", uint16(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetFloat("missing", int32(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetFloat("missing", uint32(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetFloat("missing", int64(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetFloat("missing", uint64(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetFloat("missing", float32(5))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 5 != val {
		t.Error("Wrong value returned ", val, " expecting 5")
	}

	val, serr = config.GetFloat("missing", float64(6))
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 6 != val {
		t.Error("Wrong value returned ", val, " expecting 6")
	}

	val, serr = config.GetFloat("missing", "126")
	if nil != serr {
		t.Error("Should not raise error", serr)
	}
	if 126 != val {
		t.Error("Wrong value returned ", val, " expecting 126")
	}

}

// vi:set fileencoding=utf-8 tabstop=4 ai
