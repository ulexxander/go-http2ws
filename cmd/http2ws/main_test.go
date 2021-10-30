package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type headers map[string]string

func (h headers) String() string {
	b, err := json.Marshal(h)
	if err != nil {
		panic(fmt.Errorf("stringifying headers: %v", err))
	}
	return string(b)
}

func TestParseHeaders(t *testing.T) {
	tt := []struct {
		arg     string
		headers headers
		err     bool
	}{
		{
			arg:     "",
			headers: nil,
		},
		{
			arg: "Content-Type:application/json",
			headers: headers{
				"Content-Type": "application/json",
			},
		},
		{
			arg: "Authorization=123",
			err: true,
		},
		{
			arg: "Host:developer.mozilla.org,Connection:keep-alive",
			headers: headers{
				"Host":       "developer.mozilla.org",
				"Connection": "keep-alive",
			},
		},
		{
			arg: `Referer:https\://developer.mozilla.org/testpage.html`,
			headers: headers{
				"Referer": "https://developer.mozilla.org/testpage.html",
			},
		},
	}
	for _, tc := range tt {
		name := tc.arg
		if tc.arg == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			h, err := parseHeaders(tc.arg)
			if tc.err && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.err && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(headers(h), tc.headers) {
				t.Fatalf("expected headers %s, got: %s", tc.headers, headers(h))
			}
		})
	}
}

func TestSplitEscaped(t *testing.T) {
	tt := []struct {
		s        string
		sep, esc rune
		splitted []string
	}{
		{
			s:        "",
			sep:      ',',
			esc:      '\\',
			splitted: []string{},
		},
		{
			s:        "1,2,3",
			sep:      ',',
			esc:      '\\',
			splitted: []string{"1", "2", "3"},
		},
		{
			s:        "1,2\\,3",
			sep:      ',',
			esc:      '\\',
			splitted: []string{"1", "2,3"},
		},
	}
	for _, tc := range tt {
		name := fmt.Sprintf("s=%s sep=%s esc=%s", tc.s, string(tc.sep), string(tc.esc))
		t.Run(name, func(t *testing.T) {
			result := splitEscaped(tc.s, tc.sep, tc.esc)
			if !reflect.DeepEqual(result, tc.splitted) {
				t.Fatalf("expected %v, got: %v", tc.splitted, result)
			}
		})
	}
}
