package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
	"transaction-crud-svc-go-postgres/models"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var token string

var tx1 = models.Transaction{
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

var tx2 = models.Transaction{
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

var tx3 = models.Transaction{
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

var user models.User

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	user = models.User{
		UserName: os.Getenv("USER_NAME"),
		Password: os.Getenv("PASSWORD"),
	}

	tok, err := login()
	if err != nil {
		fmt.Println(err)
	}

	token = tok
	fmt.Println(token)

	resp, err := deleteAll()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Delete All status:", resp.Status)
	}
}

func Test_Create(t *testing.T) {

	_, err := create(tx1)
	if err != nil {
		fmt.Println(err)
	}

	_, err = create(tx2)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := getList()
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx1, tx2}
	validateData(t, resp, expectedData[:])

	t.Cleanup(cleanUp)
}

func Test_Get(t *testing.T) {

	_, err := create(tx1)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	_, err = create(tx2)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	resp, err := getList()
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

	resp, err = get(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx[0]}
	validateData(t, resp, expectedData[:])

	t.Cleanup(cleanUp)
}

func Test_Update(t *testing.T) {

	_, err := create(tx1)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	_, err = create(tx2)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	resp, err := getList()
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

	// Get before update
	resp, err = get(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx[0]}
	validateData(t, resp, expectedData[:])

	_, err = update(tx[0].Id, tx3)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	// Get after update
	resp, err = get(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData = [...]models.Transaction{tx3}
	validateData(t, resp, expectedData[:])

	t.Cleanup(cleanUp)
}

func Test_Delete(t *testing.T) {

	_, err := create(tx1)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	_, err = create(tx2)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	resp, err := getList()
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
	resp, err = get(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	expectedData := [...]models.Transaction{tx[0]}
	validateData(t, resp, expectedData[:])

	_, err = delete(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	// Get after delete
	resp, err = get(tx[0].Id)
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, err.Error())
	}

	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	t.Cleanup(cleanUp)
}

func login() (string, error) {
	client := &http.Client{}

	json_data, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	request, err := http.NewRequest("POST", "http://localhost:3000/login", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var token map[string]string
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return token["token"], nil
}

func create(data models.Transaction) (*http.Response, error) {
	client := &http.Client{}

	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request, err := http.NewRequest("POST", "http://localhost:3000/transactions", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func getList() (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://localhost:3000/transactions", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func get(id uint) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://localhost:3000/transactions/"+strconv.FormatUint(uint64(id), 10), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func update(id uint, data models.Transaction) (*http.Response, error) {
	client := &http.Client{}

	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request, err := http.NewRequest("PUT", "http://localhost:3000/transactions/"+strconv.FormatUint(uint64(id), 10), bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func delete(id uint) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("DELETE", "http://localhost:3000/transactions/"+strconv.FormatUint(uint64(id), 10), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func deleteAll() (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("DELETE", "http://localhost:3000/transactions", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err, "NewRequest")
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err, "Do")
		return nil, err
	}

	return resp, nil
}

func validateData(t *testing.T, resp *http.Response, expectedData []models.Transaction) {
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

	for i := 0; i < len(txs); i++ {
		assert.Equal(expectedData[i].AccountId, txs[i].AccountId, "Account Id mismatch")
		assert.Equal(expectedData[i].Amout, txs[i].Amout, "Amount mismatch")
		assert.Equal(expectedData[i].Currency, txs[i].Currency, "Currency mismatch")
		assert.Equal(expectedData[i].Description, txs[i].Description, "Description mismatch")
		assert.Equal(expectedData[i].MerchantId, txs[i].MerchantId, "Merchant Id mismatch")
		assert.Equal(expectedData[i].MerchantName, txs[i].MerchantName, "Merchant name mismatch")
		assert.Equal(expectedData[i].Metadata, txs[i].Metadata, "Metadata mismatch")
		assert.Equal(expectedData[i].Status, txs[i].Status, "Status mismatch")
		assert.True(expectedData[i].CreatedAt.Before(txs[i].CreatedAt) || expectedData[i].CreatedAt.Equal(txs[i].CreatedAt), "CreatedAt mismatch")
		assert.True(expectedData[i].ModifiedAt.Before(txs[i].ModifiedAt) || expectedData[i].ModifiedAt.Equal(txs[i].ModifiedAt), "ModifiedAt mismatch")
		assert.Equal(expectedData[i].Type, txs[i].Type, "Type mismatch")
		assert.Equal(user.UserName, txs[i].User, "Type mismatch")
	}
}

func cleanUp() {
	_, err := deleteAll()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//fmt.Println("Delete All status:", resp.Status)
	}
}
