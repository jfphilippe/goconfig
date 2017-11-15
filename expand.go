/*
 Copyright (c) 2017 Jean-François PHILIPPE
 Package goconfig read config files.
*/

package goconfig

import (
	"bytes"
	"errors"
	"fmt"
	"strconv" // Itoa
	"strings"
)

const (
	// Max recursive expand
	MaxRecursion = 5
)

// ErrMaxRecursion Error returned when max recursion is reached
var ErrMaxRecursion = errors.New("Max recursion reached :" + strconv.Itoa(MaxRecursion))

// Find matching } of ${ in a string.
// val is the remaining of the string. i.e : after ${
// return pos of } in string or -1
func (c *ConfigImpl) matchEnd(val string) int {
	// For now first } found
	// Later may handle ${xx${yy}zz} (nested items)
	return strings.Index(val, "}")
}

// expand expand substitutions.
func (c *ConfigImpl) expand(buffer *bytes.Buffer, val string, deep uint) error {
	// Safe guard against infinite recursion
	if deep > MaxRecursion {
		return ErrMaxRecursion
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
			key := strings.TrimSpace(remain[:end])
			remain = remain[end+1:]
			subs, exists := c.find(key)
			if exists && nil != subs {
				// Convert found item into string
				substr := fmt.Sprint(subs)
				// enventually expand found value.
				err := c.expand(buffer, substr, deep+1)
				if err != nil {
					return err
				}
			} else {
				return errors.New("Missing key '" + key + "'")
			}
		} else {
			buffer.WriteString("${")
		}
	}
	buffer.WriteString(remain)
	return nil
}

// Expand expand a variable, replace ${var} within value.
func (c *ConfigImpl) Expand(value string) (string, error) {
	start := strings.Index(value, "${")
	if 0 <= start {
		// May need sustitutions ...
		// build a buffer with start of string
		buffer := bytes.NewBufferString(value[:start])
		buffer.Grow(len(value) * 2)

		err := c.expand(buffer, value[start:], 0)
		if err != nil {
			return value, err
		}
		return buffer.String(), nil
	} else {
		// Nothing to do
		return value, nil
	}
}

// vi:set fileencoding=utf-8 tabstop=4 ai
