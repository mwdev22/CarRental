package types

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // disable escaping for elements like ">", "<"...

	return encoder.Encode(v)
}

func WriteZIP(w http.ResponseWriter, zippedBuff []byte, fileName string) error {
	fileName = fileName + ".zip"

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(zippedBuff)))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	_, err := w.Write(zippedBuff)
	return err
}

func ParseBool(value string) bool {
	v, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return v != 0
}
