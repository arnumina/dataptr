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

	"github.com/arnumina/failure"
	"github.com/dolmen-go/jsonptr"
	"gopkg.in/yaml.v3"
)

var (
	// ErrNotFound AFAIRE.
	ErrNotFound = errors.New("not found")
)

type (
	// DataPtr AFAIRE.
	DataPtr struct {
		data interface{}
	}
)

// New AFAIRE.
func New(data interface{}) *DataPtr {
	return &DataPtr{data: data}
}

// Empty AFAIRE.
func Empty() *DataPtr {
	return &DataPtr{data: map[string]interface{}(nil)}
}

// FromJSON AFAIRE.
func FromJSON(data []byte) (*DataPtr, error) {
	dp := &DataPtr{}
	if err := json.Unmarshal(data, &dp.data); err != nil {
		return nil, err
	}

	return dp, nil
}

// FromYAML AFAIRE.
func FromYAML(data []byte) (*DataPtr, error) {
	dp := &DataPtr{}
	if err := yaml.Unmarshal(data, &dp.data); err != nil {
		return nil, err
	}

	return dp, nil
}

// Data AFAIRE.
func (dp *DataPtr) Data() interface{} {
	return dp.data
}

// Get AFAIRE.
func (dp *DataPtr) Get(keys ...string) (string, *DataPtr, error) {
	ptr := fmt.Sprintf("/%s", strings.Join(append([]string{}, keys...), "/"))

	if ptr == "/" {
		return ptr, dp, nil
	}

	d, err := jsonptr.Get(dp.data, ptr)
	if err != nil {
		if errors.Is(err, jsonptr.ErrProperty) {
			return ptr, nil,
				failure.New(ErrNotFound).
					Set("pointer", ptr).
					Msg("this data does not exist") ////////////////////////////////////////////////////////////////////
		}

		return ptr, nil, err
	}

	return ptr, &DataPtr{data: d}, nil
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

/*
######################################################################################################## @(°_°)@ #######
*/
