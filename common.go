package restio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/bitrise-io/go-utils/log"
)

// Panic ...
func Panic(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

// RespondJSON ...
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// RespondError ...
func RespondError(w http.ResponseWriter, sc int, f string, args ...interface{}) {
	RespondJSON(w, sc, map[string]string{"error": fmt.Sprintf(f, args...)})
}

func dynamicFunc(f interface{}) error {
	fnValue := reflect.ValueOf(f)
	arguments := []reflect.Value{}
	fnResults := fnValue.Call(arguments)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	for _, r := range fnResults {
		if r.Type().Implements(errorInterface) {
			intf := r.Interface()
			if intf != nil {
				return intf.(error)
			}
			return nil
		}
	}
	return fmt.Errorf("function has no error return value")
}
