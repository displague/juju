// Copyright 2012, 2013, 2014 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

// +build windows

package tools_test

var toolsTestsPlatformSpecific = []struct {
	in          string
	expected    string
	expectedErr error
}{{
	in:          "C:/home/foo",
	expected:    "file://\\\\localhost\\C$/home/foo/tools",
	expectedErr: nil,
}, {
	in:          "C:/home/foo/tools",
	expected:    "file://\\\\localhost\\C$/home/foo/tools",
	expectedErr: nil,
}}
