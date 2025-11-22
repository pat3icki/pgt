package server

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgproto3"
)

// BackendKeyData represents the secret key used for cancel requests
type BackendKeyData struct {
	ProcessID uint32
	SecretKey uint32
}

type Error struct {
	ErrorHanlder HandlerTag
	Err          error
	Level        string
	Comment      string
	Line         string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Current Handler: %s, Error: %s, Level: %s", e.ErrorHanlder, e.Err, e.Level)
}

type COPY int8

const (
	CopyFROM COPY = iota
	CopyTO
	CopyFAIL
	CopyData
)

type HandlerTag string

const (
	HandlerTagNone          HandlerTag = "none"
	HandlerTagStartUpMsg    HandlerTag = "start-up"
	HandlerTagQuery         HandlerTag = "query"
	HandlerTagParse         HandlerTag = "parse"
	HandlerTagBind          HandlerTag = "bind"
	HandlerTagExcute        HandlerTag = "excute"
	HandlerTagSync          HandlerTag = "sync"
	HandlerTagCopy          HandlerTag = "copy"
	HandlerTagCancelRequest HandlerTag = "cancel-request"
	HandlerTagFunctionCall  HandlerTag = "function-call"
	HandlerTagTerminate     HandlerTag = "terminate"
	HandlerTagClose         HandlerTag = "close"
	HandlerTagServerContext HandlerTag = "context"
	HandlerTagEOF           HandlerTag = "connection closed"
)

func (h HandlerTag) String() string {
	return string(h)
}

type Handler interface {
	HandlerStartupMessage(StartUpCtx) error

	HandlerQuery(Context, *pgproto3.Query) error
	HandlerParse(Context, *pgproto3.Parse) error
	HandlerBind(Context, *pgproto3.Bind) error
	HandlerExecute(Context, *pgproto3.Execute) error
	HandlerSync(Context, *pgproto3.Sync) error
	HandlerCopy(COPY, Context) error
	HandlerCancelRequest(Context, *pgproto3.CancelRequest) error
	HandlerFunctionCall(Context, *pgproto3.FunctionCall) error
	HandlerTerminate(Context, *pgproto3.Terminate) error
	HandlerClose(Context, *pgproto3.Close) error

	HandleError(*Error)
}
