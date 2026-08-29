package nodeconfigs

import (
	"testing"

	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

func TestUAMPolicyFirewallScope(t *testing.T) {
	policy := NewUAMPolicy()
	if scope := policy.FirewallScope(); scope != firewallconfigs.FirewallScopeGlobal {
		t.Fatalf("默认作用域应为 global，实际为 %q", scope)
	}

	policy.Firewall.Scope = firewallconfigs.FirewallScopeServer
	if scope := policy.FirewallScope(); scope != firewallconfigs.FirewallScopeServer {
		t.Fatalf("显式 server 作用域应保持不变，实际为 %q", scope)
	}

	policy.Firewall.Scope = ""
	if scope := policy.FirewallScope(); scope != firewallconfigs.FirewallScopeGlobal {
		t.Fatalf("空作用域应回退为 global，实际为 %q", scope)
	}

	var nilPolicy *UAMPolicy
	if scope := nilPolicy.FirewallScope(); scope != firewallconfigs.FirewallScopeGlobal {
		t.Fatalf("nil 策略应安全回退为 global，实际为 %q", scope)
	}
}
