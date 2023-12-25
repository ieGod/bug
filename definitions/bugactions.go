package definitions

type BugAction int

const (
	BugActionIdle BugAction = iota
	BugActionForwardRun
	BugActionSideRun
	BugActionReverseRun
	BugActionGlitch
)
