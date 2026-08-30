// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

// HTTPCCThreshold CC 请求阈值。
//
// 1.3.9 Plus 使用三个 int32 字段保存周期、最大请求数和封禁时长；
// 保持相同字段类型可以避免既有 JSON/配置在自维护版本中出现语义偏差。
type HTTPCCThreshold struct {
	PeriodSeconds int32 `yaml:"periodSeconds" json:"periodSeconds"` // 统计周期，单位秒
	MaxRequests   int32 `yaml:"maxRequests" json:"maxRequests"`     // 周期内允许的最大请求数
	BlockSeconds  int32 `yaml:"blockSeconds" json:"blockSeconds"`   // 触发后封禁时长，单位秒
}

func NewHTTPCCThreshold() *HTTPCCThreshold {
	return &HTTPCCThreshold{}
}

// Merge 用 threshold 中大于 0 的值覆盖当前阈值。
func (this *HTTPCCThreshold) Merge(threshold *HTTPCCThreshold) {
	if threshold == nil {
		return
	}
	if threshold.PeriodSeconds > 0 {
		this.PeriodSeconds = threshold.PeriodSeconds
	}
	if threshold.MaxRequests > 0 {
		this.MaxRequests = threshold.MaxRequests
	}
	if threshold.BlockSeconds > 0 {
		this.BlockSeconds = threshold.BlockSeconds
	}
}

// MergeIfEmpty 只在当前值未设置时，从 threshold 补入大于 0 的值。
func (this *HTTPCCThreshold) MergeIfEmpty(threshold *HTTPCCThreshold) {
	if threshold == nil {
		return
	}
	if this.PeriodSeconds <= 0 && threshold.PeriodSeconds > 0 {
		this.PeriodSeconds = threshold.PeriodSeconds
	}
	if this.MaxRequests <= 0 && threshold.MaxRequests > 0 {
		this.MaxRequests = threshold.MaxRequests
	}
	if this.BlockSeconds <= 0 && threshold.BlockSeconds > 0 {
		this.BlockSeconds = threshold.BlockSeconds
	}
}

func (this *HTTPCCThreshold) Clone() *HTTPCCThreshold {
	if this == nil {
		return nil
	}
	return &HTTPCCThreshold{
		PeriodSeconds: this.PeriodSeconds,
		MaxRequests:   this.MaxRequests,
		BlockSeconds:  this.BlockSeconds,
	}
}

// DefaultHTTPCCThresholds 是 1.3.9 Plus 的默认三档 CC 阈值。
var DefaultHTTPCCThresholds = []*HTTPCCThreshold{
	{PeriodSeconds: 5, MaxRequests: 60, BlockSeconds: 600},
	{PeriodSeconds: 60, MaxRequests: 150, BlockSeconds: 2400},
	{PeriodSeconds: 300, MaxRequests: 300, BlockSeconds: 3600},
}

// CloneDefaultHTTPCCThresholds 返回默认阈值的独立副本，避免调用方修改全局默认值。
func CloneDefaultHTTPCCThresholds() []*HTTPCCThreshold {
	thresholds := make([]*HTTPCCThreshold, 0, len(DefaultHTTPCCThresholds))
	for _, threshold := range DefaultHTTPCCThresholds {
		thresholds = append(thresholds, threshold.Clone())
	}
	return thresholds
}

// DefaultHTTPCCConfig 默认的 CC 配置。
func DefaultHTTPCCConfig() *HTTPCCConfig {
	config := NewHTTPCCConfig()
	config.IsOn = true
	config.Thresholds = CloneDefaultHTTPCCThresholds()
	return config
}

// HTTPCCConfig HTTP CC 防护配置。
//
// 字段顺序与可信 1.3.9 Plus edge-node 的 Go 运行时类型元数据保持一致。
// Level、WithRequestPath、IgnoreCommonAgents 和 Action 在公开 !plus 版本中曾被裁掉，
// 因此仍需保留这些字段以兼容既有 JSON 和历史配置。
//
// 可信 1.3.9 Plus edge-node 静态审计进一步确认：WithRequestPath 会被 doCC() 读取；
// Level、IgnoreCommonAgents、Action 则不会被 doCC() 直接读取，HTTPCCConfig.Init() / MatchURL()
// 也不会把它们转换成其他运行时配置。自维护版本不得仅凭字段名为这三个兼容字段新增
// 请求级跳过、放行或阻断语义。
type HTTPCCConfig struct {
	IsPrior bool   `yaml:"isPrior" json:"isPrior"` // 是否覆盖父级
	IsOn    bool   `yaml:"isOn" json:"isOn"`       // 是否启用
	Level   string `yaml:"level" json:"level"`     // 历史兼容字段；1.3.9 Node 请求链不直接读取

	WithRequestPath      bool               `yaml:"withRequestPath" json:"withRequestPath"`           // 请求计数是否区分路径
	UseDefaultThresholds bool               `yaml:"useDefaultThresholds" json:"useDefaultThresholds"` // 是否使用默认阈值
	Thresholds           []*HTTPCCThreshold `yaml:"thresholds" json:"thresholds"`                     // 自定义阈值
	IgnoreCommonFiles    bool               `yaml:"ignoreCommonFiles" json:"ignoreCommonFiles"`       // 是否忽略常见静态文件
	IgnoreCommonAgents   bool               `yaml:"ignoreCommonAgents" json:"ignoreCommonAgents"`     // 历史兼容字段；1.3.9 Node 请求链不直接读取
	Action               string             `yaml:"action" json:"action"`                             // 历史兼容字段；1.3.9 Node 请求链不直接读取

	OnlyURLPatterns   []*shared.URLPattern `yaml:"onlyURLPatterns" json:"onlyURLPatterns"`     // 仅限的 URL
	ExceptURLPatterns []*shared.URLPattern `yaml:"exceptURLPatterns" json:"exceptURLPatterns"` // 排除的 URL

	EnableFingerprint bool `yaml:"enableFingerprint" json:"enableFingerprint"` // 是否启用浏览器指纹校验
	EnableGET302      bool `yaml:"enableGET302" json:"enableGET302"`           // 是否启用 GET 302 校验
	MinQPSPerIP       int  `yaml:"minQPSPerIP" json:"minQPSPerIP"`             // 启用要求的单 IP 最低平均 QPS
}

func NewHTTPCCConfig() *HTTPCCConfig {
	return &HTTPCCConfig{
		EnableFingerprint:    true,
		EnableGET302:         true,
		UseDefaultThresholds: true,
		IgnoreCommonFiles:    true,
	}
}

func (this *HTTPCCConfig) Init() error {
	for _, pattern := range this.OnlyURLPatterns {
		if pattern == nil {
			continue
		}
		if err := pattern.Init(); err != nil {
			return err
		}
	}

	for _, pattern := range this.ExceptURLPatterns {
		if pattern == nil {
			continue
		}
		if err := pattern.Init(); err != nil {
			return err
		}
	}

	return nil
}

func (this *HTTPCCConfig) MatchURL(url string) bool {
	if len(this.ExceptURLPatterns) > 0 {
		for _, pattern := range this.ExceptURLPatterns {
			if pattern != nil && pattern.Match(url) {
				return false
			}
		}
	}

	if len(this.OnlyURLPatterns) > 0 {
		for _, pattern := range this.OnlyURLPatterns {
			if pattern != nil && pattern.Match(url) {
				return true
			}
		}
		return false
	}

	return true
}
