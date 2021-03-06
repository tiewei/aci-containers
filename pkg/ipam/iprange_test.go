// Copyright 2017 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipam

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type range2cidrTest struct {
	input  []string
	output []string
}

var range2cidrTests = []range2cidrTest{
	{
		[]string{"0.0.0.0", "0.0.0.0"},
		[]string{"0.0.0.0/32"},
	},
	{
		[]string{"0.0.0.0", "0.0.0.1"},
		[]string{"0.0.0.0/31"},
	},
	{
		[]string{"0.0.0.1", "0.0.0.2"},
		[]string{"0.0.0.1/32", "0.0.0.2/32"},
	},
	{
		[]string{"255.255.255.255", "255.255.255.255"},
		[]string{"255.255.255.255/32"},
	},
	{
		[]string{"0.0.0.254", "0.0.1.0"},
		[]string{"0.0.0.254/31", "0.0.1.0/32"},
	},
	{
		[]string{"0.0.0.0", "255.255.255.255"},
		[]string{"0.0.0.0/0"},
	},
	{
		[]string{"1.2.3.4", "5.6.7.8"},
		[]string{"1.2.3.4/30",
			"1.2.3.8/29",
			"1.2.3.16/28",
			"1.2.3.32/27",
			"1.2.3.64/26",
			"1.2.3.128/25",
			"1.2.4.0/22",
			"1.2.8.0/21",
			"1.2.16.0/20",
			"1.2.32.0/19",
			"1.2.64.0/18",
			"1.2.128.0/17",
			"1.3.0.0/16",
			"1.4.0.0/14",
			"1.8.0.0/13",
			"1.16.0.0/12",
			"1.32.0.0/11",
			"1.64.0.0/10",
			"1.128.0.0/9",
			"2.0.0.0/7",
			"4.0.0.0/8",
			"5.0.0.0/14",
			"5.4.0.0/15",
			"5.6.0.0/22",
			"5.6.4.0/23",
			"5.6.6.0/24",
			"5.6.7.0/29",
			"5.6.7.8/32"},
	},
	{
		[]string{"0.0.0.1", "255.255.255.254"},
		[]string{"0.0.0.1/32",
			"0.0.0.2/31",
			"0.0.0.4/30",
			"0.0.0.8/29",
			"0.0.0.16/28",
			"0.0.0.32/27",
			"0.0.0.64/26",
			"0.0.0.128/25",
			"0.0.1.0/24",
			"0.0.2.0/23",
			"0.0.4.0/22",
			"0.0.8.0/21",
			"0.0.16.0/20",
			"0.0.32.0/19",
			"0.0.64.0/18",
			"0.0.128.0/17",
			"0.1.0.0/16",
			"0.2.0.0/15",
			"0.4.0.0/14",
			"0.8.0.0/13",
			"0.16.0.0/12",
			"0.32.0.0/11",
			"0.64.0.0/10",
			"0.128.0.0/9",
			"1.0.0.0/8",
			"2.0.0.0/7",
			"4.0.0.0/6",
			"8.0.0.0/5",
			"16.0.0.0/4",
			"32.0.0.0/3",
			"64.0.0.0/2",
			"128.0.0.0/2",
			"192.0.0.0/3",
			"224.0.0.0/4",
			"240.0.0.0/5",
			"248.0.0.0/6",
			"252.0.0.0/7",
			"254.0.0.0/8",
			"255.0.0.0/9",
			"255.128.0.0/10",
			"255.192.0.0/11",
			"255.224.0.0/12",
			"255.240.0.0/13",
			"255.248.0.0/14",
			"255.252.0.0/15",
			"255.254.0.0/16",
			"255.255.0.0/17",
			"255.255.128.0/18",
			"255.255.192.0/19",
			"255.255.224.0/20",
			"255.255.240.0/21",
			"255.255.248.0/22",
			"255.255.252.0/23",
			"255.255.254.0/24",
			"255.255.255.0/25",
			"255.255.255.128/26",
			"255.255.255.192/27",
			"255.255.255.224/28",
			"255.255.255.240/29",
			"255.255.255.248/30",
			"255.255.255.252/31",
			"255.255.255.254/32"},
	},
	{
		[]string{"::", "::"},
		[]string{"::/128"},
	},
	{
		[]string{"::", "::0001"},
		[]string{"::/127"},
	},
	{
		[]string{"::0001", "::0002"},
		[]string{"::1/128", "::2/128"},
	},
	{
		[]string{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"},
		[]string{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128"},
	},
	{
		[]string{"::00fe", "::0100"},
		[]string{"::fe/127", "::100/128"},
	},
	{
		[]string{"::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"},
		[]string{"::/0"},
	},
	{
		[]string{"::0102:0304", "::0506:0708"},
		[]string{"::102:304/126",
			"::102:308/125",
			"::102:310/124",
			"::102:320/123",
			"::102:340/122",
			"::102:380/121",
			"::102:400/118",
			"::102:800/117",
			"::102:1000/116",
			"::102:2000/115",
			"::102:4000/114",
			"::102:8000/113",
			"::103:0/112",
			"::104:0/110",
			"::108:0/109",
			"::110:0/108",
			"::120:0/107",
			"::140:0/106",
			"::180:0/105",
			"::200:0/103",
			"::400:0/104",
			"::500:0/110",
			"::504:0/111",
			"::506:0/118",
			"::506:400/119",
			"::506:600/120",
			"::506:700/125",
			"::506:708/128"},
	},
}

func TestRange2Cidr(t *testing.T) {
	for _, rt := range range2cidrTests {
		in := IpRange{net.ParseIP(rt.input[0]), net.ParseIP(rt.input[1])}
		var out []string
		for _, n := range Range2Cidr(in.Start, in.End) {
			out = append(out, n.String())
		}
		assert.Equal(t, rt.output, out,
			fmt.Sprintf("%s - %s", rt.input[0], rt.input[1]))
	}
}
