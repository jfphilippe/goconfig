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
	"path"
	"strings"
)

// ConfigBuilder Used to create Config Objects and parse config files.
type ConfigBuilder struct {
	conf               *ConfigImpl
	ignoreMissingFiles bool
}

// NewBuilder Instantiate a new builder
//  prefix : prefix for env variable
//  defaults : any defaults values may be nil.
func NewBuilder(prefix string, defaults map[string]interface{}) *ConfigBuilder {
	prefix = strings.ToUpper(prefix)
	obj := make(map[string]interface{})
	def := &ConfigDefault{prefix: prefix, values: defaults, maxRecursion: 5}
	conf := &ConfigImpl{values: obj, parent: nil, def: def}
	result := &ConfigBuilder{conf: conf, ignoreMissingFiles: false}

	return result
}

// SetIgnoreMissingFiles should builder ignore missing files or not
func (b *ConfigBuilder) SetIgnoreMissingFiles(value bool) {
	b.ignoreMissingFiles = value
}

// IgnoreMissingFiles check if builder ignore missing files or not
func (b *ConfigBuilder) IgnoreMissingFiles() bool {
	return b.ignoreMissingFiles
}

// Config return current config
func (b *ConfigBuilder) Config() GoConfig {
	return b.conf
}

// Prefix get current prefix
func (b *ConfigBuilder) Prefix() string {
	return b.conf.def.prefix
}

// AddDefault Add a default value
func (b *ConfigBuilder) AddDefault(key string, value interface{}) {
	b.conf.def.AddDefault(key, value)
}

// SetMaxRecursion configure max expand recursion.
// once the limit reached an error will be returned
// set to 0 to disable expansion.
func (b *ConfigBuilder) SetMaxRecursion(max uint) {
	b.conf.def.SetMaxRecursion(max)
}

// MaxRecursion return current value
func (b *ConfigBuilder) MaxRecursion() uint {
	return b.conf.def.maxRecursion
}

// LoadJSON Load a map from a Json Stream
// merge loaded value with previous one.
func (b *ConfigBuilder) LoadJSON(r io.Reader) (GoConfig, error) {

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

// LoadTxt Load a map from a text Stream
// merge loaded value with previous one.
func (b *ConfigBuilder) LoadTxt(r io.Reader) (GoConfig, error) {

	scanner := bufio.NewScanner(r)
	lineNb := 0
	obj := make(map[string]interface{})
	conf := &ConfigImpl{values: obj, parent: nil, def: b.conf.def}
	for scanner.Scan() {
		lineNb++
		line := strings.TrimSpace(scanner.Text())
		// ignore empty lines and comments
		if "" != line && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			// Parse key=value
			words := strings.SplitN(line, "=", 2)
			if len(words) != 2 {
				return b.conf, &ParseError{line: lineNb, msg: "missing '=' : '" + line + "'"}
			}
			key := strings.TrimSpace(words[0])
			value := strings.TrimSpace(words[1])
			// Set Value in a spare config
			conf.SetValue(key, value)
		}

	}
	// Merge new config and current one.
	mergeMap(conf.values, b.conf.values)

	return b.conf, nil
}

// LoadJSONFile load from a file
func (b *ConfigBuilder) LoadJSONFile(filename string) (GoConfig, error) {
	f, err := os.Open(filename)
	if nil != err {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return b.LoadJSON(r)
}

// LoadTxtFile load from a file
func (b *ConfigBuilder) LoadTxtFile(filename string) (GoConfig, error) {
	f, err := os.Open(filename)
	if nil != err {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return b.LoadTxt(r)
}

// LoadFiles load from files. Guess file type by reading extension.
// When extension is .json parse it as a json file,
// otherwise as a txt file
func (b *ConfigBuilder) LoadFiles(filenames ...string) (GoConfig, error) {
	// wich method to use to parse the file, default to Txt
	parser := b.LoadTxt
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if nil != err {
			if !(b.ignoreMissingFiles && os.IsNotExist(err)) {
				return nil, err
			}
		} else {
			// Choose a parser
			if ".json" == path.Ext(filename) {
				parser = b.LoadJSON
			} else {
				parser = b.LoadTxt
			}

			// parse file
			r := bufio.NewReader(f)
			_, err = parser(r)
			f.Close()

			// if any error stop
			if nil != err {
				return nil, err
			}
		}
	}
	return b.conf, nil
}

// vi:set fileencoding=utf-8 tabstop=4 ai
