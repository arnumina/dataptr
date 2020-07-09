/*
#######
##           __     __            __
##       ___/ /__ _/ /____ ____  / /_____
##      / _  / _ `/ __/ _ `/ _ \/ __/ __/
##      \_,_/\_,_/\__/\_,_/ .__/\__/_/
##                       /_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package dataptr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arnumina/failure"
	"github.com/dolmen-go/jsonptr"
	"gopkg.in/yaml.v3"
)

var (
	// ErrBadType AFAIRE.
	ErrBadType = errors.New("bad type")
	// ErrNotFound AFAIRE.
	ErrNotFound = errors.New("not found")
)

type (
	// DataPtr AFAIRE.
	DataPtr struct {
		value interface{}
	}
)

// New AFAIRE.
func New(value interface{}) *DataPtr {
	return &DataPtr{value: value}
}

// Empty AFAIRE.
func Empty() *DataPtr {
	return New(map[string]interface{}(nil))
}

// FromJSON AFAIRE.
func FromJSON(data []byte) (*DataPtr, error) {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return New(value), nil
}

// FromYAML AFAIRE.
func FromYAML(data []byte) (*DataPtr, error) {
	var value interface{}

	if err := yaml.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return New(value), nil
}

// Value AFAIRE.
func (dp *DataPtr) Value() interface{} {
	return dp.value
}

// Get AFAIRE.
func (dp *DataPtr) Get(keys ...string) (string, *DataPtr, error) {
	ptr := fmt.Sprintf("/%s", strings.Join(append([]string{}, keys...), "/"))

	if ptr == "/" {
		return ptr, dp, nil
	}

	value, err := jsonptr.Get(dp.value, ptr)
	if err != nil {
		if errors.Is(err, jsonptr.ErrProperty) {
			return ptr, nil,
				failure.New(ErrNotFound).
					Set("pointer", ptr).
					Msg("the referenced data does not exist") //////////////////////////////////////////////////////////
		}

		return ptr, nil, err
	}

	return ptr, New(value), nil
}

// MaybeGet AFAIRE.
func (dp *DataPtr) MaybeGet(keys ...string) (*DataPtr, error) {
	_, dptr, err := dp.Get(keys...)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Empty(), nil
		}

		return nil, err
	}

	return dptr, nil
}

// MapString AFAIRE.
func (dp *DataPtr) MapString(keys ...string) (map[string]*DataPtr, error) {
	ptr, dptr, err := dp.Get(keys...)
	if err != nil {
		return nil, err
	}

	ms, ok := dptr.value.(map[string]interface{})
	if !ok {
		return nil,
			failure.New(ErrBadType).
				Set("pointer", ptr).
				Msg("this pointer do not reference a value of type 'map[string]interface{}'") //////////////////////////
	}

	value := make(map[string]*DataPtr)

	for k, v := range ms {
		value[k] = New(v)
	}

	return value, nil
}

// Slice AFAIRE.
func (dp *DataPtr) Slice(keys ...string) ([]*DataPtr, error) {
	ptr, dptr, err := dp.Get(keys...)
	if err != nil {
		return nil, err
	}

	slice, ok := dptr.value.([]interface{})
	if !ok {
		return nil,
			failure.New(ErrBadType).
				Set("pointer", ptr).
				Msg("this pointer do not reference a value of type '[]interface{}'") ///////////////////////////////////
	}

	value := []*DataPtr{}

	for _, v := range slice {
		value = append(value, New(v))
	}

	return value, nil
}

// Bool AFAIRE.
func (dp *DataPtr) Bool(keys ...string) (bool, error) {
	ptr, dptr, err := dp.Get(keys...)
	if err != nil {
		return false, err
	}

	value, ok := dptr.value.(bool)
	if !ok {
		return false,
			failure.New(ErrBadType).
				Set("pointer", ptr).
				Msg("this pointer do not reference a value of type 'bool'") ////////////////////////////////////////////
	}

	return value, nil
}

// DBool AFAIRE.
func (dp *DataPtr) DBool(d bool, keys ...string) (bool, error) {
	value, err := dp.Bool(keys...)
	if errors.Is(err, ErrNotFound) {
		return d, nil
	}

	return value, err
}

// Int AFAIRE.
func (dp *DataPtr) Int(keys ...string) (int, error) {
	ptr, dptr, err := dp.Get(keys...)
	if err != nil {
		return 0, err
	}

	value, ok := dptr.value.(int)
	if !ok {
		return 0,
			failure.New(ErrBadType).
				Set("pointer", ptr).
				Msg("this pointer do not reference a value of type 'int'") /////////////////////////////////////////////
	}

	return value, nil
}

// DInt AFAIRE.
func (dp *DataPtr) DInt(d int, keys ...string) (int, error) {
	value, err := dp.Int(keys...)
	if errors.Is(err, ErrNotFound) {
		return d, nil
	}

	return value, err
}

// String AFAIRE.
func (dp *DataPtr) String(keys ...string) (string, error) {
	ptr, dptr, err := dp.Get(keys...)
	if err != nil {
		return "", err
	}

	value, ok := dptr.value.(string)
	if !ok {
		return "",
			failure.New(ErrBadType).
				Set("pointer", ptr).
				Msg("this pointer do not reference a value of type 'string'") //////////////////////////////////////////
	}

	return value, nil
}

// DString AFAIRE.
func (dp *DataPtr) DString(d string, keys ...string) (string, error) {
	value, err := dp.String(keys...)
	if errors.Is(err, ErrNotFound) {
		return d, nil
	}

	return value, err
}

// Duration AFAIRE.
func (dp *DataPtr) Duration(keys ...string) (time.Duration, error) {
	ptr, dptr, err := dp.Get(keys...)
	if err != nil {
		return 0, err
	}

	s, ok := dptr.value.(string)
	if ok {
		value, err := time.ParseDuration(s)
		if err == nil {
			return value, nil
		}
	}

	return 0,
		failure.New(ErrBadType).
			Set("pointer", ptr).
			Msg("this pointer do not reference a value of type 'time.Duration'") ///////////////////////////////////////
}

// DDuration AFAIRE.
func (dp *DataPtr) DDuration(d time.Duration, keys ...string) (time.Duration, error) {
	value, err := dp.Duration(keys...)
	if errors.Is(err, ErrNotFound) {
		return d, nil
	}

	return value, err
}

/*
######################################################################################################## @(°_°)@ #######
*/
