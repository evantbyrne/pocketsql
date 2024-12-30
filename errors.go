package pocketsql

import (
	"log"
	"net/http"

	"github.com/evantbyrne/trance"
)

type errorSql struct {
	Message string
}

func (err errorSql) Error() string {
	if err.Message != "" {
		return err.Message
	}
	return "Bad request"
}

func (err errorSql) Status() int {
	return http.StatusBadRequest
}

func errorHandler(strand *trance.Strand) {
	if strand.Error == nil {
		return
	}
	if _, ok := strand.Error.(errorSql); ok {
		strand.Response.Header().Add("Content-Type", "text/vnd.turbo-stream.html")
		errorResponse(strand.Response, errorSqlStreamTemplate.Execute(strand.Response, map[string]any{
			"message": strand.Error.Error(),
		}))
		return
	}
	errorResponse(strand.Response, strand.Error)
}

func errorResponse(response http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("%T: %s\n", err, err)
		status := http.StatusInternalServerError
		if errWithStatus, ok := err.(trance.ErrorWithStatus); ok {
			status = errWithStatus.Status()
		}
		log.Println(err)
		http.Error(response, err.Error(), status)
	}
}
