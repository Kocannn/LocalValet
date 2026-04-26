package terminal

type LaunchOptions struct {
	ProjectDir        string
	PreferredTerminal string
}

type LaunchResult struct {
	Terminal string
	WorkDir  string
}

type Manager interface {
	Launch(options LaunchOptions) (LaunchResult, error)
}
