/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	//"path"
	"strings"
)

type ConfigBuilder struct {
	conf *ConfigImpl
}

/*
  Instantiate a new builder
  prefix : prefix for env variable
  defaults : any defaults values may be nil.
*/
func NewBuilder(prefix string, defaults map[string]interface{}) *ConfigBuilder {
	prefix = strings.ToUpper(prefix)
	obj := make(map[string]interface{})
	def := &ConfigDefault{prefix: prefix, values: defaults, maxRecursion: 5}
	conf := &ConfigImpl{values: obj, parent: nil, def: def}
	result := &ConfigBuilder{conf: conf}

	return result
}

// GetConfig return current config
func (b *ConfigBuilder) GetConfig() GoConfig {
	return b.conf
}

// GetPrefix get current prefix
func (b *ConfigBuilder) GetPrefix() string {
	return b.conf.def.prefix
}

// AddDefault Add a default value
func (b *ConfigBuilder) AddDefault(key string, value interface{}) {
	b.conf.def.AddDefault(key, value)
}

func (b *ConfigBuilder) SetMaxRecursion(max uint) {
	b.conf.def.SetMaxRecursion(max)
}

// LoadJson Load a map from a Json Stream
// merge loaded value with previous one.
func (b *ConfigBuilder) LoadJson(r io.Reader) (GoConfig, error) {

	jsonBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		return nil, err
	}
	mergeMap(obj, b.conf.values)
	return b.conf, nil
}

// LoadJsonFile load from a file
func (b *ConfigBuilder) LoadJsonFile(filename string) (GoConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return b.LoadJson(r)
}

func (b *ConfigBuilder) LoadJsonFiles(ignoremissing bool, filenames ...string) (GoConfig, error) {
	for _, filename := range filenames {
		_, err := b.LoadJsonFile(filename)
		if nil != err && !(ignoremissing && os.IsNotExist(err)) {
			return nil, err
		}
	}
	return b.conf, nil
}

// vi:set fileencoding=utf-8 tabstop=4 ai
