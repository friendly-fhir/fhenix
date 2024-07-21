package templatefuncs

import "github.com/friendly-fhir/fhenix/internal/ansi"

type ANSIModule struct{}

func (m *ANSIModule) Reset(s string) string {
	return ansi.Reset.Format(s)
}

func (m *ANSIModule) Bold(s string) string {
	return ansi.Bold.Format(s)
}

func (m *ANSIModule) Faint(s string) string {
	return ansi.Faint.Format(s)
}

func (m *ANSIModule) Italic(s string) string {
	return ansi.Italic.Format(s)
}

func (m *ANSIModule) Underline(s string) string {
	return ansi.Underline.Format(s)
}

////////////////////////////////////////////////////////////////////////////////

func (m *ANSIModule) FGBlack(s string) string {
	return ansi.FGBlack.Format(s)
}

func (m *ANSIModule) FGRed(s string) string {
	return ansi.FGRed.Format(s)
}

func (m *ANSIModule) FGGreen(s string) string {
	return ansi.FGGreen.Format(s)
}

func (m *ANSIModule) FGYellow(s string) string {
	return ansi.FGYellow.Format(s)
}

func (m *ANSIModule) FGBlue(s string) string {
	return ansi.FGBlue.Format(s)
}

func (m *ANSIModule) FGMagenta(s string) string {
	return ansi.FGMagenta.Format(s)
}

func (m *ANSIModule) FGCyan(s string) string {
	return ansi.FGCyan.Format(s)
}

func (m *ANSIModule) FGWhite(s string) string {
	return ansi.FGWhite.Format(s)
}

func (m *ANSIModule) FGGray(s string) string {
	return ansi.FGGray.Format(s)
}

func (m *ANSIModule) FGBrightRed(s string) string {
	return ansi.FGBrightRed.Format(s)
}

func (m *ANSIModule) FGBrightGreen(s string) string {
	return ansi.FGBrightGreen.Format(s)
}

func (m *ANSIModule) FGBrightYellow(s string) string {
	return ansi.FGBrightYellow.Format(s)
}

func (m *ANSIModule) FGBrightBlue(s string) string {
	return ansi.FGBrightBlue.Format(s)
}

func (m *ANSIModule) FGBrightMagenta(s string) string {
	return ansi.FGBrightMagenta.Format(s)
}

func (m *ANSIModule) FGBrightCyan(s string) string {
	return ansi.FGBrightCyan.Format(s)
}

func (m *ANSIModule) FGBrightWhite(s string) string {
	return ansi.FGBrightWhite.Format(s)
}

////////////////////////////////////////////////////////////////////////////////

func (m *ANSIModule) BGBlack(s string) string {
	return ansi.BGBlack.Format(s)
}

func (m *ANSIModule) BGRed(s string) string {
	return ansi.BGRed.Format(s)
}

func (m *ANSIModule) BGGreen(s string) string {
	return ansi.BGGreen.Format(s)
}

func (m *ANSIModule) BGYellow(s string) string {
	return ansi.BGYellow.Format(s)
}

func (m *ANSIModule) BGBlue(s string) string {
	return ansi.BGBlue.Format(s)
}

func (m *ANSIModule) BGMagenta(s string) string {
	return ansi.BGMagenta.Format(s)
}

func (m *ANSIModule) BGCyan(s string) string {
	return ansi.BGCyan.Format(s)
}

func (m *ANSIModule) BGWhite(s string) string {
	return ansi.BGWhite.Format(s)
}

func (m *ANSIModule) BGGray(s string) string {
	return ansi.BGGray.Format(s)
}

func (m *ANSIModule) BGBrightRed(s string) string {
	return ansi.BGBrightRed.Format(s)
}

func (m *ANSIModule) BGBrightGreen(s string) string {
	return ansi.BGBrightGreen.Format(s)
}

func (m *ANSIModule) BGBrightYellow(s string) string {
	return ansi.BGBrightYellow.Format(s)
}

func (m *ANSIModule) BGBrightBlue(s string) string {
	return ansi.BGBrightBlue.Format(s)
}

func (m *ANSIModule) BGBrightMagenta(s string) string {
	return ansi.BGBrightMagenta.Format(s)
}

func (m *ANSIModule) BGBrightCyan(s string) string {
	return ansi.BGBrightCyan.Format(s)
}

func (m *ANSIModule) BGBrightWhite(s string) string {
	return ansi.BGBrightWhite.Format(s)
}
