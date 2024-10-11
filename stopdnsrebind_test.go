package stopdnsrebind

import (
	"context"
	"net"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

// testHandler
type testHandler struct {
	Response *test.Case
	Next     plugin.Handler
}

type testcase struct {
	Expected int
	test     test.Case
	config   string
}

func (t *testHandler) Name() string { return "test-handler" }

func (t *testHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	d := new(dns.Msg)
	d.SetReply(r)
	if t.Response != nil {
		d.Answer = t.Response.Answer
		d.Rcode = t.Response.Rcode
	}
	w.WriteMsg(d)
	return 0, nil
}

func TestBlockingResponse(t *testing.T) {
	var tests = []testcase{
		{
			Expected: dns.RcodeSuccess,
			test: test.Case{
				Answer: []dns.RR{test.A("example.org. 0 IN A 1.1.1.1")},
				Qname:  "example.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 169.254.169.254")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 10.0.0.1")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 172.16.0.1")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 192.168.0.1")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 0.0.0.0")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.org. 0 IN A 224.0.0.0")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 127.0.0.1")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.A("example.refused.org. 0 IN A 192.0.2.1")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeA,
			},
			config: "yep",
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("example.refused.org. 0 IN AAAA ::1")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeAAAA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("example.refused.org. 0 IN AAAA ::ffff:0a00:0001")},
				Qname:  "example.refused.org.",
				Qtype:  dns.TypeAAAA,
			},
		},
		{
			Expected: dns.RcodeSuccess,
			test: test.Case{
				Answer: []dns.RR{test.MX("example.org. 585 IN MX 50 mx01.example.org.")},
				Qname:  "example.org.",
				Qtype:  dns.TypeMX,
			},
		},
		{
			Expected: dns.RcodeSuccess,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("example.test.valid.ipv6. 0 IN AAAA 2a04:4e42:200::644")},
				Qname:  "example.test.valid.ipv6.",
				Qtype:  dns.TypeAAAA,
			},
		},
	}

	for _, tc := range tests {

		m := new(dns.Msg)
		m.SetQuestion(tc.test.Qname, tc.test.Qtype)

		tHandler := &testHandler{
			Response: &tc.test,
			Next:     nil,
		}
		o := &Stopdnsrebind{Next: tHandler}
		w := dnstest.NewRecorder(&test.ResponseWriter{})

		_, ipNet, _ := net.ParseCIDR("192.0.2.1/24")

		if tc.config != "" {
			o.AllowList = []string{"hello.com."}
			o.DenyList = []net.IPNet{*ipNet}
		}
		_, err := o.ServeDNS(context.TODO(), w, m)

		if err != nil {
			t.Errorf("Error %q", err)
		}

		if w.Rcode != tc.Expected {
			t.Error(tc.test.Qname, "failed", "| ANSWER: ", tc.test.Answer[0], "| Rcode:", w.Rcode)
		}
	}
}
