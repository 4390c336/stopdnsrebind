package stopdnsrebind

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

func init() { plugin.Register("stopdnsrebind", setup) }

func setup(c *caddy.Controller) error {
	allowList, err := parse(c)

	//parsing err
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Stopdnsrebind{Next: next, AllowList: allowList}
	})

	return nil
}

func parse(c *caddy.Controller) ([]string, error) {
	allowList := []string{}
	for c.Next() {
		for c.NextBlock() {
			if c.Val() != "allow" {
				return nil, plugin.Error("stopdnsrebind", c.Err("only allow operation is supported"))
			}

			for _, d := range c.RemainingArgs() {
				_, valid := dns.IsDomainName(d)
				if !valid {
					return nil, plugin.Error("stopdnsrebind", c.Errf("%s is not a valid domain", d))
				}

				allowList = append(allowList, d)
			}
		}
	}

	return allowList, nil
}
