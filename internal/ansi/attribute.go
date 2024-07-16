package ansi

// Attribute is a unit type which represents individual ANSI formatting attributes.
type Attribute interface {
	Display
	isAttribute()
}

const (
	// Reset is an ANSI attribute for resetting formatting
	Reset attribute = 0

	// Bold is an ANSI attribute for setting a bold format
	Bold attribute = 1

	// Faint is an ANSI attribute for setting a faint format
	Faint attribute = 2

	// Italice is an ANSI attribute for setting an italic format
	Italic attribute = 3

	// Underline is an ANSI attribute for setting an underline format
	Underline attribute = 4

	// DefaultFont is an ANSI attribute for setting the default font
	DefaultFont attribute = 10

	// FGBlack is an ANSI attribute for the foreground color black
	FGBlack attribute = 30

	// FGRed is an ANSI attribute for the foreground color red
	FGRed attribute = 31

	// FGGreen is an ANSI attribute for the foreground color green
	FGGreen attribute = 32

	// FGYellow is an ANSI attribute for the foreground color yellow
	FGYellow attribute = 33

	// FGBlue is an ANSI attribute for the foreground color blue
	FGBlue attribute = 34

	// FGMagenta is an ANSI attribute for the foreground color magenta
	FGMagenta attribute = 35

	// FGCyan is an ANSI attribute for the foreground color cyan
	FGCyan attribute = 36

	// FGWhite is an ANSI attribute for the foreground color white
	FGWhite attribute = 37

	// FGGray is an ANSI attribute for the foreground color gray
	FGGray attribute = 90

	// FGBrightRed is an ANSI attribute for the foreground color brightred
	FGBrightRed attribute = 91

	// FGBrightGreen is an ANSI attribute for the foreground color brightgreen
	FGBrightGreen attribute = 92

	// FGBrightYellow is an ANSI attribute for the foreground color brightyellow
	FGBrightYellow attribute = 93

	// FGBrightBlue is an ANSI attribute for the foreground color brightblue
	FGBrightBlue attribute = 94

	// FGBrightMagenta is an ANSI attribute for the foreground color brightmagenta
	FGBrightMagenta attribute = 95

	// FGBrightCyan is an ANSI attribute for the foreground color brightcyan
	FGBrightCyan attribute = 96

	// FGBrightWhite is an ANSI attribute for the foreground color brightwhite
	FGBrightWhite attribute = 97

	// FGDefault is an ANSI attribute for setting the default foreground color
	FGDefault attribute = 39

	// BGBlack is an ANSI attribute for the background color black
	BGBlack attribute = 40

	// BGRed is an ANSI attribute for the background color red
	BGRed attribute = 41

	// BGGreen is an ANSI attribute for the background color green
	BGGreen attribute = 42

	// BGYellow is an ANSI attribute for the background color yellow
	BGYellow attribute = 43

	// BGBlue is an ANSI attribute for the background color blue
	BGBlue attribute = 44

	// BGMagenta is an ANSI attribute for the background color magenta
	BGMagenta attribute = 45

	// BGCyan is an ANSI attribute for the background color cyan
	BGCyan attribute = 46

	// BGWhite is an ANSI attribute for the background color white
	BGWhite attribute = 47

	// BGGray is an ANSI attribute for the background color gray
	BGGray attribute = 100

	// BGBrightRed is an ANSI attribute for the background color brightred
	BGBrightRed attribute = 101

	// BGBrightGreen is an ANSI attribute for the background color brightgreen
	BGBrightGreen attribute = 102

	// BGBrightYellow is an ANSI attribute for the background color brightyellow
	BGBrightYellow attribute = 103

	// BGBrightBlue is an ANSI attribute for the background color brightblue
	BGBrightBlue attribute = 104

	// BGBrightMagenta is an ANSI attribute for the background color brightmagenta
	BGBrightMagenta attribute = 105

	// BGBrightCyan is an ANSI attribute for the background color brightcyan
	BGBrightCyan attribute = 106

	// BGBrightWhite is an ANSI attribute for the background color brightwhite
	BGBrightWhite attribute = 107

	// BGDefault is an ANSI attribute for setting the default background color
	BGDefault attribute = 49
)

// attribute represents an ANSI attribute formatting
type attribute uint8

func (a attribute) Format(format string, args ...any) string {
	return Format(a).Format(format, args...)
}

func (a attribute) String() string {
	return createFormatFunc(byte(a))
}

func (a attribute) FormatString() string {
	return ansiFormat(byte(a))
}

func (a attribute) codes() []byte {
	if !enabled {
		return nil
	}
	return []byte{byte(a)}
}

func (a attribute) len() int {
	return 1
}

func (c attribute) isAttribute() {

}

var _ Display = (*attribute)(nil)
var _ Attribute = (*attribute)(nil)
