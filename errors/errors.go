package errors

import "net/http"

type Op string
type ErrKind int
type ErrClientMessage string
type ErrPayload struct {
	Data interface{}
}

const (
	KindBadRequest    ErrKind = http.StatusBadRequest
	KindUnauthorized  ErrKind = http.StatusUnauthorized
	KindNotFound      ErrKind = http.StatusNotFound
	KindUnprocessable ErrKind = http.StatusUnprocessableEntity
	KindServerError   ErrKind = http.StatusInternalServerError
)

type Error struct {
	Op            Op               // Operation (method) that failed
	Kind          ErrKind          // Kind of error
	ClientMessage ErrClientMessage // Message for client
	Payload       ErrPayload       // Payload
	Err           error            // Actual error

	// (optional) application specific data
}

func (e *Error) Error() string {
	return e.Err.Error()
}

/* Example:

func main() {
	err := CreateUser()
	if err != nil {
		fmt.Println(errors.Kind(err))
		fmt.Println(errors.ClientMessage(err))
		fmt.Println(errors.Payload(err))

		er := err.(*errors.Error)
		fmt.Println(errors.Ops(er))
		return
	}

	fmt.Println("Success")
}

func CreateUser() error {
	const op errors.Op = "CreateUser"

	err := InsertUser()
	return errors.E(err, op)
}

func InsertUser() error {
	const op errors.Op = "InsertUser"

	var kind errors.ErrKind = errors.KindUnprocessable
	var clientMessage errors.ErrClientMessage = "Username should be unique"
	var payload errors.ErrPayload = errors.ErrPayload{
		Data: map[string]interface{}{
			"id": 1,
		},
	}

	err := goerrors.New("Duplicate username")
	return errors.E(op, kind, clientMessage, payload, err)
}
*/
