package rpc

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"log"
	"net/http"
)

type validatorResponse struct {
	Validators []db.Validator
}

func validators(w http.ResponseWriter, r *http.Request) {
	var errResponse errorResponse
	validators, err := db.GetValidators()
	if err != nil {
		errResponse.Message = err.Error()
		b, err := json.Marshal(errResponse)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
		return
	}

	validatorResponse := validatorResponse{Validators: validators}
	b, err := json.Marshal(validatorResponse)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
	return

}
