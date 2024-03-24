/*
Copyright © 2024 buty4649

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package iproute2

import (
	"reflect"
	"testing"
)

func TestUnmarshalInterfacesData(t *testing.T) {
	testCases := []struct {
		desc         string
		input        string
		expected     Interfaces
		expectingErr bool
	}{
		{
			desc:  "Valid input",
			input: `[{"ifindex":1,"ifname":"lo","flags":["LOOPBACK","UP","LOWER_UP"],"mtu":65536,"qdisc":"noqueue","operstate":"UNKNOWN","group":"default","txqlen":1000,"link_type":"loopback","address":"00:00:00:00:00:00","broadcast":"00:00:00:00:00:00","addr_info":[{"family":"inet","local":"127.0.0.1","prefixlen":8,"scope":"host","label":"lo","valid_life_time":4294967295,"preferred_life_time":4294967295},{"family":"inet6","local":"::1","prefixlen":128,"scope":"host","valid_life_time":4294967295,"preferred_life_time":4294967295}]}]`,
			expected: Interfaces{
				{
					Ifindex:   1,
					Ifname:    "lo",
					Flags:     []string{"LOOPBACK", "UP", "LOWER_UP"},
					Mtu:       65536,
					Qdisc:     "noqueue",
					Operstate: OperStateUnkwon,
					Group:     "default",
					Txqlen:    1000,
					LinkType:  "loopback",
					Address:   "00:00:00:00:00:00",
					Broadcast: "00:00:00:00:00:00",
					AddrInfo: []AddressInfo{
						{Family: "inet", Local: "127.0.0.1", Prefixlen: 8, Scope: "host", Label: "lo", ValidLifeTime: 4294967295, PreferredLifeTime: 4294967295},
						{Family: "inet6", Local: "::1", Prefixlen: 128, Scope: "host", Label: "", ValidLifeTime: 4294967295, PreferredLifeTime: 4294967295},
					},
				},
			},
			expectingErr: false,
		},
		{
			desc:         "Invalid input",
			input:        `invalid JSON`,
			expected:     nil,
			expectingErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := unmarshalInterfacesData(tc.input)
			if (err != nil) != tc.expectingErr {
				t.Errorf("unmarshalInterfacesData() error = %v, expectingErr %v", err, tc.expectingErr)
				return
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("unmarshalInterfacesData() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestUnmarshalLinksData(t *testing.T) {
	testCases := []struct {
		desc         string
		input        string
		expected     Links
		expectingErr bool
	}{
		{
			desc:  "Valid input",
			input: `[{"ifindex":1,"ifname":"lo","flags":["LOOPBACK"],"mtu":65536,"qdisc":"noop","operstate":"DOWN","linkmode":"DEFAULT","group":"default","txqlen":1000,"link_type":"loopback","address":"00:00:00:00:00:00","broadcast":"00:00:00:00:00:00"}]`,
			expected: Links{
				{
					Ifindex:   1,
					Ifname:    "lo",
					Flags:     []string{"LOOPBACK"},
					Mtu:       65536,
					Qdisc:     "noop",
					Operstate: OperStateDown,
					Linkmode:  "DEFAULT",
					Group:     "default",
					Txqlen:    1000,
					LinkType:  "loopback",
					Address:   "00:00:00:00:00:00",
					Broadcast: "00:00:00:00:00:00",
				},
			},
			expectingErr: false,
		},
		{
			desc:         "Invalid input",
			input:        `invalid JSON`,
			expected:     nil,
			expectingErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := unmarshalLinksData(tc.input)
			if (err != nil) != tc.expectingErr {
				t.Errorf("unmarshalLinksData() error = %v, expectingErr %v", err, tc.expectingErr)
				return
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("unmarshalLinksData() = %v, want %v", got, tc.expected)
			}
		})
	}
}
