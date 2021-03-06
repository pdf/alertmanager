// Copyright 2015 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"strings"
)

var (
	// DefaultWebhookConfig defines default values for Webhook configurations.
	DefaultWebhookConfig = WebhookConfig{
		NotifierConfig: NotifierConfig{
			VSendResolved: true,
		},
	}

	// DefaultEmailConfig defines default values for Email configurations.
	DefaultEmailConfig = EmailConfig{
		NotifierConfig: NotifierConfig{
			VSendResolved: false,
		},
		HTML: `{{ template "email.default.html" . }}`,
	}

	// DefaultEmailSubject defines the default Subject header of an Email.
	DefaultEmailSubject = `{{ template "email.default.subject" . }}`

	// DefaultPagerdutyConfig defines default values for PagerDuty configurations.
	DefaultPagerdutyConfig = PagerdutyConfig{
		NotifierConfig: NotifierConfig{
			VSendResolved: true,
		},
		Description: `{{ template "pagerduty.default.description" .}}`,
		Client:      `{{ template "pagerduty.default.client" . }}`,
		ClientURL:   `{{ template "pagerduty.default.clientURL" . }}`,
		Details: map[string]string{
			"firing":       `{{ template "pagerduty.default.instances" .Alerts.Firing }}`,
			"resolved":     `{{ template "pagerduty.default.instances" .Alerts.Resolved }}`,
			"num_firing":   `{{ .Alerts.Firing | len }}`,
			"num_resolved": `{{ .Alerts.Resolved | len }}`,
		},
	}

	// DefaultSlackConfig defines default values for Slack configurations.
	DefaultSlackConfig = SlackConfig{
		NotifierConfig: NotifierConfig{
			VSendResolved: true,
		},
		Color:     `{{ if eq .Status "firing" }}danger{{ else }}good{{ end }}`,
		Username:  `{{ template "slack.default.username" . }}`,
		Title:     `{{ template "slack.default.title" . }}`,
		TitleLink: `{{ template "slack.default.titlelink" . }}`,
		Pretext:   `{{ template "slack.default.pretext" . }}`,
		Text:      `{{ template "slack.default.text" . }}`,
		Fallback:  `{{ template "slack.default.fallback" . }}`,
	}

	// DefaultOpsGenieConfig defines default values for OpsGenie configurations.
	DefaultOpsGenieConfig = OpsGenieConfig{
		NotifierConfig: NotifierConfig{
			VSendResolved: true,
		},
		Description: `{{ template "opsgenie.default.description" . }}`,
		Source:      `{{ template "opsgenie.default.source" . }}`,
		// TODO: Add a details field with all the alerts.
	}
)

// NotifierConfig contains base options common across all notifier configurations.
type NotifierConfig struct {
	VSendResolved bool `yaml:"send_resolved"`
}

func (nc *NotifierConfig) SendResolved() bool {
	return nc.VSendResolved
}

// EmailConfig configures notifications via mail.
type EmailConfig struct {
	NotifierConfig

	// Email address to notify.
	To        string            `yaml:"to"`
	From      string            `yaml:"from"`
	Smarthost string            `yaml:"smarthost,omitempty"`
	Headers   map[string]string `yaml:"headers"`
	HTML      string            `yaml:"html"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *EmailConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultEmailConfig
	type plain EmailConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.To == "" {
		return fmt.Errorf("missing to address in email config")
	}
	// Header names are case-insensitive, check for collisions.
	normalizedHeaders := map[string]string{}
	for h, v := range c.Headers {
		normalized := strings.ToTitle(h)
		if _, ok := normalizedHeaders[normalized]; ok {
			return fmt.Errorf("duplicate header %q in email config", normalized)
		}
		normalizedHeaders[normalized] = v
	}
	c.Headers = normalizedHeaders

	return checkOverflow(c.XXX, "email config")
}

// PagerdutyConfig configures notifications via PagerDuty.
type PagerdutyConfig struct {
	NotifierConfig

	ServiceKey  Secret            `yaml:"service_key"`
	URL         string            `yaml:"url"`
	Client      string            `yaml:"client"`
	ClientURL   string            `yaml:"client_url"`
	Description string            `yaml:"description"`
	Details     map[string]string `yaml:"details"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *PagerdutyConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultPagerdutyConfig
	type plain PagerdutyConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.ServiceKey == "" {
		return fmt.Errorf("missing service key in PagerDuty config")
	}
	return checkOverflow(c.XXX, "pagerduty config")
}

// SlackConfig configures notifications via Slack.
type SlackConfig struct {
	NotifierConfig

	APIURL Secret `yaml:"api_url"`

	// Slack channel override, (like #other-channel or @username).
	Channel  string `yaml:"channel"`
	Username string `yaml:"username"`
	Color    string `yaml:"color"`

	Title     string `yaml:"title"`
	TitleLink string `yaml:"title_link"`
	Pretext   string `yaml:"pretext"`
	Text      string `yaml:"text"`
	Fallback  string `yaml:"fallback"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *SlackConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultSlackConfig
	type plain SlackConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Channel == "" {
		return fmt.Errorf("missing channel in Slack config")
	}
	return checkOverflow(c.XXX, "slack config")
}

// WebhookConfig configures notifications via a generic webhook.
type WebhookConfig struct {
	NotifierConfig

	// URL to send POST request to.
	URL string `yaml:"url"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *WebhookConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultWebhookConfig
	type plain WebhookConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.URL == "" {
		return fmt.Errorf("missing URL in webhook config")
	}
	return checkOverflow(c.XXX, "slack config")
}

// OpsGenieConfig configures notifications via OpsGenie.
type OpsGenieConfig struct {
	NotifierConfig

	APIKey      Secret            `yaml:"api_key"`
	APIHost     string            `yaml:"api_host"`
	Description string            `yaml:"description"`
	Source      string            `yaml:"source"`
	Details     map[string]string `yaml:"details"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *OpsGenieConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultOpsGenieConfig
	type plain OpsGenieConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.APIKey == "" {
		return fmt.Errorf("missing API key in OpsGenie config")
	}
	return checkOverflow(c.XXX, "opsgenie config")
}
