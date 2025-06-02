package tests

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
	"transaction-crud-svc-go-postgres/handlers"
	"transaction-crud-svc-go-postgres/middlewares"
	"transaction-crud-svc-go-postgres/models"
	"transaction-crud-svc-go-postgres/storage"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var tx1_mock = models.Transaction{
	Id:           1,
	AccountId:    "act1",
	Type:         models.Credit,
	Amout:        123.45,
	Currency:     "USD",
	CreatedAt:    time.Now(),
	ModifiedAt:   time.Now(),
	Description:  "test desc 1",
	Status:       models.Pending,
	MerchantId:   "Merch_1",
	MerchantName: "MerchName_1",
	Metadata:     "{'channnel': 'mobile_app1', 'location': 'Seattle, WA'}",
}

var tx2_mock = models.Transaction{
	Id:           2,
	AccountId:    "act2",
	Type:         models.Debit,
	Amout:        456.78,
	Currency:     "USD",
	CreatedAt:    time.Now(),
	ModifiedAt:   time.Now(),
	Description:  "test desc 2",
	Status:       models.Failed,
	MerchantId:   "Merch_2",
	MerchantName: "MerchName_2",
	Metadata:     "{'channnel': 'mobile_app2', 'location': 'San francisco, CA'}",
}

var tx3_mock = models.Transaction{
	Id:           3,
	AccountId:    "act3",
	Type:         models.Credit,
	Amout:        789.50,
	Currency:     "USD",
	CreatedAt:    time.Now(),
	ModifiedAt:   time.Now(),
	Description:  "test desc 3",
	Status:       models.Pending,
	MerchantId:   "Merch_3",
	MerchantName: "MerchName_3",
	Metadata:     "{'channnel': 'mobile_app3', 'location': 'Portland, OR'}",
}

var mockUser models.User

var mockStore storage.MockTransactionStore

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	mockUser = models.User{
		UserName: os.Getenv("USER_NAME"),
		Password: os.Getenv("PASSWORD"),
	}

	mockStore = storage.MockTransactionStore{}

	handlers.DataHandler = handlers.DataStore{
		Store: &mockStore,
	}

	mockCleanUp()
}

func Test_PostHandler(t *testing.T) {

	_, err := postHandler(tx1_mock)
	if err != nil {
		fmt.Println(err)
	}

	_, err = postHandler(tx2_mock)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := getListHandler()
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx1_mock, tx2_mock}
	validateMockData(t, resp, expectedData[:])

	t.Cleanup(mockCleanUp)
}

func Test_GetHandler(t *testing.T) {

	_, err := postHandler(tx1_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	_, err = postHandler(tx2_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	resp, err := getListHandler()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	//jsonBody := string(body)
	//fmt.Println(jsonBody)

	var tx []models.Transaction
	err = json.Unmarshal(body, &tx)
	if err != nil {
		fmt.Println("Invalid json:", err.Error())
		assert.Fail(t, "Invalid json:", err.Error())
	}

	resp, err = getHandler(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx[0]}
	validateMockData(t, resp, expectedData[:])

	t.Cleanup(mockCleanUp)
}

func Test_PutHandler(t *testing.T) {

	_, err := postHandler(tx1_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	_, err = postHandler(tx2_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	resp, err := getListHandler()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	var tx []models.Transaction
	err = json.Unmarshal(body, &tx)
	if err != nil {
		fmt.Println("Invalid json:", err.Error())
		assert.Fail(t, "Invalid json:", err.Error())
	}

	// Get before update
	resp, err = getHandler(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx[0]}
	validateMockData(t, resp, expectedData[:])

	_, err = putHandler(tx[0].Id, tx3_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	// Get after update
	resp, err = getHandler(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	tx3DeepCopy := DeepCopy(tx3_mock)
	tx3DeepCopy.Id = tx[0].Id
	expectedData = [...]models.Transaction{tx3DeepCopy}
	validateMockData(t, resp, expectedData[:])

	t.Cleanup(mockCleanUp)
}

func Test_DeleteHandler(t *testing.T) {

	_, err := postHandler(tx1_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	_, err = postHandler(tx2_mock)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	resp, err := getListHandler()
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	//jsonBody := string(body)
	//fmt.Println(jsonBody)

	var tx []models.Transaction
	err = json.Unmarshal(body, &tx)
	if err != nil {
		fmt.Println("Invalid json:", err.Error())
		assert.Fail(t, "Invalid json:", err.Error())
	}

	// Get before delete
	resp, err = getHandler(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx[0]}
	validateMockData(t, resp, expectedData[:])

	_, err = deleteHandler(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	// Get after delete
	resp, err = getHandler(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	t.Cleanup(mockCleanUp)
}

func postHandler(data models.Transaction) (*http.Response, error) {
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/transactions", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req = req.WithContext(context.WithValue(req.Context(), middlewares.UserKey, mockUser.UserName))

	rr := httptest.NewRecorder()

	handlers.PostHandler(rr, req)

	return rr.Result(), nil
}

func getListHandler() (*http.Response, error) {
	req, err := http.NewRequest("GET", "http://localhost:3000/transactions", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req = req.WithContext(context.WithValue(req.Context(), middlewares.UserKey, mockUser.UserName))

	rr := httptest.NewRecorder()

	handlers.GetListHandler(rr, req)

	return rr.Result(), nil
}

func getHandler(id uint) (*http.Response, error) {
	strId := strconv.FormatUint(uint64(id), 10)
	req := httptest.NewRequest("GET", "http://localhost:3000/transactions/"+strId, nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", strId)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	req = req.WithContext(context.WithValue(req.Context(), middlewares.UserKey, mockUser.UserName))

	rr := httptest.NewRecorder()

	handlers.GetHandler(rr, req)

	return rr.Result(), nil
}

func putHandler(id uint, data models.Transaction) (*http.Response, error) {
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	strId := strconv.FormatUint(uint64(id), 10)

	req := httptest.NewRequest("PUT", "http://localhost:3000/transactions/"+strId, bytes.NewBuffer(json_data))

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", strId)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
	req = req.WithContext(context.WithValue(req.Context(), middlewares.UserKey, mockUser.UserName))

	rr := httptest.NewRecorder()

	handlers.PutHandler(rr, req)

	return rr.Result(), nil
}

func deleteHandler(id uint) (*http.Response, error) {
	strId := strconv.FormatUint(uint64(id), 10)

	req := httptest.NewRequest("DELETE", "http://localhost:3000/transactions/"+strconv.FormatUint(uint64(id), 10), nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", strId)

	ctx := req.Context()
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
	ctx = context.WithValue(ctx, middlewares.UserKey, mockUser.UserName)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handlers.DeleteHandler(rr, req)

	return rr.Result(), nil
}

func deleteAllHandler() (*http.Response, error) {
	req, err := http.NewRequest("DELETE", "http://localhost:3000/transactions", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req = req.WithContext(context.WithValue(req.Context(), middlewares.UserKey, mockUser.UserName))

	rr := httptest.NewRecorder()

	handlers.DeleteAllHandler(rr, req)

	return rr.Result(), nil
}

func validateMockData(t *testing.T, resp *http.Response, expectedData []models.Transaction) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	var txs []models.Transaction
	err = json.Unmarshal(body, &txs)
	if err != nil {
		var tx models.Transaction
		err = json.Unmarshal(body, &tx)
		if err != nil {
			fmt.Println("Invalid json:", err.Error())
			assert.Fail(t, "Invalid json:", err.Error())
			return
		}

		txs = []models.Transaction{tx}
	}

	assert := assert.New(t)
	assert.Equal(len(txs), len(expectedData), "Length mismatch")

	expectedTxMap := make(map[uint]models.Transaction)
	for _, tx := range expectedData {
		expectedTxMap[tx.Id] = tx
	}

	resultTxMap := make(map[uint]models.Transaction)
	for _, tx := range txs {
		resultTxMap[tx.Id] = tx
	}

	for _, expectedTx := range expectedTxMap {
		tx, ok := resultTxMap[expectedTx.Id]

		assert.True(ok, "Tx id not exist")
		assert.Equal(expectedTx.AccountId, tx.AccountId, "Account Id mismatch")
		assert.Equal(expectedTx.Amout, tx.Amout, "Amount mismatch")
		assert.Equal(expectedTx.Currency, tx.Currency, "Currency mismatch")
		assert.Equal(expectedTx.Description, tx.Description, "Description mismatch")
		assert.Equal(expectedTx.MerchantId, tx.MerchantId, "Merchant Id mismatch")
		assert.Equal(expectedTx.MerchantName, tx.MerchantName, "Merchant name mismatch")
		assert.Equal(expectedTx.Metadata, tx.Metadata, "Metadata mismatch")
		assert.Equal(expectedTx.Status, tx.Status, "Status mismatch")
		assert.True(expectedTx.CreatedAt.Before(tx.CreatedAt) || expectedTx.CreatedAt.Equal(tx.CreatedAt), "CreatedAt mismatch")
		assert.True(expectedTx.ModifiedAt.Before(tx.ModifiedAt) || expectedTx.ModifiedAt.Equal(tx.ModifiedAt), "ModifiedAt mismatch")
		assert.Equal(expectedTx.Type, tx.Type, "Type mismatch")
		assert.Equal(mockUser.UserName, tx.User, "Type mismatch")
	}
}

func mockCleanUp() {
	resp, err := deleteAllHandler()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = ReadResponse(resp)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		//fmt.Println("Delete All status:", resp.Status, body)
	}
}

func ReadResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		return string(body), nil
	}
}

func DeepCopy[T any](src T) T {
	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(src)
	var dst T
	_ = gob.NewDecoder(&buf).Decode(&dst)
	return dst
}
