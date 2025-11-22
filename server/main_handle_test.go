// BSD 3-Clause License

// Copyright (c) 2025, Patrick Innocent

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package server

import (
	"errors"
	"io"
	"log"

	"github.com/jackc/pgx/v5/pgproto3"
)

type TestHandlerStruct struct {
	name  string
	query string
}

func NewTestHandler() Handler {
	return &TestHandlerStruct{}
}

func (h *TestHandlerStruct) HandlerStartupMessage(StartUpCtx StartUpCtx) error {
	// 	var parameter map[string]string

	//	return startup.DefaultStartUpMsgHandler(&startup.StartUpParams{
	//		Ctx:      StartUpCtx,
	//		AuthType: pgproto3.AuthTypeOk,
	//		BackendKey: BackendKeyData{
	//			ProcessID: 23456,
	//			SecretKey: 876543,
	//		},
	//	}, parameter)
	return nil
}

func (h *TestHandlerStruct) HandlerQuery(ctx Context, msg *pgproto3.Query) error {
	if msg.String == "" {
		ctx.WriteMsg(&pgproto3.EmptyQueryResponse{})
	}
	if msg.String != "SELECT 1" {
		ctx.WriteMsg(&pgproto3.ErrorResponse{
			TableName: "pg-server",
			Severity:  "Fatal",
			Message:   "Query must be SELECT 1",
		})
	} else {
		ctx.WriteMsg(&pgproto3.CommandComplete{
			CommandTag: []byte("SELECT 1"),
		})
	}
	ctx.ReadyForQuery('I')
	return nil
}

func (h *TestHandlerStruct) HandlerParse(ctx Context, msg *pgproto3.Parse) error {
	h.name = msg.Name
	h.query = msg.Query

	buf, err := (&pgproto3.ParseComplete{}).Encode(nil)
	if err != nil {
		return err
	}
	ctx.Write(buf)
	// ctx.WriteMsg(&pgproto3.ParseComplete{})

	return nil
}

func (h *TestHandlerStruct) HandlerBind(ctx Context, msg *pgproto3.Bind) error {

	// fmt.Println("whats happening ?:", msg.DestinationPortal)
	if msg == nil {
		return errors.New("ERROR!")
	}
	// if h.name != msg.DestinationPortal {
	// 	return errors.New("ERROR! - 2")
	// }

	ctx.WriteMsg(&pgproto3.BindComplete{})
	// buf, err := (&pgproto3.BindComplete{}).Encode(nil)
	// if err != nil {
	// 	panic(err)
	// }
	// ctx.Write(buf)
	return nil
}

func (h *TestHandlerStruct) HandlerExecute(ctx Context, msg *pgproto3.Execute) error {
	return nil
}

func (h *TestHandlerStruct) HandlerSync(ctx Context, msg *pgproto3.Sync) error {
	return nil
}

func (h *TestHandlerStruct) HandlerCopy(COPY COPY, ctx Context) error {
	return nil
}
func (h *TestHandlerStruct) HandlerCancelRequest(ctx Context, msg *pgproto3.CancelRequest) error {
	return nil
}
func (h *TestHandlerStruct) HandlerFunctionCall(ctx Context, msg *pgproto3.FunctionCall) error {
	return nil
}

func (h *TestHandlerStruct) HandlerTerminate(ctx Context, msg *pgproto3.Terminate) error {
	return nil
}
func (h *TestHandlerStruct) HandlerClose(ctx Context, msg *pgproto3.Close) error {
	return nil
}

func (h *TestHandlerStruct) HandleError(err *Error) {
	if err.Err != io.EOF || err.Err != io.ErrUnexpectedEOF {
		return
	}
	log.Fatalf(`
	   Server Error: 
	Current Server Handler: %s, 
	Level: %s, 
	Error: %s, 
	Comment: %s, 
	`, err.ErrorHanlder, err.Level, err.Err.Error(), err.Comment)
}
