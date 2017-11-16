/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

/*
  define interface.
*/
type GoConfig interface {
	GetConfig(key string) (GoConfig, error)
	GetString(key string, deflt ...interface{}) (string, error)
	//	GetInt64(key string, defaultValue interface{}) (int64, error)
	//	GetFloat(key string, defaultValue interface{}) (float64, error)
	//	GetBool(key string, defaultValue interface{}) (bool, error)
	//	GetDuration(key string, defaultValue interface{}) (Duration, error)
	//	GetAs(key string, target interface{}) error
	// Expand expand a value, replace
	Expand(value string) (string, error)
}

// vi:set fileencoding=utf-8 tabstop=4 ai
