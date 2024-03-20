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

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Netns map[string]Netns `yaml:"netns"`
}

type Netns struct {
	Ethernets    map[string]Ethernet   `yaml:"ethernets,omitempty"`
	DummyDevices map[string]Ethernet   `yaml:"dummy-devices,omitempty"`
	VethDevices  map[string]VethDevice `yaml:"veth-devices,omitempty"`
}

type Ethernet struct {
	Addresses []string `yaml:"addresses"`
	Routes    []Route  `yaml:"routes,omitempty"`
}

type VethDevice struct {
	Addresses []string `yaml:"addresses"`
	Routes    []Route  `yaml:"routes,omitempty"`
	Peer      Peer     `yaml:"peer"`
}

type Peer struct {
	Name      string   `yaml:"name"`
	Netns     string   `yaml:"netns,omitempty"`
	Addresses []string `yaml:"addresses"`
	Routes    []Route  `yaml:"routes,omitempty"`
}

type Route struct {
	To  string `yaml:"to"`
	Via string `yaml:"via"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
