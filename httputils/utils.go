package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//GetIP gets the remote IP from a request
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

//LogDivider prints a divider in the console
func LogDivider() {
	log.Printf("------------------------------")
}

//MalformedRequest struct that allow to handle the problems when reading decoding the Filters json info
type MalformedRequest struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func JSONPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

//DecodeJSONFromRequest decodes a Json request and returns an error if needed
func DecodeJSONFromRequest(w http.ResponseWriter, r *http.Request, dst interface{}, origin string) error {
	if r.Header.Get("Content-Type") != "application/json" {
		msg := fmt.Sprintf(`{"message":[%v] Incorrect Content-Type}`, origin)
		log.Println(msg)
		return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
	}
	//read the body and log it
	body, errr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if errr != nil {
		panic(errr)
	}
	log.Printf("[%v] Search request received:", origin)
	log.Println("\n", JSONPrettyPrint(string(body)))
	//create a new body with the same information and continue processing the request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("[%v] Request body contains badly-formed JSON (at position %d)", origin, syntaxError.Offset)
			log.Println(msg)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("[%v] Request body contains badly-formed JSON", origin)
			log.Println(msg)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("[%v] Request body contains an invalid value for the %q field (at position %d)", origin, unmarshalTypeError.Field, unmarshalTypeError.Offset)
			log.Println(msg)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("[%v] Request body contains unknown field %s", origin, fieldName)
			log.Println(msg)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := fmt.Sprintf("[%v] Request body must not be empty", origin)
			log.Println(msg)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := fmt.Sprintf("[%v] Request body must not be larger than 1MB", origin)
			log.Println(msg)
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return err
		}
	}

	if dec.More() {
		msg := fmt.Sprintf("[%v] Request body must only contain a single JSON object", origin)
		log.Println(msg)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	log.Printf("[%v] Request body parsed correctly", origin)

	return nil
}
