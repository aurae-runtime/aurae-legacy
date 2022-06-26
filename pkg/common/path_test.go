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

import (
	"testing"
)

func TestPath(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{
			path:     "    boops",
			expected: "/boops",
		},
		{
			path:     "    beeps //boops",
			expected: "/beeps/boops",
		},
		{
			path:     "   beeps    ",
			expected: "/beeps",
		},
		{
			path:     "    beeps /\\/\\///////boops",
			expected: "/beeps/boops",
		},
		{
			path:     "    beeps /\\/\\///////boops",
			expected: "/beeps/boops",
		},
		{
			path:     "a/b/\\   c/  d /e / f   ",
			expected: "/a/b/c/d/e/f",
		},
		{
			path:     "*",
			expected: "/*",
		},
		{
			path:     "",
			expected: "/",
		},
		{
			path:     " ",
			expected: "/",
		},
		{
			path:     "    \t\n\r",
			expected: "/",
		},
	}
	for _, c := range cases {
		actual := Path(c.path)
		if actual != c.expected {
			t.Errorf("Expected: %s, Actual: %s", c.expected, actual)
		} else {
			//t.Logf("Expected: %s, Actual: %s", c.expected, actual)
		}
	}
}
