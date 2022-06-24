/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package common

import "testing"

func TestParseRawName(t *testing.T) {
	cases := []struct {
		raw      string
		node     string
		domain   string
		expected string
	}{
		{
			raw:      "hack",
			node:     "alice",
			domain:   "nivenly.com",
			expected: "hack@alice@nivenly.com",
		},
		{
			raw:      "hack@nivenly.com",
			node:     "alice",
			domain:   "nivenly.com",
			expected: "hack@alice@nivenly.com",
		},
		{
			raw:      "hack@alice",
			node:     "alice",
			domain:   "nivenly.com",
			expected: "hack@alice@nivenly.com",
		},
		{
			raw:      "hack@alice@nivenly.com",
			node:     "alice",
			domain:   "nivenly.com",
			expected: "hack@alice@nivenly.com",
		},
	}
	for _, testCase := range cases {
		actual := ParseRawName(testCase.raw, testCase.node, testCase.domain)
		if actual != testCase.expected {
			t.Errorf("Actual: %s, Expected: %s", actual, testCase.expected)
		}
	}
}
