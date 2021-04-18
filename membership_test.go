// +build !windows

package userutil

import (
	"reflect"
	"testing"
)

var tuples = []byte(`0,100,20(staff),12(everyone),61(localaccounts),79(_appserverusr),80(admin)`)

func TestParseMembership(t *testing.T) {

	expected := []Group{
		{Name: "staff", Gid: "20"},
		{Name: "everyone", Gid: "12"},
		{Name: "localaccounts", Gid: "61"},
		{Name: "_appserverusr", Gid: "79"},
		{Name: "admin", Gid: "80"},
	}

	groups := parseMembership(tuples, true)
	if !reflect.DeepEqual(expected, groups) {
		t.Errorf("TestParseMembership failed. Wanted %v, got: %v", expected, groups)
	}
}
