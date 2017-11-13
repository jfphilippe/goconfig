/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

 package goconfig

import (
	//"encoding/json"
	//"io"
	//"io/ioutil"
	//"path"
 	"errors"
 	"fmt"
 	"os"
	"strings"
)

type ConfigDefault struct {
	prefix string
	def map[string]interface{}
}

// GetPrefix read the prefix for env var.
func (c *ConfigDefault) GetPrefix() string {
	return c.prefix
}

// GetValue 
func (c *ConfigDefault) GetValue(key string) (interface{}, error) {
	// try first default value
	result, found := c.def[key]
	if ! found {
		// try Env vars
		name := c.prefix + strings.ToUpper(strings.Replace(key, ".", "_", -1))
		result, found = os.LookupEnv(name)
	}
	if found {
		return result, nil
	} else {
		return nil, nil
	}
}

/*
  define interface.   
 */
type ConfigImpl struct {
	values map[string]interface{}
	parent *ConfigImpl
	def *ConfigDefault
}

func (c *ConfigImpl) GetConfig(key string, defaultValue interface{}) (*GoConfig, error) {
	return nil, nil
}


// Get a String. the key mais be expressed with . to reach a nested item (aka key.sub.sub).
// If nothing is found and a default value is given, will return the default value. 
func (c *ConfigImpl) GetString(key string, deflt ...interface{}) (string, error) {
	// Get raw value
	raw, ok := c.get(key)
	strraw := ""
	// If not exists, 
	if !ok {
		if len(deflt) > 0 {
			// have a default value
			raw = deflt[0]
		} else {
			return "", errors.New("Key '" + key + "' does not exsists")
		}
	}
	switch v := raw.(type) {
	case string:
		strraw = v
	default:
		// Convert to string
		strraw = fmt.Sprint(v)
	}
	return strraw, nil
}

// sectionA extract a sub part of the map.
// if create is true an empty map will be created.
// may return nil if create is false and no map is found
func (c *ConfigImpl) sectionA(keys []string, create bool) *map[string]interface{} {
	vals := c.values
	for _, k := range keys {
		k := strings.TrimSpace(k)
		if k != "" {  // Ignore empty keys !!
			sub, ok := vals[k]
			if ok {
				// Check if can be casted
				if entry, ok := sub.(map[string]interface{}); ok {
					vals = entry
				} else {
					// Something eles, int, string , ...
					// return nil as we are trying to find a map
					return nil
				}
			} else {
				// Missing entry
				// create one if needed.
				if create {
					entry := make(map[string]interface{})
					vals[k] = &entry
					vals = entry
				} else {
					return nil
				}
			}
		}
	}

	return &vals
}

// get return the stored value as-is if exists
func (c *ConfigImpl) get(key string) (raw interface{}, exists bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	name := keys[len(keys)-1]
	entries := c.sectionA(section, false)
	if entries != nil {
		item, found := (*entries)[name]
		if found {
			return item, true
		}
	}
	return nil, false
}

// vi:set fileencoding=utf-8 tabstop=4 ai}
