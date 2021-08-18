package stopdnsrebind

import (
	"context"
	"net"
	"reflect"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/miekg/dns"
)

type Stopdnsrebind struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (a Stopdnsrebind) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	nw := nonwriter.New(w)

	rcode, err := plugin.NextOrFailure(a.Name(), a.Next, ctx, nw, r)

	if err != nil {
		return rcode, err
	}

	for _, ans := range nw.Msg.Answer {
		//we only care about A and AAAA types
		if ans.Header().Rrtype != dns.TypeA && ans.Header().Rrtype != dns.TypeAAAA {
			continue
		}

		//get te field directly and convert it to net.IP
		ip := net.IP(reflect.ValueOf(ans).Elem().Field(1).Bytes())

		//check if private
		if isPrivate(ip) {
			m := new(dns.Msg)
			m.SetRcode(r, dns.RcodeRefused)
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}
	}

	w.WriteMsg(nw.Msg)

	return 0, nil
}

var reservedIPv4Nets = []net.IPNet{
	String2IPNet("192.0.2.1/24"),
	String2IPNet("10.0.0.1/8"),
	String2IPNet("127.0.0.1/8"),
	String2IPNet("169.254.0.0/16"),
}

func String2IPNet(cidr string) net.IPNet {
	_, ipnet, _ := net.ParseCIDR(cidr)
	return *ipnet
}

func isPrivate(ip net.IP) bool {
	if ip.To4() == nil && !ip.IsGlobalUnicast() {
		return true
	}

	for _, privnet := range reservedIPv4Nets {
		if privnet.Contains(ip) {
			return true
		}
	}
	return false
}

// Name implements the Handler interface.
func (a Stopdnsrebind) Name() string { return "stopdnsrebind" }
