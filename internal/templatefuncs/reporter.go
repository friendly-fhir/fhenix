package templatefuncs

type Reporter interface {
	Report(err error)
}

type ReporterFunc func(error)

func (f ReporterFunc) Report(err error) {
	f(err)
}
