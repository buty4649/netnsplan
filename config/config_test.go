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
	"reflect"
	"testing"
)

const testYAML = `
netns:
  netns1:
    ethernets:
      eth0:
        addresses:
          - "192.168.1.1/24"
        routes:
          - to: "0.0.0.0/0"
            via: "192.168.1.254"
    dummy-devices:
      dummy0:
        addresses:
          - "10.0.0.1/8"
    veth-devices:
      veth0:
        addresses:
          - "10.1.0.1/24"
        peer:
          name: "veth0-peer"
          netns: "netns2"
          addresses:
            - "10.1.0.2/24"
    post-script: |
        echo 'Hello'
        echo 'World!'
`

func TestLoadConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testYAML)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	expected := &Config{
		Netns: map[string]Netns{
			"netns1": {
				Ethernets: map[string]Ethernet{
					"eth0": {
						Addresses: []string{"192.168.1.1/24"},
						Routes: []Route{
							{To: "0.0.0.0/0", Via: "192.168.1.254"},
						},
					},
				},
				DummyDevices: map[string]Ethernet{
					"dummy0": {
						Addresses: []string{"10.0.0.1/8"},
					},
				},
				VethDevices: map[string]VethDevice{
					"veth0": {
						Addresses: []string{"10.1.0.1/24"},
						Peer: Peer{
							Name:      "veth0-peer",
							Netns:     "netns2",
							Addresses: []string{"10.1.0.2/24"},
						},
					},
				},
				PostScript: "echo 'Hello'\necho 'World!'\n",
			},
		},
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Config does not match expected\nGot: %#v\nWant: %#v", config, expected)
	}
}
