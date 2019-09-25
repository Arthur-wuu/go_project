package api

type Error struct{
	Err int
	ErrMsg string
}

func NewError(err int, msg string) *Error {
	return &Error{Err:err, ErrMsg:msg}
}