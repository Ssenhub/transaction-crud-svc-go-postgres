package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"transaction-crud-svc-go-postgres/middlewares"
	"transaction-crud-svc-go-postgres/models"
	"transaction-crud-svc-go-postgres/storage"

	"github.com/go-chi/chi/v5"
)

type DataStore struct {
	Store storage.TransactionStore
}

var DataHandler DataStore

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	txModels, err := DataHandler.Store.GetList()

	if err != nil {
		WriteError(w, "Failed to find books. "+err.Error(), 512)
	} else {
		resp, err := json.Marshal(txModels)

		if err != nil {
			WriteError(w, "Failed to serialize books. "+err.Error(), 513)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var txModel = models.Transaction{}
	var err error
	var id uint64

	if id, err = strconv.ParseUint(chi.URLParam(r, "id"), 10, 0); err == nil {
		txModel, err = DataHandler.Store.Get(uint(id))
	} else {
		WriteError(w, "Failed to retrieve id from URL. "+err.Error(), 514)
		return
	}

	if err != nil {
		WriteError(w, "Id ("+strconv.FormatUint(id, 10)+") not found. "+err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(txModel)

	if err != nil {
		WriteError(w, "Failed to serialize books", 515)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var id uint64
	var curTx models.Transaction

	if id, err = strconv.ParseUint(chi.URLParam(r, "id"), 10, 0); err == nil {
		curTx, err = DataHandler.Store.Get(uint(id))
	} else {
		WriteError(w, "Failed to retrieve id from URL. "+err.Error(), 514)
		return
	}

	if err != nil {
		WriteError(w, "Id ("+string(id)+") not found. "+err.Error(), http.StatusNotFound)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		WriteError(w, "Failed to read body. "+err.Error(), http.StatusBadRequest)
	} else {
		//fmt.Println("Body: ", string(body))
	}

	var putTx models.Transaction
	err2 := json.Unmarshal(body, &putTx)

	if err2 != nil {
		WriteError(w, "Invalid json. "+err.Error(), http.StatusBadRequest)
	} else {
		//fmt.Println(putTx)
	}

	putTx.Id = uint(id)
	putTx.CreatedAt = curTx.CreatedAt
	putTx.User = middlewares.GetUser(r)
	putTx.ModifiedAt = time.Now()

	err = DataHandler.Store.Update(putTx)
	if err != nil {
		WriteError(w, "Failed to save transaction. "+err.Error(), 514)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		WriteError(w, "Failed to read body. "+err.Error(), http.StatusBadRequest)
	} else {
		//fmt.Println("Body: ", string(body))
	}

	var tx models.Transaction
	err2 := json.Unmarshal(body, &tx)
	if err2 != nil {
		WriteError(w, "Invalid json. "+err2.Error(), http.StatusBadRequest)
	} else {
		//fmt.Println(tx)
	}

	tx.User = middlewares.GetUser(r)
	tx.CreatedAt = time.Now()
	tx.ModifiedAt = tx.CreatedAt

	createErr := DataHandler.Store.Create(tx)
	if createErr != nil {
		WriteError(w, "Failed to create book. "+createErr.Error(), 516)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var id uint64

	if id, err = strconv.ParseUint(chi.URLParam(r, "id"), 10, 0); err == nil {
		err = DataHandler.Store.Delete(uint(id))
	} else {
		WriteError(w, "Failed to retrieve id from URL. "+err.Error(), 514)
		return
	}

	if err != nil {
		WriteError(w, "Id ("+strconv.FormatUint(id, 10)+") not found. "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	err := DataHandler.Store.DeleteAll()

	if err != nil {
		WriteError(w, "Delete failed. "+err.Error(), 517)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func WriteError(w http.ResponseWriter, errorMsg string, statusCode int) {
	fmt.Println(errorMsg)
	http.Error(w, errorMsg, statusCode)
}
