// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package network_test

import (
	"sort"

	"github.com/juju/errors"
	"github.com/juju/juju/core/network"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type subnetSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&subnetSuite{})

func (*subnetSuite) TestFindSubnetIDsForAZ(c *gc.C) {
	testCases := []struct {
		name           string
		zoneName       string
		subnetsToZones map[network.Id][]string
		expected       []network.Id
		expectedErr    func(error) bool
	}{
		{
			name:           "empty",
			zoneName:       "",
			subnetsToZones: make(map[network.Id][]string),
			expected:       make([]network.Id, 0),
			expectedErr:    errors.IsNotFound,
		},
		{
			name:     "no match",
			zoneName: "fuzz",
			subnetsToZones: map[network.Id][]string{
				"bar": {"foo", "baz"},
			},
			expected:    make([]network.Id, 0),
			expectedErr: errors.IsNotFound,
		},
		{
			name:     "match",
			zoneName: "foo",
			subnetsToZones: map[network.Id][]string{
				"bar": {"foo", "baz"},
			},
			expected: []network.Id{"bar"},
		},
		{
			name:     "multi-match",
			zoneName: "foo",
			subnetsToZones: map[network.Id][]string{
				"bar":   {"foo", "baz"},
				"other": {"aaa", "foo", "xxx"},
			},
			expected: []network.Id{"bar", "other"},
		},
	}

	for i, t := range testCases {
		c.Logf("test %d: %s", i, t.name)

		res, err := network.FindSubnetIDsForAvailabilityZone(t.zoneName, t.subnetsToZones)
		if t.expectedErr != nil {
			c.Check(t.expectedErr(err), jc.IsTrue)
		} else {
			c.Assert(err, gc.IsNil)
			c.Check(res, gc.DeepEquals, t.expected)
		}
	}
}

func (subnetSuite) TestSubnetSetSize(c *gc.C) {
	// Empty sets are empty.
	s := network.MakeSubnetSet()
	c.Assert(s.Size(), gc.Equals, 0)

	// Size returns number of unique values.
	s = network.MakeSubnetSet("foo", "foo", "bar")
	c.Assert(s.Size(), gc.Equals, 2)
}

func (subnetSuite) TestSubnetSetEmpty(c *gc.C) {
	s := network.MakeSubnetSet()
	assertValues(c, s)
}

func (subnetSuite) TestSubnetSetInitialValues(c *gc.C) {
	values := []network.Id{"foo", "bar", "baz"}
	s := network.MakeSubnetSet(values...)
	assertValues(c, s, values...)
}

func (subnetSuite) TestSubnetSetIsEmpty(c *gc.C) {
	// Empty sets are empty.
	s := network.MakeSubnetSet()
	c.Assert(s.IsEmpty(), gc.Equals, true)

	// Non-empty sets are not empty.
	s = network.MakeSubnetSet("foo")
	c.Assert(s.IsEmpty(), gc.Equals, false)
}

func (subnetSuite) TestSubnetSetAdd(c *gc.C) {
	s := network.MakeSubnetSet()
	s.Add("foo")
	s.Add("foo")
	s.Add("bar")
	assertValues(c, s, "foo", "bar")
}

func (subnetSuite) TestSubnetSetContains(c *gc.C) {
	s := network.MakeSubnetSet("foo", "bar")
	c.Assert(s.Contains("foo"), gc.Equals, true)
	c.Assert(s.Contains("bar"), gc.Equals, true)
	c.Assert(s.Contains("baz"), gc.Equals, false)
}

// Helper methods for the tests.
func assertValues(c *gc.C, s network.SubnetSet, expected ...network.Id) {
	values := s.Values()

	// Expect an empty slice, not a nil slice for values.
	if expected == nil {
		expected = make([]network.Id, 0)
	}

	sort.Slice(expected, func(i, j int) bool {
		return expected[i] < expected[j]
	})
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	c.Assert(values, gc.DeepEquals, expected)
	c.Assert(s.Size(), gc.Equals, len(expected))

	// Check the sorted values too.
	sorted := s.SortedValues()
	c.Assert(sorted, gc.DeepEquals, expected)
}
