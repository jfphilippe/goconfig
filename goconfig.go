/*
Copyright (c) Jean-François PHILIPPE 2017-2018
Package goconfig read config files.
*/

package goconfig

import (
	"fmt"
	"time"
)

// GoConfig define interface.
type GoConfig interface {
	// Extract a sub part of config.
	GetConfig(key string) (GoConfig, error)
	GetString(key string, deflt ...interface{}) (string, error)
	GetInt(key string, defaultValue ...interface{}) (int64, error)
	GetUint(key string, defaultValue ...interface{}) (uint64, error)
	GetFloat(key string, defaultValue ...interface{}) (float64, error)
	GetBool(key string, deflt ...interface{}) (bool, error)
	GetDuration(key string, deflt ...interface{}) (time.Duration, error)
	// GetString(key, deflt string) string
	// GetBool(key string, deflt bool) bool
	Expand(value string) (string, error)
}

// ParseError Error for missing Key
type ParseError struct {
	line int
	msg  string
}

// Error interface implementation
func (m ParseError) Error() string {
	return fmt.Sprintf("Parse Error line '%d' : '%s'", m.line, m.msg)
}

// MissingKeyError Error for missing Key
type MissingKeyError struct {
	key string
}

// Error interface implementation
func (m MissingKeyError) Error() string {
	return fmt.Sprintf("Missing key : '%s'", m.key)
}

// ExpandKeyError Error while expanding ${xx} values (missing key)
type ExpandKeyError struct {
	key string
}

// Error interface implementation
func (m ExpandKeyError) Error() string {
	return fmt.Sprintf("Missing key : '%s'", m.key)
}

// ExpandRecursionError Error max recursion reached while expanding
type ExpandRecursionError struct {
	step uint
}

// Error interface implementation
func (m ExpandRecursionError) Error() string {
	return fmt.Sprintf("Expand key, max recursion reached : %d", m.step)
}

// vi:set fileencoding=utf-8 tabstop=4 ai
