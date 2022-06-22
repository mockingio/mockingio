package api

import (
	"github.com/smockyio/smocky/engine/mock"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/engine/persistent"
)

func GetMocksHandler(w http.ResponseWriter, r *http.Request) {
	db := persistent.GetDefault()
	mocks, err := db.GetMocks(r.Context())
	if err != nil {
		log.WithError(err).Error("get configs")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, mocks)
}

func CreateMockHandler(w http.ResponseWriter, r *http.Request) {
	db := persistent.GetDefault()
	m := mock.New()

	if err := db.SetMock(r.Context(), m); err != nil {
		log.WithError(err).Error("create new mock")
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusCreated, map[string]any{"id": m.ID})
}
