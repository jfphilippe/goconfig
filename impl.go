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

// Store commons values
type ConfigDefault struct {
	prefix       string
	values       map[string]interface{}
	maxRecursion uint
}

// GetMaxRecursion
func (c *ConfigDefault) GetMaxRecursion() uint {
	return c.maxRecursion
}

// SetMaxRecursion set maxRecursion Value
func (c *ConfigDefault) SetMaxRecursion(max uint) {
	c.maxRecursion = max
}

// GetPrefix read the prefix for env var.
func (c *ConfigDefault) GetPrefix() string {
	return c.prefix
}

func (c *ConfigDefault) getValue(keys []string) (interface{}, bool) {
	vals := c.values
	section := keys[:len(keys)-1]
	// name is last part
	name := keys[len(keys)-1]
	// Search for the required section
	for _, k := range section {
		k := strings.TrimSpace(k)
		if k != "" { // Ignore empty keys !!
			sub, ok := vals[k]
			if ok {
				// Check if can be casted
				if entry, ok := sub.(map[string]interface{}); ok {
					vals = entry
				} else {
					// Something eles, int, string , ...
					// return nil as we are trying to find a map
					return nil, false
				}
			}
		}
	}

	val, ok := vals[name]
	return val, ok
}

// GetValue try to get a value from defaults.
// search first a value in default map, the into Env vars.
func (c *ConfigDefault) GetValue(key string) (interface{}, bool) {
	found := false
	var result interface{}

	// try first default value
	if nil != c.values {
		result, found = c.values[key]
	}
	if !found {
		// try Env vars
		name := c.prefix + key
		// Convert to UpperCase but first replace '.' with '_'
		name = strings.ToUpper(strings.Replace(name, ".", "_", -1))
		result, found = os.LookupEnv(name)
		if found {
			return result, true
		} else {
			return nil, false
		}
	}

	return result, found
}

// Add a default value
func (b *ConfigDefault) AddDefault(key string, value interface{}) {
	if nil == b.values {
		b.values = make(map[string]interface{})
	}
	b.values[key] = value
}

/*
  implements interface.
*/
type ConfigImpl struct {
	values map[string]interface{}
	parent *ConfigImpl
	def    *ConfigDefault
}

// Create a config using a subtree of the currents values
func (c *ConfigImpl) GetConfig(key string) (GoConfig, error) {
	keys := strings.Split(key, ".")
	values := c.sectionA(keys, false)
	if nil == values {
		return nil, errors.New("Key '" + key + "' does not exsists")
	}
	return &ConfigImpl{values: *values, parent: c, def: c.def}, nil
}

// Get a String. the key mais be expressed with . to reach a nested item (aka key.sub.sub).
// If nothing is found and a default value is given, will return the default value.
func (c *ConfigImpl) GetString(key string, deflt ...interface{}) (string, error) {
	// Get raw value
	raw, err := c.getExpand(key, deflt...)
	if nil != raw {
		// Convert to string....
		switch v := raw.(type) {
		case string:
			return v, err
		default:
			// Convert to string
			return fmt.Sprint(v), err
		}
	}
	return "", err
}

// sectionA extract a sub part of the map.
// if create is true an empty map will be created.
// may return nil if create is false and no map is found or if the item found is not a map
func (c *ConfigImpl) sectionA(keys []string, create bool) *map[string]interface{} {
	vals := c.values
	for _, k := range keys {
		k := strings.TrimSpace(k)
		if k != "" { // Ignore empty keys !!
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

// getExpand return the stored value, or default, and expand if value is a string
func (c *ConfigImpl) getExpand(key string, deflt ...interface{}) (raw interface{}, err error) {
	result, found := c.get(key)
	if !found {
		if len(deflt) > 0 {
			result = deflt[0]
			found = nil != result
		}
		if !found {
			return nil, &MissingKeyError{key: key}
		}
	}

	switch v := result.(type) {
	case string:
		return c.Expand(v)
	default:
		return result, nil
	}

}

// get return the stored value as-is if exists
// return the value and 'true' if key was found.
func (c *ConfigImpl) get(key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	// name is last part
	name := keys[len(keys)-1]
	entries := c.sectionA(section, false)
	if entries != nil {
		item, found := (*entries)[name]
		if found {
			return item, true
		}
	}
	// Otherwise try defaults
	return c.def.getValue(keys)
}

// find return the stored value, search eventualy in parents Config and Default.
func (c *ConfigImpl) find(key string) (raw interface{}, exists bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	name := keys[len(keys)-1]
	conf := c
	var entries *map[string]interface{}
	for conf != nil {
		entries = conf.sectionA(section, false)
		if entries != nil {
			item, found := (*entries)[name]
			if found {
				return item, true
			}
		}
		conf = conf.parent
	}
	// fail over, search in defaults
	// first full name
	return c.def.GetValue(key)
}

// vi:set fileencoding=utf-8 tabstop=4 ai
