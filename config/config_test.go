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
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadYamlFiles(t *testing.T) {
	wd, _ := os.Getwd()
	testdataDir := filepath.Join(wd, "..", "testdata", "config")

	expected := &Config{
		Netns: map[string]Netns{
			"sample1": {
				Loopback: Ethernet{
					Addresses: []string{"127.0.0.53/8"},
					Routes: []Route{{
						To:  "10.10.0.0/24",
						Via: "127.0.0.53",
					}},
				},
				Ethernets: map[string]Ethernet{
					"eth0": {
						Addresses: []string{
							"192.168.0.1/24",
							"2001:db8:beaf:cafe::1/112",
						},
						Routes: []Route{{
							To:  "default",
							Via: "192.168.0.254",
						}},
					},
					"eth1": {
						Addresses: []string{
							"192.168.1.1/24",
							"10.0.0.1/24",
						},
					},
				},
				DummyDevices: map[string]Ethernet{
					"dummy0": {
						Addresses: []string{"192.168.10.1/24"},
						Routes: []Route{{
							To:  "192.168.11.0/24",
							Via: "192.168.10.254",
						}},
					},
				},
				VethDevices: map[string]VethDevice{
					"veth0": {
						Addresses: []string{"192.168.20.1/24"},
						Routes: []Route{{
							To:  "192.168.21.0/24",
							Via: "192.168.20.254",
						}},
						Peer: Peer{
							Name:      "veth0-peer",
							Netns:     "sample2",
							Addresses: []string{"192.168.20.2/24"},
							Routes: []Route{{
								To:  "192.168.21.0/24",
								Via: "192.168.20.2",
							}},
						},
					},
				},
				PostScript: "echo 'Hello, World!'\n",
			},
			"sample2": {
				Ethernets: map[string]Ethernet{
					"eth2": {
						Addresses: []string{"172.16.0.1/24"},
					},
				},
			},
		},
	}

	result, err := LoadYamlFiles(testdataDir)
	if err != nil {
		t.Fatalf("LoadYamlFiles returned an error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
