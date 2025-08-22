package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"golang/tutorial/todo/internal/apperr"
)

// NOTE: ヘルパー関数のためanyを許容
func WriteJSON(w http.ResponseWriter, statusCode int, v ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(v) > 0 {
		_ = json.NewEncoder(w).Encode(v[0])
	}
}

type ErrorJSON struct {
	Code    string
	Message string
}

func writeErrorJSON(w http.ResponseWriter, statusCode int, errCode string, message string) {
	WriteJSON(w, statusCode, ErrorJSON{Code: errCode, Message: message})
}
func WriteError(w http.ResponseWriter, err error) {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		writeErrorJSON(w, ae.Code.HTTPStatus(), ae.Code.String(), ae.Error())
		return
	}
	writeErrorJSON(w, http.StatusInternalServerError, apperr.CodeUnknown.String(), err.Error())
}

func ReadJSON(r *http.Request, v any) error {
	// Content-Type
	ct := r.Header.Get("Content-Type")
	if ct != "" && !strings.HasPrefix(ct, "application/json") {
		return apperr.E(apperr.CodeInvalid, "Content-Type must be application/json", nil)
	}

	defer r.Body.Close()

	/*
		Reference: O'REILLY「実用GO言語」8.1 p.169
		json.Decoder(): io.Readerインターフェースを満たしている型（os.Stdin, http.Response.Body...）
		json.Unmarshal(): []byteを扱う場合
	*/
	decoder := json.NewDecoder(r.Body)

	/*
		Reference: O'REILLY「実用GO言語」8.1 p.174
		デコード時に未知のフィールドがある場合にエラーとする。
		エラーメッセージ: json: unknown field radius
	*/
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF):
			return apperr.E(apperr.CodeInvalid, "request body is empty", err)
		case errors.As(err, &syntaxErr):
			return apperr.E(apperr.CodeInvalid, "invalid JSON syntax", err)
		case errors.As(err, &unmarshalTypeErr):
			return apperr.E(apperr.CodeInvalid, "invalid field type", err)
		default:
			return apperr.E(apperr.CodeInvalid, "invalid JSON", err)
		}
	}

	// 余りがないかチェック（配列など2個目のトークンがないか）
	if decoder.More() {
		return apperr.E(apperr.CodeInvalid, "multiple JSON values in body", nil)
	}

	return nil
}

func ParseID(r *http.Request, key string) (string, error) {
	raw := r.PathValue(key)
	if raw == "" {
		return "", apperr.E(apperr.CodeInvalid, "missing ID in path", nil)
	}
	return raw, nil
}
