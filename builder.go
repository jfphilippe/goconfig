/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

import (
	"encoding/json"
	"io"
	"io/ioutil"
	//"path"
	"strings"
)

type ConfigBuilder struct {
	def *ConfigDefault
}

/*
  Instantiate a new builder
  prefix : prefix for env variable, nil if must be ignored
  defaults : any defaults values.
*/
func NewBuilder(prefix string, defaults map[string]interface{}) *ConfigBuilder {
	prefix = strings.ToUpper(prefix)
	def := &ConfigDefault{prefix: prefix, def: defaults}
	result := &ConfigBuilder{def: def}

	return result
}

func (b *ConfigBuilder) GetPrefix() string {
	return b.def.prefix
}

func (b *ConfigBuilder) AddDefault(key string, value interface{}) {
	b.def.AddDefault(key, value)
}

func (b *ConfigBuilder) LoadJson(r io.Reader) (GoConfig, error) {

	jsonBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		return nil, err
	}
	return &ConfigImpl{values: obj, parent: nil, def: b.def}, nil
}

// vi:set fileencoding=utf-8 tabstop=4 ai
