package stopdnsrebind

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("stopdnsrebind", setup) }

func setup(c *caddy.Controller) error {
	allowList := []string{}

	for c.Next() {
		for c.NextBlock() {
			if c.Val() != "allow" {
				return plugin.Error("stopdnsrebind", c.Err("only allow operation is supported"))
			}
			allowList = append(allowList, c.RemainingArgs()...)
		}
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Stopdnsrebind{Next: next, AllowList: allowList}
	})

	return nil
}
