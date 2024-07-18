package spinner

import "time"

type Spinner struct {
	perFrame    time.Duration
	accumulator time.Duration

	lastUpdate time.Time
	frames     []string
}

func (s *Spinner) Update() string {
	now := time.Now()
	delta := now.Sub(s.lastUpdate)
	s.lastUpdate = now

	s.accumulator += delta
	frame := int(s.accumulator/s.perFrame) % len(s.frames)
	return s.frames[frame]
}

func (s *Spinner) Frame() string {
	return s.frames[int(s.accumulator/s.perFrame)%len(s.frames)]
}

func New(total time.Duration, frames ...string) *Spinner {
	return &Spinner{
		perFrame:   total / time.Duration(len(frames)),
		frames:     frames,
		lastUpdate: time.Now(),
	}
}

func Dots(cycleTime time.Duration) *Spinner {
	return New(cycleTime, "⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷")
}

func Lines(cycleTime time.Duration) *Spinner {
	return New(cycleTime, "⠂", "⠄", "⠆", "⠇", "⠋", "⠙", "⠸", "⠰", "⠠", "⠰", "⠸", "⠙", "⠋", "⠇", "⠆", "⠄")
}
