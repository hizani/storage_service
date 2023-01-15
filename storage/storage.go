package storage

type RequiredMissingError struct {
	msg string // name of missing field
}

func (e *RequiredMissingError) Error() string { return e.msg }
