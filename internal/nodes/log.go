package nodes

type Log struct{}

func (Log) Type() string       { return "log" }
func (Log) OutPorts() []string { return []string{"out"} }

func (Log) Run(env *Env, cfg map[string]any) error {
	msg := Interpolate(str(cfg, "message"), env.Vars)
	if env.Logf != nil {
		env.Logf("log node: %s", msg)
	}
	return nil
}

type If struct{}

func (If) Type() string       { return "if" }
func (If) OutPorts() []string { return []string{"true", "false"} }

// Run is a stub: it always takes the true branch until the branch-expression row lands.
func (If) Run(env *Env, cfg map[string]any) error { return nil }
