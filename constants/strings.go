package constants

type DefaultStrings struct {
	Loading      string
	ThanksTxt    string
	VersionInfo  string
	Title        string
	IntroName    string
	SwatName     string
	AsphodelName string
	PressEnter   string
	Targeted     string
	Splat        string
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
	Strings.AsphodelName = "asphodel"
	Strings.PressEnter = "Press Enter to Play"
	Strings.Targeted = "DANGER"
	Strings.Splat = "SWATTED"
}
