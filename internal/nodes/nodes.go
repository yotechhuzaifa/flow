// Package nodes holds one file per node type. Adding a node type means adding a file here
// and registering it in Register — nothing else in the codebase changes.
package nodes

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Env is everything a node may touch. Features added later get what they need here,
// never by reaching into package globals.
type Env struct {
	Vars   map[string]string
	HTTP   *http.Client
	Logf   func(format string, args ...any)
	Ctx    context.Context
	Output *any // set by a node that produces a value for downstream nodes
}

// Node runs one step. Implementations must be safe to call once.
type Node interface {
	Type() string
	OutPorts() []string
	Run(env *Env, cfg map[string]any) error
}

var registry = map[string]Node{}

// Register adds node types to the registry. Called once from internal/app.
func Register(ns ...Node) {
	for _, n := range ns {
		registry[n.Type()] = n
	}
}

// Lookup returns the implementation for a node type.
func Lookup(t string) (Node, error) {
	n, ok := registry[t]
	if !ok {
		return nil, fmt.Errorf("unknown node type %q", t)
	}
	return n, nil
}

// OutPorts returns the declared output ports for a type, for graph validation.
func OutPorts(t string) []string {
	if n, ok := registry[t]; ok {
		return n.OutPorts()
	}
	return nil
}

// Interpolate replaces {{name}} from env vars.
func Interpolate(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{"+k+"}}", v)
	}
	return s
}

func str(cfg map[string]any, key string) string {
	v, _ := cfg[key].(string)
	return v
}

func num(cfg map[string]any, key string) float64 {
	switch v := cfg[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	return 0
}

var defaultHTTP = &http.Client{Timeout: 30 * time.Second}
