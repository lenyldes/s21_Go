package handler

import (
	"encoding/json"
	"errors"
	"io"
)

func decodeJSONBody(body io.Reader, request any) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(request); err != nil {
		return err
	}

	// После первого объекта допустимы только пробелы и конец тела запроса.
	var extra json.RawMessage
	if err := decoder.Decode(&extra); errors.Is(err, io.EOF) {
		return nil
	}
	return errors.New("тело запроса должно содержать только один JSON-объект")
}
