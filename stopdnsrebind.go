package stopdnsrebind

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Stopdnsrebind struct {
	Next      plugin.Handler
	AllowList []string
	DenyList  []net.IPNet
}

// ServeDNS implements the plugin.Handler interface.
func (a Stopdnsrebind) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	state := request.Request{W: w, Req: r}

	//ignore if on the allow list
	for _, allowed := range a.AllowList {
		if allowed == state.QName() {
			return plugin.NextOrFailure(a.Name(), a.Next, ctx, w, r)
		}
	}

	nw := nonwriter.New(w)

	rcode, err := plugin.NextOrFailure(a.Name(), a.Next, ctx, nw, r)

	if err != nil {
		return rcode, err
	}

	for _, ans := range nw.Msg.Answer {
		var ip net.IP

		switch ans.Header().Rrtype {
		case dns.TypeA:
			ip = ans.(*dns.A).A
		case dns.TypeAAAA:
			ip = ans.(*dns.AAAA).AAAA
		default:
			//we only care about A and AAA
			continue
		}

		/*
			🚀 Default blocking rules:

			🔒 Loopback Addresses: 127.0.0.1/8
			🔒 Private Addresses:
				- 10.0.0.0/8
				- 172.16.0.0/12
				- 192.168.0.0/16
			🔒 Link Local Addresses: 169.254.0.0/16
			🔒 Unspecified: 0.0.0.0
			🔒 Interface Local Multicast: 224.0.0.0/24
			🔒 DenyList: Add your entries in the plugin configuration

			// Keeping the network secure!
		*/

		if !ip.IsGlobalUnicast() || ip.IsInterfaceLocalMulticast() ||
			ip.IsPrivate() || shouldDeny(ip, a.DenyList) {
			m := new(dns.Msg)
			m.SetRcode(r, dns.RcodeRefused)
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}
	}

	w.WriteMsg(nw.Msg)

	return 0, nil
}

func shouldDeny(ip net.IP, denyList []net.IPNet) bool {
	for _, ipNetDenied := range denyList {
		if ipNetDenied.Contains(ip) {
			return true
		}
	}
	return false
}

// Name implements the Handler interface.
func (a Stopdnsrebind) Name() string { return "stopdnsrebind" }
