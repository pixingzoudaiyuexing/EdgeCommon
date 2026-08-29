// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package nodeconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

const (
	// DefaultHTTPCCPolicyMaxConnectionsPerIP 是 1.3.9 Plus 的单 IP 默认最大并发连接数。
	DefaultHTTPCCPolicyMaxConnectionsPerIP = 30

	// 下列参数用于识别异常连续重定向行为。
	DefaultHTTPCCPolicyRedirectsCheckingValidatePath     = "/GE/CC/VALIDATOR"
	DefaultHTTPCCPolicyRedirectsCheckingDurationSeconds = 120
	DefaultHTTPCCPolicyRedirectsCheckingMaxRedirects    = 30
	DefaultHTTPCCPolicyRedirectsCheckingBlockSeconds    = 3600
)

// HTTPCCRedirectsCheckingConfig 连续重定向检测配置。
type HTTPCCRedirectsCheckingConfig struct {
	ValidatePath    string `json:"validatePath" yaml:"validatePath"`
	DurationSeconds int    `json:"durationSeconds" yaml:"durationSeconds"`
	MaxRedirects    int    `json:"maxRedirects" yaml:"maxRedirects"`
	BlockSeconds    int    `json:"blockSeconds" yaml:"blockSeconds"`
}

// HTTPCCFirewallConfig 定义 CC 封禁的作用范围。
type HTTPCCFirewallConfig struct {
	Scope firewallconfigs.FirewallScope `json:"scope" yaml:"scope"`
}

// HTTPCCPolicy 集群 CC 策略。
type HTTPCCPolicy struct {
	IsOn                bool                             `json:"isOn" yaml:"isOn"`
	Thresholds          []*serverconfigs.HTTPCCThreshold `json:"thresholds" yaml:"thresholds"` // 阈值
	MaxConnectionsPerIP int                              `json:"maxConnectionsPerIP" yaml:"maxConnectionsPerIP"`
	RedirectsChecking   HTTPCCRedirectsCheckingConfig    `json:"redirectsChecking" yaml:"redirectsChecking"`
	Firewall            HTTPCCFirewallConfig             `json:"firewall" yaml:"firewall"`
}

func NewHTTPCCPolicy() *HTTPCCPolicy {
	return &HTTPCCPolicy{
		IsOn: true,
		Firewall: HTTPCCFirewallConfig{
			Scope: firewallconfigs.FirewallScopeGlobal,
		},
	}
}

func (this *HTTPCCPolicy) Init() error {
	return nil
}

// FirewallScope 返回 CC 封禁范围；旧配置未写 scope 时按 1.3.9 行为回退为 global。
func (this *HTTPCCPolicy) FirewallScope() firewallconfigs.FirewallScope {
	if len(this.Firewall.Scope) == 0 {
		return firewallconfigs.FirewallScopeGlobal
	}
	return this.Firewall.Scope
}
