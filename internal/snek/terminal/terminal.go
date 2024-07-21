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
func (c *Terminal) Close() {
	c.cursor.Show()
}

// HideCursor hides the cursor.
func (c *Terminal) HideCursor() {
	c.cursor.Hide()
}

// ShowCursor shows the cursor.
func (c *Terminal) ShowCursor() {
	c.cursor.Show()
}

// Println prints the given arguments to the terminal, followed by a newline.
func (c *Terminal) Println(n int, args ...any) (int, error) {
	formatted := fmt.Sprintln(args...)

	c.lock.Lock()
	defer c.lock.Unlock()
	c.setRow(n)
	c.setColumn(0)
	c.cursor.ClearLine()
	return c.write([]byte(formatted))
}

// Print prints the given arguments to the terminal.
func (c *Terminal) Print(n int, args ...any) (int, error) {
	formatted := fmt.Sprint(args...)

	c.lock.Lock()
	defer c.lock.Unlock()
	c.setRow(n)
	c.setColumn(0)
	return c.write([]byte(formatted))
}

// Printf prints the given arguments to the terminal, using the given format.
func (c *Terminal) Printf(n int, format string, args ...any) (int, error) {
	formatted := fmt.Sprintf(format, args...)

	c.lock.Lock()
	defer c.lock.Unlock()
	c.setRow(n)
	c.setColumn(0)
	return c.write([]byte(formatted))
}

// Offset returns the current row and column of the cursor, offset from where
// the cursor was initially created.
func (c *Terminal) Offset() (int, int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.row, c.col
}

// Row returns the current row of the cursor.
func (c *Terminal) Row() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.row
}

// Column returns the current column of the cursor.
func (c *Terminal) Column() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.col
}

// Write writes the given bytes to the cursor, updating the row and column.
func (c *Terminal) Write(p []byte) (int, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.write(p)
}

// Up moves the cursor up the given number of rows. It cannot exceed the point
// where the cursor was initially created.
func (c *Terminal) Up(rows int) *Terminal {
	if rows < 0 {
		return c.Down(-rows)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	offset := min(max(rows, 0), c.row)
	c.row -= offset
	c.cursor.Up(offset)
	return c
}

func (c *Terminal) Down(rows int) *Terminal {
	if rows < 0 {
		return c.Up(-rows)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.row += rows
	c.cursor.Down(rows)
	c.maxRow = max(c.maxRow, c.row)
	return c
}

func (c *Terminal) Left(cols int) *Terminal {
	if cols < 0 {
		return c.Right(-cols)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	offset := min(max(cols, 0), c.col)
	c.col -= offset
	c.cursor.Left(offset)
	return c
}

func (c *Terminal) Right(cols int) *Terminal {
	if cols < 0 {
		return c.Left(-cols)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.col += cols
	c.cursor.Right(cols)
	return c
}

func (c *Terminal) Top() *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.setRow(0)
	return c
}

func (c *Terminal) Bottom() *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.translateRow(c.maxRow - c.row)
}

func (c *Terminal) Move(row, col int) *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.translateColumn(row).translateRow(col)
}

func (c *Terminal) MoveColumn(col int) *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.translateColumn(col)
}

func (c *Terminal) MoveRow(row int) *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.translateRow(row)
}

func (c *Terminal) SetPosition(row, col int) *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.setRow(row).setColumn(col)
}

func (c *Terminal) SetRow(n int) *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.setRow(n)
}

func (c *Terminal) SetColumn(n int) *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.setColumn(n)
}

func (c *Terminal) Clear() *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cursor.Clear()
	c.col--
	return c
}

func (c *Terminal) ClearLine() *Terminal {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cursor.ClearLine()
	c.col = 0
	return c
}

func (c *Terminal) write(p []byte) (int, error) {
	for _, b := range p {
		if b == '\n' {
			c.row++
			c.col = 0
			continue
		}
		if b == '\r' {
			c.col = 0
			continue
		}
		if b == '\b' {
			c.col--
			continue
		}
		c.col++
	}
	c.maxRow = max(c.maxRow, c.row)
	return c.w.Write(p)
}

func (c *Terminal) translateColumn(col int) *Terminal {
	c.col += col
	if col > 0 {
		c.cursor.Right(col)
	} else if col < 0 {
		c.cursor.Left(-col)
	}
	return c
}

func (c *Terminal) translateRow(row int) *Terminal {
	c.row += row
	if row > 0 {
		c.cursor.Down(row)
	} else if row < 0 {
		c.cursor.Up(-row)
	}
	c.maxRow = max(c.maxRow, c.row)
	return c
}

func (c *Terminal) setRow(n int) *Terminal {
	n = max(0, n)

	delta := min(n-c.row, c.maxRow-c.row)
	if delta > 0 {
		c.cursor.Down(delta)
	} else if delta < 0 {
		c.cursor.Up(-delta)
	}
	c.row = n
	return c
}

func (c *Terminal) setColumn(n int) *Terminal {
	n = max(0, n)

	delta := n - c.col
	if delta > 0 {
		c.cursor.Right(delta)
	} else if delta < 0 {
		c.cursor.Left(-delta)
	}
	c.col = n
	return c
}
