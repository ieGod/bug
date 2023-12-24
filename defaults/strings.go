package defaults

type DefaultStrings struct {
	Loading     string
	ThanksTxt   string
	VersionInfo string
	Title       string
	IntroName   string
	SwatName    string
	PressEnter  string
}

var (
	Strings DefaultStrings
)

func init() {
	Strings.VersionInfo = "bsoft games bug v0.01"
	Strings.Title = "bug"
	Strings.Loading = "Loading"
	Strings.ThanksTxt = "buh bye thanks for playing"
	Strings.IntroName = "intro"
	Strings.SwatName = "swat"
	Strings.PressEnter = "Press Enter to Play"
}
