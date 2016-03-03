package ts3

import (
	"fmt"
	"strconv"
	"strings"
)

const okError = "error id=0 msg=ok"

type Error struct {
	Id      int
	Message string
}

func NewError(id int, msg string) Error {
	return Error{
		Id:      id,
		Message: msg,
	}
}

func NewErrorOk() Error { return NewError(0, "ok") }

func NewErrorString(input string) Error {
	if input == okError {
		return NewErrorOk()
	}

	input = strings.TrimPrefix(input, "error ")
	space := strings.Index(input, " ")
	id := strings.TrimPrefix(input[:space], "id=")
	msg := strings.TrimPrefix(input[space+1:], "msg=")

	err := Error{}
	err.Id, _ = strconv.Atoi(id)
	err.Message = Unquote(msg)

	return err
}

func (e Error) Error() string {
	return fmt.Sprintf("ts3.Error id=%d msg=%q", e.Id, e.Message)
}
