package server

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgproto3"
)

// Handler StartUp Message
func SendStartUpMessage(io *_IO) (err error) {
	err = io.Send(&pgproto3.StartupMessage{
		ProtocolVersion: 3,
		Parameters: map[string]string{
			"user":        "pgserver",
			"application": "pg-server",
		},
	})
	if err != nil {
		return
	}

	for {
		var raw_msg pgproto3.BackendMessage

		raw_msg, err = io.Receive() // something is wrong here
		if err != nil {
			return
		}

		switch raw_msg.(type) {
		case *pgproto3.AuthenticationOk:

			return

		}

	}
}

// Handler SSLRequest
func SendSslRequest(io *_IO) (err error) {
	err = io.Send(&pgproto3.SSLRequest{})
	if err != nil {
		return
	}

	rbuf := make([]byte, 1)
	_, err = io.Read(rbuf)
	if err != nil {
		return
	}

	str := bytes.NewBuffer(rbuf).String()
	if !(strings.Contains(str, "N") ||
		strings.Contains(str, "S")) {
		return errors.New("invaild response from ssl request")
	}
	return
}

// Handler Bind

func SendBind(io *_IO) (err error) {
	ret := pgproto3.Bind{
		DestinationPortal:    "test", // unnamed portal
		PreparedStatement:    "get_user",
		ParameterFormatCodes: []int16{0, 0}, // both parameters in text format
		Parameters: [][]byte{
			[]byte("123"),  // value for $1 (id)
			[]byte("john"), // value for $2 (name)
		},
		ResultFormatCodes: []int16{0}, // all result columns in text format
	}
	err = io.Send(&ret)
	if err != nil {
		return errors.New(":fyh")
	}

	for {
		var raw_msg pgproto3.BackendMessage

		raw_msg, err = io.Receive()
		if err != nil {
			err = fmt.Errorf("experted type *pgproto3.BindComplete")
			return

		}
		switch raw_msg.(type) {
		case *pgproto3.BindComplete:
			return nil

		default:
			err = fmt.Errorf("experted type *pgproto3.BindComplete but got %s", reflect.TypeOf(raw_msg))
			return
		}
	}

}

// Handler Parse
func SendParse(io *_IO) (err error) {
	ret := pgproto3.Parse{
		Name:  "test",
		Query: "SELECT 1",
	}
	err = io.Send(&ret)
	if err != nil {
		return
	}

	for {
		var raw_msg pgproto3.BackendMessage

		raw_msg, err = io.Receive()
		if err != nil {
			return fmt.Errorf("experted type *pgproto3.ParseComplete")
		}
		switch raw_msg.(type) {
		case *pgproto3.ParseComplete:
			return nil

		default:
			err = fmt.Errorf("experted type *pgproto3.ParseComplete")
			return
		}
	}
}

// Handler Query
func SendQuery(io *_IO) (err error) {
	ret := pgproto3.Query{
		String: "SELECT 1",
	}
	err = io.Send(&ret)
	if err != nil {
		return err
	}
	for {
		var raw_msg pgproto3.BackendMessage

		raw_msg, err = io.Receive()
		if err != nil {
			return err
		}

		switch raw_msg.(type) {
		case *pgproto3.CommandComplete:
		case *pgproto3.ErrorResponse:
		case *pgproto3.ReadyForQuery:
			return nil

		default:
			return errors.New("invaild response from query handler")

		}
	}
}
