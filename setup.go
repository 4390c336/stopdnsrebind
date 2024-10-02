package stopdnsrebind

import (
	"net"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

func init() { plugin.Register("stopdnsrebind", setup) }

func setup(c *caddy.Controller) error {
	allowList, denyList, err := parse(c)

	//parsing err
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Stopdnsrebind{Next: next, AllowList: allowList, DenyList: denyList}
	})

	return nil
}

func parse(c *caddy.Controller) ([]string, []net.IPNet, error) {
	allowList := []string{}
	denyList := []net.IPNet{}
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "allow":
				for _, d := range c.RemainingArgs() {
					_, valid := dns.IsDomainName(d)
					if !valid {
						return nil, nil, plugin.Error("stopdnsrebind", c.Errf("%s is not a valid domain", d))
					}

					allowList = append(allowList, d)
				}
			case "deny":
				for _, cidr := range c.RemainingArgs() {
					_, ipNet, err := net.ParseCIDR(cidr)
					if err != nil {
						return nil, nil, plugin.Error("stopdnsrebind", c.Errf("%s is not a valid cidr", cidr))
					}

					denyList = append(denyList, *ipNet)
				}
			default:
				return nil, nil, plugin.Error("stopdnsrebind", c.Err("only allow and deny operations are supported"))
			}
		}
	}

	return allowList, denyList, nil
}
