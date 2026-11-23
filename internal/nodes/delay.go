package nodes

import "time"

type Delay struct{}

func (Delay) Type() string       { return "delay" }
func (Delay) OutPorts() []string { return []string{"out"} }

func (Delay) Run(env *Env, cfg map[string]any) error {
	ms := num(cfg, "ms")
	if ms <= 0 {
		ms = 100
	}
	t := time.NewTimer(time.Duration(ms) * time.Millisecond)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-env.Ctx.Done():
		return env.Ctx.Err()
	}
}
