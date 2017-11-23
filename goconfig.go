/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

import (
	"fmt"
)

/*
  define interface.
*/
type GoConfig interface {
	GetConfig(key string) (GoConfig, error)
	GetString(key string, deflt ...interface{}) (string, error)
	//	GetInt64(key string, defaultValue interface{}) (int64, error)
	//	GetFloat(key string, defaultValue interface{}) (float64, error)
	GetBool(key string, deflt ...interface{}) (bool, error)
	//	GetDuration(key string, defaultValue interface{}) (Duration, error)
	//	GetAs(key string, target interface{}) error
	// Expand expand a value, replace
	// GetString(key, deflt string) string
	// GetString(key string) (string, error)
	// GetBool(key string, deflt bool) bool
	// GetBool(key string) (bool, error)
	Expand(value string) (string, error)
}

// Error for missing Key
type MissingKeyError struct {
	key string
}

// Error interface implementation
func (m MissingKeyError) Error() string {
	return fmt.Sprintf("Missing key : '%s'", m.key)
}

// Error while expanding ${xx} values (missing key)
type ExpandKeyError struct {
	key string
}

// Error interface implementation
func (m ExpandKeyError) Error() string {
	return fmt.Sprintf("Missing key : '%s'", m.key)
}

// Error max recursion reached while expanding
type ExpandRecursionError struct {
	step uint
}

// Error interface implementation
func (m ExpandRecursionError) Error() string {
	return fmt.Sprintf("Expand key, max recursion reached : %d", m.step)
}

// vi:set fileencoding=utf-8 tabstop=4 ai
