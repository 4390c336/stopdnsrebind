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
