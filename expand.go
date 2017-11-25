/*
Copyright (c) 2017 Jean-François PHILIPPE
Package goconfig read config files.
*/

package goconfig

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// matchEnd Find matching } of ${ in a string.
// val is the remaining of the string. i.e : after ${
// return pos of } in string or -1
func (c *ConfigImpl) matchEnd(val string) int {
	level := 1
	dfound := false
	for pos, rune := range val {
		if '{' == rune && dfound {
			// another ${ found
			level++
		} else if '}' == rune {
			level--
			if 0 == level {
				// match found !!
				return pos
			}
		}
		dfound = '$' == rune
	}
	return -1
}

// expandBuffer expand substitutions.
func (c *ConfigImpl) expandBuffer(buffer *bytes.Buffer, val string, deep uint) error {
	// Safe guard against infinite recursion
	if deep >= c.def.maxRecursion {
		return &ExpandRecursionError{step: deep}
	}
	remain := val // remaining of current string.

	var start, end int // Of ${xxx}

	// Search for ${
	start = strings.Index(remain, "${")
	for ; start >= 0; start = strings.Index(remain, "${") {
		buffer.WriteString(remain[:start])
		remain = remain[start+2:]
		end = c.matchEnd(remain)
		if end >= 0 {
			// extract key, and expand it if needed
			key, err := c.expand(strings.TrimSpace(remain[:end]), deep+1)
			if err != nil {
				return err
			}
			// Extra TrimSpace for keys.
			key = strings.TrimSpace(key)
			remain = remain[end+1:]
			subs, exists := c.find(key)
			if exists && nil != subs {
				// Convert found item into string
				substr := fmt.Sprint(subs)
				// enventually expand found value.
				err = c.expandBuffer(buffer, substr, deep+1)
				if err != nil {
					return err
				}
			} else {
				return &ExpandKeyError{key: key}
			}
		} else {
			buffer.WriteString("${")
		}
	}
	buffer.WriteString(remain)
	return nil
}

// expand expand a variable, replace ${var} within value.
func (c *ConfigImpl) expand(value string, deep uint) (string, error) {
	// if no recursion allowed return value.
	if 0 == c.def.maxRecursion {
		return value, nil
	}
	if deep >= c.def.maxRecursion {
		return value, &ExpandRecursionError{step: deep}
	}
	start := strings.Index(value, "${")
	if 0 <= start {
		// May need sustitutions ...
		// build a buffer with start of string
		buffer := bytes.NewBufferString(value[:start])
		buffer.Grow(len(value) * 2)

		err := c.expandBuffer(buffer, value[start:], deep)
		if err != nil {
			return value, err
		}
		return buffer.String(), nil
	}
	// Nothing to do
	return value, nil
}

// Expand expand a variable, replace ${var} within value.
func (c *ConfigImpl) Expand(value string) (string, error) {
	if 0 == c.def.maxRecursion {
		// No recursion allowed
		return value, nil
	} else {
		return c.expand(value, 0)
	}
}

// from https://gist.github.com/hvoecking/10772475  :
// The MIT License (MIT)
//
// Copyright (c) 2014 Heye Vöcking
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// Make a deep copy of an item, and expand any given string within.
func (c *ConfigImpl) Translate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	c.translateRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func (c *ConfigImpl) translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		c.translateRecursive(copy.Elem(), originalValue)

		// If it is an interface (which is very similar to a pointer), do basically the
		// same as for the pointer. Though a pointer is not the same as an interface so
		// note that we have to call Elem() after creating a new object because otherwise
		// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		c.translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

		// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			c.translateRecursive(copy.Field(i), original.Field(i))
		}

		// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			c.translateRecursive(copy.Index(i), original.Index(i))
		}

		// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			c.translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

		// Otherwise we cannot traverse anywhere so this finishes the the recursion

		// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		translatedString, _ := c.Expand(original.Interface().(string))
		copy.SetString(translatedString)

		// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}

// vi:set fileencoding=utf-8 tabstop=4 ai
