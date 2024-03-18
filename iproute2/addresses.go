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

import "encoding/json"

type AddressInfo struct {
	Family            string `json:"family"`
	Local             string `json:"local"`
	Prefixlen         int    `json:"prefixlen"`
	Scope             string `json:"scope"`
	Label             string `json:"label"`
	ValidLifeTime     uint64 `json:"valid_life_time"`
	PreferredLifeTime uint64 `json:"preferred_life_time"`
}

type InterfaceInfo struct {
	Ifindex   int           `json:"ifindex"`
	Ifname    string        `json:"ifname"`
	Flags     []string      `json:"flags"`
	Mtu       int           `json:"mtu"`
	Qdisc     string        `json:"qdisc"`
	Operstate string        `json:"operstate"`
	Group     string        `json:"group"`
	Txqlen    int           `json:"txqlen"`
	LinkType  string        `json:"link_type"`
	Address   string        `json:"address"`
	Broadcast string        `json:"broadcast"`
	AddrInfo  []AddressInfo `json:"addr_info"`
}

type Addresses []InterfaceInfo

func (i *Iproute2) ListAddresses() (*Addresses, error) {
	data, err := i.executeWithStdout("-json", "address", "show")
	if err != nil {
		return nil, err
	}

	var addresses Addresses
	err = json.Unmarshal([]byte(data), &addresses)
	if err != nil {
		return nil, err
	}

	return &addresses, nil
}
