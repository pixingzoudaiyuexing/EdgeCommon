// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package nodeconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"

func init() {
	_ = DefaultUAMPolicy.Init()
}

// DefaultUAMPolicy 保持 1.3.9 Plus 的集群 UAM 默认行为。
var DefaultUAMPolicy = &UAMPolicy{
	IsOn:               true,
	AllowSearchEngines: true,
	DenySpiders:        true,
	MaxFails:           30,
	BlockSeconds:       1800,
	IncludeSubdomains:  true,
	UITitle:            "",
	UIBody:             "",
	KeyLife:            3600,
	Firewall: UAMFirewallConfig{
		Scope: firewallconfigs.FirewallScopeGlobal,
	},
}

// UAMFirewallConfig 定义 UAM 校验结果的防火墙作用范围。
type UAMFirewallConfig struct {
	Scope firewallconfigs.FirewallScope `yaml:"scope" json:"scope"`
}

// UAMPolicy 集群 UAM（5 秒盾）策略。
type UAMPolicy struct {
	IsOn               bool `yaml:"isOn" json:"isOn"`
	AllowSearchEngines bool `yaml:"allowSearchEngines" json:"allowSearchEngines"` // 直接跳过常见搜索引擎
	DenySpiders        bool `yaml:"denySpiders" json:"denySpiders"`               // 拦截常见爬虫

	MaxFails          int  `yaml:"maxFails" json:"maxFails"`                   // 最大校验失败次数
	BlockSeconds      int  `yaml:"blockSeconds" json:"blockSeconds"`           // 校验失败后的封禁时长
	IncludeSubdomains bool `yaml:"includeSubdomains" json:"includeSubdomains"` // 是否允许子域共享校验结果

	UITitle string `yaml:"uiTitle" json:"uiTitle"` // 页面标题
	UIBody  string `yaml:"uiBody" json:"uiBody"`   // 页面内容
	KeyLife int    `yaml:"keyLife" json:"keyLife"` // 校验 Key 有效期

	Firewall UAMFirewallConfig `yaml:"firewall" json:"firewall"`
}

func NewUAMPolicy() *UAMPolicy {
	policy := *DefaultUAMPolicy
	return &policy
}

func (this *UAMPolicy) Init() error {
	return nil
}

// FirewallScope 返回 UAM 连续校验失败后的封禁范围。
// 兼容旧策略：没有写入 scope 时按 1.3.9 默认行为回退为全局范围。
func (this *UAMPolicy) FirewallScope() firewallconfigs.FirewallScope {
	if this == nil || len(this.Firewall.Scope) == 0 {
		return firewallconfigs.FirewallScopeGlobal
	}
	return this.Firewall.Scope
}
