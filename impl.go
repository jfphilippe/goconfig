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
	"strconv"
	"strings"
	"time"
)

// ConfigDefault Store commons values
type ConfigDefault struct {
	prefix       string
	values       map[string]interface{}
	maxRecursion uint
}

// GetMaxRecursion return current max recursion.
// Return 0 if disabled.
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
	if nil == c.values {
		return nil, false
	}
	vals := c.values
	section := keys[:len(keys)-1]
	// name is last part
	name := strings.TrimSpace(keys[len(keys)-1])
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
		keys := strings.Split(key, ".")
		section := keys[:len(keys)-1]
		// name is last part
		name := keys[len(keys)-1]

		m := subMap(&c.values, section, true)
		if nil != m {
			smap := *m
			result, found = smap[name]
		}
	}
	if !found {
		// try Env vars
		name := c.prefix + key
		// Convert to UpperCase but first replace '.' with '_'
		name = strings.ToUpper(strings.Replace(name, ".", "_", -1))
		result, found = os.LookupEnv(name)
		if found {
			return result, true
		}
		return nil, false
	}

	return result, found
}

// AddDefault Add a default value
func (c *ConfigDefault) AddDefault(key string, value interface{}) bool {
	if nil != value {
		if nil == c.values {
			c.values = make(map[string]interface{})
		}
		keys := strings.Split(key, ".")
		section := keys[:len(keys)-1]
		// name is last part
		name := keys[len(keys)-1]

		m := subMap(&c.values, section, true)
		if nil != m {
			smap := *m
			smap[name] = value
			return true
		}
	}
	return false
}

// ConfigImpl implements GoConfig interface
type ConfigImpl struct {
	values map[string]interface{}
	parent *ConfigImpl
	def    *ConfigDefault
}

// GetConfig Create a config using a subtree of the currents values
func (c *ConfigImpl) GetConfig(key string) (GoConfig, error) {
	keys := strings.Split(key, ".")
	values := subMap(&c.values, keys, false)
	if nil == values {
		return nil, errors.New("Key '" + key + "' does not exsists")
	}
	return &ConfigImpl{values: *values, parent: c, def: c.def}, nil
}

// SetValue store a value (value may be a map[string]interface{})
func (c *ConfigImpl) SetValue(key string, value interface{}) bool {
	if nil != value {
		keys := strings.Split(key, ".")
		section := keys[:len(keys)-1]
		// name is last part
		name := strings.TrimSpace(keys[len(keys)-1])
		entries := subMap(&c.values, section, true)
		if nil != entries {
			// Do Not override an existing entry
			item, found := (*entries)[name]
			if found {
				// if value AND item are map[string]interface merge recursively
				// otherwise do nothing
				switch tsrc := value.(type) {
				case map[string]interface{}:
					switch tdest := item.(type) {
					case map[string]interface{}:
						mergeMap(tsrc, tdest)
						return true
					}
				}
			} else {
				(*entries)[name] = value
				return true
			}
		} // entries should always be not nil
	}
	return false
}

// GetString  get a String. the key may be expressed with . to reach a nested item (aka key.sub.sub).
// If nothing is found and a default value is given, will return the default value.
func (c *ConfigImpl) GetString(key string, deflt ...interface{}) (string, error) {
	// Get raw value
	raw, err := c.getExpand(key, deflt...)
	// If not exists,
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

// GetBool return a value as a boolean
func (c *ConfigImpl) GetBool(key string, defaultValue ...interface{}) (bool, error) {
	// Get raw value
	raw, err := c.getExpand(key, defaultValue...)
	// If not exists,
	if nil != raw {
		// Convert to string....
		switch v := raw.(type) {
		case bool:
			return v, err
		case string:
			return strconv.ParseBool(v)
		default:
			// Convert to string
			strval := fmt.Sprint(v)
			return strconv.ParseBool(strval)
		}
	}
	return false, err
}

// GetDuration read a Duration from configuration.
func (c *ConfigImpl) GetDuration(key string, defaultValue ...interface{}) (time.Duration, error) {
	// Get raw value
	raw, err := c.getExpand(key, defaultValue...)
	// If not exists,
	if nil != raw {
		// Convert to string....
		switch v := raw.(type) {
		case time.Duration:
			return v, err
		case string:
			return time.ParseDuration(v)
		default:
			// Convert to string
			strval := fmt.Sprint(v)
			return time.ParseDuration(strval)
		}
	}
	return 0 * time.Second, err
}

// GetInt read an Int from configuration.
func (c *ConfigImpl) GetInt(key string, defaultValue ...interface{}) (int64, error) {
	// Get raw value
	raw, err := c.getExpand(key, defaultValue...)
	// If not exists,
	if nil != raw {
		switch val := raw.(type) {
		case int:
			return int64(val), nil
		case uint:
			return int64(val), nil
		case int8:
			return int64(val), nil
		case uint8:
			return int64(val), nil
		case int16:
			return int64(val), nil
		case uint16:
			return int64(val), nil
		case int32:
			return int64(val), nil
		case uint32:
			return int64(val), nil
		case int64:
			return int64(val), nil
		case uint64:
			return int64(val), nil
		case float32:
			return int64(val), nil
		case float64:
			return int64(val), nil
		case string:
			return strconv.ParseInt(val, 0, 64)
		default:
			// Convert to string
			strval := fmt.Sprint(val)
			return strconv.ParseInt(strval, 0, 64)
		}
	}
	return 0, err
}

// GetUint read an uint from configuration.
func (c *ConfigImpl) GetUint(key string, defaultValue ...interface{}) (uint64, error) {
	// Get raw value
	raw, err := c.getExpand(key, defaultValue...)
	// If not exists,
	if nil != raw {
		switch val := raw.(type) {
		case int:
			return uint64(val), nil
		case uint:
			return uint64(val), nil
		case int8:
			return uint64(val), nil
		case uint8:
			return uint64(val), nil
		case int16:
			return uint64(val), nil
		case uint16:
			return uint64(val), nil
		case int32:
			return uint64(val), nil
		case uint32:
			return uint64(val), nil
		case int64:
			return uint64(val), nil
		case uint64:
			return uint64(val), nil
		case float32:
			return uint64(val), nil
		case float64:
			return uint64(val), nil
		case string:
			return strconv.ParseUint(val, 0, 64)
		default:
			// Convert to string
			strval := fmt.Sprint(val)
			return strconv.ParseUint(strval, 0, 64)
		}
	}
	return 0, err
}

// GetFloat read a float from configuration.
func (c *ConfigImpl) GetFloat(key string, defaultValue ...interface{}) (float64, error) {
	// Get raw value
	raw, err := c.getExpand(key, defaultValue...)
	// If not exists,
	if nil != raw {
		switch val := raw.(type) {
		case int:
			return float64(val), nil
		case uint:
			return float64(val), nil
		case int8:
			return float64(val), nil
		case uint8:
			return float64(val), nil
		case int16:
			return float64(val), nil
		case uint16:
			return float64(val), nil
		case int32:
			return float64(val), nil
		case uint32:
			return float64(val), nil
		case int64:
			return float64(val), nil
		case float32:
			return float64(val), nil
		case float64:
			return float64(val), nil
		case string:
			return strconv.ParseFloat(val, 64)
		default:
			// Convert to string
			strval := fmt.Sprint(val)
			return strconv.ParseFloat(strval, 64)
		}
	}
	return 0.0, err
}

// subMap extract a sub part of the map.
// if create is true an empty map will be created.
// may return nil if create is false and no map is found or if the item found is not a map
func subMap(values *map[string]interface{}, keys []string, create bool) *map[string]interface{} {
	vals := *values
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
					vals[k] = entry
					vals = entry
				} else {
					return nil
				}
			}
		}
	}

	return &vals
}

// mergeMap merge two maps. Copy all entries from the first to the second if
func mergeMap(src, dest map[string]interface{}) {
	// iterate over all key, values from src
	for k, v := range src {
		// search in dest value for same key
		v2, found := dest[k]
		if !found {
			// not found : set it
			dest[k] = v
		} else {
			// if v AND v2 are map[string]interface merge recursively
			switch tsrc := v.(type) {
			case map[string]interface{}:
				switch tdest := v2.(type) {
				case map[string]interface{}:
					mergeMap(tsrc, tdest)
				}
			}
		}
	}
}

// getExpand return the stored value, or default, and expand if value is a string
func (c *ConfigImpl) getExpand(key string, deflt ...interface{}) (raw interface{}, err error) {
	result, found := c.get(key, deflt...)
	if !found {
		return nil, &MissingKeyError{key: key}
	}

	switch v := result.(type) {
	case string:
		return c.Expand(v)
	default:
		return c.Translate(result), nil
	}

}

// get return the stored value as-is if exists
func (c *ConfigImpl) get(key string, deflt ...interface{}) (raw interface{}, exists bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	// name is last part
	name := keys[len(keys)-1]
	entries := subMap(&c.values, section, false)
	if nil != entries {
		item, found := (*entries)[name]
		if found {
			return item, true
		}
	}
	// if nothing found try defaults
	item, found := c.def.GetValue(key)
	if found {
		return item, true
	}
	// fallback try default param
	if len(deflt) > 0 {
		return deflt[0], true
	}
	// Nothing found
	return nil, false
}

// find return the stored value, search eventualy in parents Config and Default.
func (c *ConfigImpl) find(key string) (raw interface{}, exists bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	name := keys[len(keys)-1]
	conf := c
	var entries *map[string]interface{}
	for conf != nil {
		entries = subMap(&conf.values, section, false)
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
