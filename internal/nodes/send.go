package nodes

import (
	"fmt"
	"net/http"
)

type Send struct{}

func (Send) Type() string       { return "send" }
func (Send) OutPorts() []string { return []string{"out"} }

func (Send) Run(env *Env, cfg map[string]any) error {
	url := Interpolate(str(cfg, "url"), env.Vars)
	if url == "" {
		return fmt.Errorf("send: missing url")
	}
	method := str(cfg, "method")
	if method == "" {
		method = http.MethodGet
	}
	req, err := http.NewRequestWithContext(env.Ctx, method, url, nil)
	if err != nil {
		return err
	}
	c := env.HTTP
	if c == nil {
		c = defaultHTTP
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("send: http %d", resp.StatusCode)
	}
	return nil
}
