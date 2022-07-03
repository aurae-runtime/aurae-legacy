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

package name

import "testing"

func TestNames(t *testing.T) {

	cases := []struct {
		message          string
		raw              string
		domainExpected   string
		hostExpected     string
		subExpected      string
		fullhostExpected string
		serviceExpected  string
	}{
		{
			message:          "Testing basic example",
			raw:              "hack@alice@nivenly.com",
			domainExpected:   "nivenly.com",
			hostExpected:     "alice",
			subExpected:      "hack",
			fullhostExpected: "hack@alice@nivenly.com",
			serviceExpected:  "hack@alice@nivenly.com",
		},
		{
			message:          "Testing single string example",
			raw:              "hack",
			domainExpected:   "",
			hostExpected:     "hack",
			subExpected:      "",
			fullhostExpected: "@hack@",
			serviceExpected:  "hack",
		},
		{
			message:          "Testing double string experiment", // :)
			raw:              "alice@nivenly.com",
			domainExpected:   "nivenly.com",
			hostExpected:     "alice",
			subExpected:      "",
			fullhostExpected: "@alice@nivenly.com",
			serviceExpected:  "alice@nivenly.com",
		},
	}
	for _, c := range cases {
		name := New(c.raw)
		if name.sub != c.subExpected {
			t.Logf(c.message)
			t.Errorf("Failed [sub]. Expected: %s, Actual: %s", c.subExpected, name.sub)
		}
		if name.host != c.hostExpected {
			t.Logf(c.message)
			t.Errorf("Failed [host]. Expected: %s, Actual: %s", c.hostExpected, name.host)
		}
		if name.domain != c.domainExpected {
			t.Logf(c.message)
			t.Errorf("Failed [domain]. Expected: %s, Actual: %s", c.domainExpected, name.domain)
		}
		if name.service != c.serviceExpected {
			t.Logf(c.message)
			t.Errorf("Failed [service]. Expected: %s, Actual: %s", c.serviceExpected, name.service)
		}
		if name.Host() != c.fullhostExpected {
			t.Logf(c.message)
			t.Errorf("Failed [Host()]. Expected: %s, Actual: %s", c.fullhostExpected, name.Host())
		}
	}
}
