package restio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/bitrise-io/go-utils/log"
)

// Handler ...
type Handler struct {
	F interface{}
}

// Satisfies the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	fnValue := reflect.ValueOf(h.F)

	argsCount := fnValue.Type().NumIn()
	if argsCount > 3 {
		respondErrorAndLog(w, "endpoint handlers can have maximum 3 arguments, received: %d", argsCount)
		return
	}

	var args []reflect.Value

	if argsCount == 3 {
		t := fnValue.Type().In(0)
		j := (reflect.New(t).Elem()).Interface()
		rv := reflect.ValueOf(&j).Elem()
		te := rv.Elem().Type().Elem()
		rv.Set(reflect.New(te))

		dat, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondErrorAndLog(w, "failed to get request body, error: %s", err)
			return
		}

		if err := json.Unmarshal(dat, j); err != nil {
			respondErrorAndLog(w, "failed to unmarshal request body, error: %s", err)
			return
		}

		args = append(args, reflect.ValueOf(j))
	}

	args = append(args, reflect.ValueOf(w))
	args = append(args, reflect.ValueOf(r))

	fnResults := fnValue.Call(args)

	if len(fnResults) > 2 {
		respondErrorAndLog(w, "endpoint handlers can have maximum 2 return values, received: %d", len(fnResults))
		return
	}

	var err error
	var model interface{}
	for _, r := range fnResults {
		if r.Type().Implements(errorInterface) {
			intf := r.Interface()
			if intf != nil {
				err = intf.(error)
			}
		} else {
			model = r.Interface()
		}
	}

	if err != nil {
		respondErrorAndLog(w, "%s", err)
		return
	}

	if model != nil {
		RespondJSON(w, http.StatusOK, model)
	}
}

func respondErrorAndLog(w http.ResponseWriter, m string, a ...interface{}) {
	m = fmt.Sprintf(m, a...)
	RespondError(w, http.StatusInternalServerError, m)
	log.Errorf(m)
}
