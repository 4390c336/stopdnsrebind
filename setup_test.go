package stopdnsrebind

import (
	"testing"

	"github.com/coredns/caddy"
)

func Test_setup(t *testing.T) {

	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			"allow internal.example.org",
			`stopdnsrebind {
				allow internal.example.org.
			}`,
			false,
		},
		{
			"allow multiple",
			`stopdnsrebind {
				allow internal.example.org. internal.example.net.
			}`,
			false,
		},
		{
			"non supported op",
			`stopdnsrebind {
				anything internal.example.org.
			}`,
			true,
		},
		{
			"not a valid domain",
			`stopdnsrebind {
				allow ..example.org.
			}`,
			true,
		},
		{
			"deny a valid ipNet",
			`stopdnsrebind {
				deny 192.0.2.1/24
			}`,
			false,
		},
		{
			"deny multiple ipNet",
			`stopdnsrebind {
				deny 192.0.2.1/24 127.0.0.1/8
			}`,
			false,
		},
		{
			"deny invalid ipNet",
			`stopdnsrebind {
				deny 192.0.2.1
			}`,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := caddy.NewTestController("dns", tt.config)
			if err := setup(ctr); (err != nil) != tt.wantErr {
				t.Errorf("setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
