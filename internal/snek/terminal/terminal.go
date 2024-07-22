package terminal

import (
	"fmt"
	"io"
	"sync"

	"atomicgo.dev/cursor"
	"golang.org/x/term"
)

type Terminal struct {
	// lock guards the state of this terminal, and ensures that concurrent writes
	// are safe.
	lock sync.Mutex

	// col is the current column that the cursor has been moved to.
	col int

	// row is the current row that the cursor has been moved to.
	row int

	// maxRow is the maximum row that the cursor has been moved to.
	maxRow int

	// w is the writer that the cursor writes to.
	w cursor.Writer

	// cursor is the cursor that is used to move the cursor around the terminal.
	cursor *cursor.Cursor
}

// New constructs a terminal from the given writer, if it's a valid terminal.
//
// If the writer is not a terminal, then terminal will be nil, and ok will be
// false.
func New(w io.Writer) (terminal *Terminal, ok bool) {
	if fd, ok := w.(cursor.Writer); ok && term.IsTerminal(int(fd.Fd())) {
		terminal := &Terminal{
			w:      fd,
			cursor: cursor.NewCursor().WithWriter(fd),
		}
		return terminal, true
	}
	return nil, false
}

// Close closes the terminal, and re-shows the cursor if it was hidden.
func (t *Terminal) Close() {
	t.cursor.Show()
}

// HideCursor hides the cursor.
func (t *Terminal) HideCursor() {
	t.cursor.Hide()
}

// ShowCursor shows the cursor.
func (t *Terminal) ShowCursor() {
	t.cursor.Show()
}

// Println prints the given arguments to the terminal, followed by a newline.
func (t *Terminal) Println(n int, args ...any) (int, error) {
	formatted := fmt.Sprintln(args...)

	t.lock.Lock()
	defer t.lock.Unlock()
	t.setRow(n)
	t.setColumn(0)
	t.cursor.ClearLine()
	return t.write([]byte(formatted))
}

// Print prints the given arguments to the terminal.
func (t *Terminal) Print(n int, args ...any) (int, error) {
	formatted := fmt.Sprint(args...)

	t.lock.Lock()
	defer t.lock.Unlock()
	t.setRow(n)
	t.setColumn(0)
	t.cursor.ClearLine()
	return t.write([]byte(formatted))
}

// Printf prints the given arguments to the terminal, using the given format.
func (t *Terminal) Printf(n int, format string, args ...any) (int, error) {
	formatted := fmt.Sprintf(format, args...)

	t.lock.Lock()
	defer t.lock.Unlock()
	t.setRow(n)
	t.setColumn(0)
	t.cursor.ClearLine()
	return t.write([]byte(formatted))
}

// Offset returns the current row and column of the cursor, offset from where
// the cursor was initially created.
func (t *Terminal) Offset() (int, int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.row, t.col
}

// Row returns the current row of the cursor.
func (t *Terminal) Row() int {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.row
}

// Column returns the current column of the cursor.
func (t *Terminal) Column() int {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.col
}

func (t *Terminal) Width() int {
	w, _, err := term.GetSize(int(t.w.Fd()))
	if err != nil {
		return 80
	}
	return w
}

// Write writes the given bytes to the cursor, updating the row and column.
func (t *Terminal) Write(p []byte) (int, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.write(p)
}

// Up moves the cursor up the given number of rows. It cannot exceed the point
// where the cursor was initially created.
func (t *Terminal) Up(rows int) *Terminal {
	if rows < 0 {
		return t.Down(-rows)
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	offset := min(max(rows, 0), t.row)
	t.row -= offset
	t.cursor.Up(offset)
	return t
}

func (t *Terminal) Down(rows int) *Terminal {
	if rows < 0 {
		return t.Up(-rows)
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.row += rows
	t.cursor.Down(rows)
	t.maxRow = max(t.maxRow, t.row)
	return t
}

func (t *Terminal) Left(cols int) *Terminal {
	if cols < 0 {
		return t.Right(-cols)
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	offset := min(max(cols, 0), t.col)
	t.col -= offset
	t.cursor.Left(offset)
	return t
}

func (t *Terminal) Right(cols int) *Terminal {
	if cols < 0 {
		return t.Left(-cols)
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.col += cols
	t.cursor.Right(cols)
	return t
}

func (t *Terminal) Top() *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.setRow(0)
	return t
}

func (t *Terminal) Bottom() *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.translateRow(t.maxRow - t.row)
}

func (t *Terminal) Move(row, col int) *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.translateColumn(row).translateRow(col)
}

func (t *Terminal) MoveColumn(col int) *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.translateColumn(col)
}

func (t *Terminal) MoveRow(row int) *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.translateRow(row)
}

func (t *Terminal) SetPosition(row, col int) *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.setRow(row).setColumn(col)
}

func (t *Terminal) SetRow(n int) *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.setRow(n)
}

func (t *Terminal) SetColumn(n int) *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.setColumn(n)
}

func (t *Terminal) Clear() *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.cursor.Clear()
	t.col--
	return t
}

func (t *Terminal) ClearLine() *Terminal {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.cursor.ClearLine()
	t.col = 0
	return t
}

func (t *Terminal) Line(n int) *Line {
	return &Line{terminal: t, line: n}
}

func (t *Terminal) write(p []byte) (int, error) {
	for _, b := range p {
		if b == '\n' {
			t.row++
			t.col = 0
			continue
		}
		if b == '\r' {
			t.col = 0
			continue
		}
		if b == '\b' {
			t.col--
			continue
		}
		t.col++
	}
	t.maxRow = max(t.maxRow, t.row)
	return t.w.Write(p)
}

func (t *Terminal) translateColumn(col int) *Terminal {
	t.col += col
	if col > 0 {
		t.cursor.Right(col)
	} else if col < 0 {
		t.cursor.Left(-col)
	}
	return t
}

func (t *Terminal) translateRow(row int) *Terminal {
	t.row += row
	if row > 0 {
		t.cursor.Down(row)
	} else if row < 0 {
		t.cursor.Up(-row)
	}
	t.maxRow = max(t.maxRow, t.row)
	return t
}

func (t *Terminal) setRow(n int) *Terminal {
	n = max(0, n)

	delta := min(n-t.row, t.maxRow-t.row)
	if delta > 0 {
		t.cursor.Down(delta)
	} else if delta < 0 {
		t.cursor.Up(-delta)
	}
	t.row = n
	return t
}

func (t *Terminal) setColumn(n int) *Terminal {
	n = max(0, n)

	delta := n - t.col
	if delta > 0 {
		t.cursor.Right(delta)
	} else if delta < 0 {
		t.cursor.Left(-delta)
	}
	t.col = n
	return t
}
