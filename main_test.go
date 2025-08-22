package main

import (
	"testing"
)

func TestIsTableAllowed(t *testing.T) {
	tests := []struct {
		table    string
		expected bool
	}{
		{"GTPL_108_gT_40E_P_S7_200_Germany", true},
		{"GTPL_109_gT_40E_P_S7_200_Germany", true},
		{"invalid_table", false},
		{"", false},
	}

	for _, test := range tests {
		result := isTableAllowed(test.table)
		if result != test.expected {
			t.Errorf("isTableAllowed(%s) = %v, expected %v", test.table, result, test.expected)
		}
	}
}

func TestLooksLikeFaultKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"fault_code", true},
		{"overheat_protection", true},
		{"temperature", false},
		{"normal_value", false},
	}

	for _, test := range tests {
		result := looksLikeFaultKey(test.key)
		if result != test.expected {
			t.Errorf("looksLikeFaultKey(%s) = %v, expected %v", test.key, result, test.expected)
		}
	}
}

func TestIsTrueish(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{true, true},
		{false, false},
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{1, true},
		{0, false},
		{nil, false},
	}

	for _, test := range tests {
		result := isTrueish(test.value)
		if result != test.expected {
			t.Errorf("isTrueish(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestToNum(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected interface{}
	}{
		{"123.45", 123.45},
		{"", ""},
		{nil, ""},
		{42, 42},
		{3.14, 3.14},
		{"invalid", ""},
	}

	for _, test := range tests {
		result := toNum(test.value)
		if result != test.expected {
			t.Errorf("toNum(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}
