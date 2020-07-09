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

package dataptr_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/arnumina/dataptr"
)

var (
	_testData = map[string]interface{}{
		"x": 785,
		"y": map[string]interface{}{
			"foo":     "bar",
			"ok":      true,
			"null":    nil,
			"timeout": "5m",
			"map": map[string]interface{}{
				"a": "a",
				"b": "b",
				"c": "c",
			},
		},
		"z":     []interface{}{"ok", 2, 5, 7, true, 13, "30s"},
		"slice": []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
)

func TestFromJSON(t *testing.T) {
	_, err := dataptr.FromJSON([]byte(`{"x": 785, "y": {"foo": "bar", "ok": true, "null": null}}`))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestFromYAML(t *testing.T) {
	data := `x: 785
y:
    foo: bar
    ok: true
    null: null`

	_, err := dataptr.FromYAML([]byte(data))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestMapString(t *testing.T) {
	_, err := dataptr.New(_testData).MapString("x")
	if !errors.Is(err, dataptr.ErrBadType) {
		t.Errorf("Error expected: %s", dataptr.ErrBadType)
	}

	ms, err := dataptr.New(_testData).MapString("y", "map")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	for k, v := range ms {
		if !reflect.DeepEqual(k, v.Value()) {
			t.Errorf("got: %v, want %v", v.Value(), k)
		}
	}
}

func TestSlice(t *testing.T) {
	_, err := dataptr.New(_testData).Slice("y", "foo")
	if !errors.Is(err, dataptr.ErrBadType) {
		t.Errorf("Error expected => got: %v want: %v", err, dataptr.ErrBadType)
	}

	s, err := dataptr.New(_testData).Slice("slice")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	for i, v := range s {
		if !reflect.DeepEqual(i, v.Value()) {
			t.Errorf("got: %v, want %v", v.Value(), i)
		}
	}
}

func TestBool(t *testing.T) {
	_, err := dataptr.New(_testData).Bool("x")
	if !errors.Is(err, dataptr.ErrBadType) {
		t.Errorf("Error expected: %s", dataptr.ErrBadType)
	}
}

func TestDBool(t *testing.T) {
	dp := dataptr.New(_testData)

	tests := []struct {
		d    bool
		want bool
		keys []string
	}{
		{false, true, []string{"y", "ok"}},
		{false, true, []string{"z", "4"}},
		{true, true, []string{"y", "ko"}},
		{false, false, []string{"y", "ko"}},
	}

	for i, tt := range tests {
		got, err := dp.DBool(tt.d, tt.keys...)
		if err != nil {
			t.Errorf("[%02d] => error: %s", i, err)
		} else if got != tt.want {
			t.Errorf("[%02d] => got: %v, want %v", i, got, tt.want)
		}
	}
}

func TestInt(t *testing.T) {
	_, err := dataptr.New(_testData).Int("y", "foo")
	if !errors.Is(err, dataptr.ErrBadType) {
		t.Errorf("Error expected: %s", dataptr.ErrBadType)
	}
}

func TestDInt(t *testing.T) {
	dp := dataptr.New(_testData)

	tests := []struct {
		d    int
		want int
		keys []string
	}{
		{0, 785, []string{"x"}},
		{0, 13, []string{"z", "5"}},
		{123, 123, []string{"y", "ko"}},
	}

	for i, tt := range tests {
		got, err := dp.DInt(tt.d, tt.keys...)
		if err != nil {
			t.Errorf("[%02d] => error: %s", i, err)
		} else if got != tt.want {
			t.Errorf("[%02d] => got: %v, want %v", i, got, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	_, err := dataptr.New(_testData).String("z")
	if !errors.Is(err, dataptr.ErrBadType) {
		t.Errorf("Error expected: %s", dataptr.ErrBadType)
	}
}

func TestDString(t *testing.T) {
	dp := dataptr.New(_testData)

	tests := []struct {
		d    string
		want string
		keys []string
	}{
		{"", "bar", []string{"y", "foo"}},
		{"", "ok", []string{"z", "0"}},
		{"default", "default", []string{"ko"}},
	}

	for i, tt := range tests {
		got, err := dp.DString(tt.d, tt.keys...)
		if err != nil {
			t.Errorf("[%02d] => error: %s", i, err)
		} else if got != tt.want {
			t.Errorf("[%02d] => got: %v, want %v", i, got, tt.want)
		}
	}
}

func TestDuration(t *testing.T) {
	_, err := dataptr.New(_testData).Duration("y", "ok")
	if !errors.Is(err, dataptr.ErrBadType) {
		t.Errorf("Error expected: %s", dataptr.ErrBadType)
	}
}

func TestDDuration(t *testing.T) {
	dp := dataptr.New(_testData)

	tests := []struct {
		d    time.Duration
		want time.Duration
		keys []string
	}{
		{0, 5 * time.Minute, []string{"y", "timeout"}},
		{0, 30 * time.Second, []string{"z", "6"}},
		{20 * time.Millisecond, 20 * time.Millisecond, []string{"ko", "ok"}},
	}

	for i, tt := range tests {
		got, err := dp.DDuration(tt.d, tt.keys...)
		if err != nil {
			t.Errorf("[%02d] => error: %s", i, err)
		} else if got != tt.want {
			t.Errorf("[%02d] => got: %v, want %v", i, got, tt.want)
		}
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
