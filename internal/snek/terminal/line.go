package terminal

type Line struct {
	terminal *Terminal
	line     int
}

func (l *Line) Write(p []byte) (n int, err error) {
	l.terminal.lock.Lock()
	defer l.terminal.lock.Unlock()

	l.terminal.setRow(l.line)
	l.terminal.cursor.ClearLine()
	return l.terminal.write(p)
}

func (l *Line) Print(v ...any) (n int, err error) {
	return l.terminal.Print(l.line, v...)
}

func (l *Line) Printf(format string, v ...any) (n int, err error) {
	return l.terminal.Printf(l.line, format, v...)
}

func (l *Line) Println(v ...any) (n int, err error) {
	return l.terminal.Println(l.line, v...)
}

func (l *Line) Clear() {
	l.terminal.lock.Lock()
	defer l.terminal.lock.Unlock()

	l.terminal.setRow(l.line)
	l.terminal.cursor.ClearLine()
}
