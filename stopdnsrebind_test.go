package stopdnsrebind

import (
	"context"
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
				Answer: []dns.RR{test.A("example.org. 0 IN A 169.254.169.254")},
				Qname:  "example.org.",
				Qtype:  dns.TypeA,
			},
		},
		{
			Expected: dns.RcodeSuccess,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("example.org. 0 IN AAAA 2a00:1450:4009:823::200e")},
				Qname:  "example.org.",
				Qtype:  dns.TypeAAAA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("example.org. 0 IN AAAA ::1")},
				Qname:  "example.org.",
				Qtype:  dns.TypeAAAA,
			},
		},
		{
			Expected: dns.RcodeRefused,
			test: test.Case{
				Answer: []dns.RR{test.AAAA("example.org. 0 IN AAAA ::ffff:0a00:0001")},
				Qname:  "example.org.",
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

		if tc.config != "" {
			o.AllowList = []string{"hello.com."}
		}
		_, err := o.ServeDNS(context.TODO(), w, m)

		if err != nil {
			t.Errorf("Error %q", err)
		}

		if w.Rcode != tc.Expected {
			t.Error("Not the expected response", tc.test.Answer[0], "Rcode:", w.Rcode)
		}
	}
}
