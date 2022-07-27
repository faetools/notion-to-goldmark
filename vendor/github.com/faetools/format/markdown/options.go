package markdown

import (
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A Config struct has configurations for the markdown based renderers.
type Config struct {
	HardWraps bool

	// Terminal specifies if ANSI escape codes are emitted for styling.
	Terminal bool
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config { return Config{} }

// SetOption implements renderer.NodeRenderer.SetOption.
func (c *Config) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optHardWraps:
		c.HardWraps = value.(bool)
	case optTerminal:
		c.Terminal = value.(bool)
	}
}

func (r *NodeRenderer) getOptions() []renderer.Option {
	opts := []renderer.Option{
		renderer.WithNodeRenderers(util.Prioritized(r, 0)),
	}

	for _, setting := range []struct {
		enabled bool
		opt     renderer.Option
	}{
		{r.HardWraps, WithHardWraps()},
		{r.Terminal, WithTerminal()},
	} {
		if setting.enabled {
			opts = append(opts, setting.opt)
		}
	}

	opts = append(opts, r.additionalOptions...)

	return opts
}

// An Option interface sets options for markdown based renderers.
type Option interface {
	SetMarkdownOption(*Config)
	renderer.Option
}

// HardWraps is an option name used in WithHardWraps.
const optHardWraps renderer.OptionName = "HardWraps"

type withHardWraps struct{}

func (o *withHardWraps) SetConfig(c *renderer.Config) {
	c.Options[optHardWraps] = true
}

func (o *withHardWraps) SetMarkdownOption(c *Config) {
	c.HardWraps = true
}

// WithHardWraps is a functional option that indicates whether softline breaks
// should be rendered as '<br>'.
func WithHardWraps() Option {
	return &withHardWraps{}
}

// Terminal is an option name used in WithTerminal.
const optTerminal renderer.OptionName = "Terminal"

type withTerminal struct{}

func (o *withTerminal) SetConfig(c *renderer.Config) {
	c.Options[optTerminal] = true
}

func (o *withTerminal) SetMarkdownOption(c *Config) {
	c.Terminal = true
}

// WithTerminal is a functional option that indicates the result is written to the terminal.
func WithTerminal() Option {
	return &withTerminal{}
}
