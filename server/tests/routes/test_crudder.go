package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"ideashare/models"
	"ideashare/tests/testutil"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestCrudEndpoints[T models.BaseModel](t *testing.T, prefix string, emptyModel func() T, update func(T)) {
	created := testCreate(t, prefix, emptyModel)
	testRead(t, prefix, created, emptyModel)
	testUpdate(t, prefix, created, update, emptyModel)
	testDelete(t, prefix, created)
	testReadAll[T](t, prefix, func(count int) {
		for i := 0; i < count; i++ {
			testCreate(t, prefix, emptyModel)
		}
	})
}

func TestCrudEndpointsWithoutAll[T models.BaseModel](t *testing.T, prefix string, emptyModel func() T, update func(T)) {
	created := testCreate(t, prefix, emptyModel)
	testRead(t, prefix, created, emptyModel)
	testUpdate(t, prefix, created, update, emptyModel)
	testDelete(t, prefix, created)
}

func testCreate[T models.BaseModel](t *testing.T, prefix string, emptyModel func() T) T {
	fmt.Println("Testing create new record for prefix:" + " " + prefix)
	populatedModel := testutil.MakeFake(emptyModel())
	reqBody, err := json.Marshal(populatedModel)
	if err != nil {
		panic(err)
	}
	res, err := http.Post(testutil.TestApiUrl(prefix), "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	body := emptyModel()
	err = json.Unmarshal(buf, body)
	if err != nil {
		fmt.Println("Failed to unmarshal: ", string(buf))
		panic(err)
	}
	assert.NotZero(t, body.GetID(), "ID should not be zero")
	assert.NotZero(t, body.GetCreatedBy(), "CreatedBy should not be zero")
	assert.NotNil(t, body.GetCreatedAt(), "CreatedAt should not be nil")
	assert.NotNil(t, body.GetUpdatedAt(), "UpdatedAt should not be nil")
	return body
}

func testRead[T models.BaseModel](t *testing.T, prefix string, created T, emptyModel func() T) {
	fmt.Println("Testing read for prefix:" + " " + prefix + " with id: " + strconv.Itoa(created.GetID()) + "")
	res, err := http.Get(testutil.TestApiUrl(prefix + "/" + strconv.Itoa(created.GetID())))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	body := emptyModel()
	err = json.Unmarshal(buf, body)
	assert.Equal(t, created, body, "Body should be equal to created model")
}

func testUpdate[T models.BaseModel](t *testing.T, prefix string, created T, update func(T), emptyModel func() T) {
	fmt.Println("Testing update for prefix:" + " " + prefix + " with id: " + strconv.Itoa(created.GetID()) + "")
	update(created)
	reqBody, err := json.Marshal(created)

	if err != nil {
		panic(err)
	}
	res, err := http.NewRequest(http.MethodPut, testutil.TestApiUrl(prefix+"/"+strconv.Itoa(created.GetID())), bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	body := emptyModel()
	err = json.Unmarshal(buf, body)
	if err != nil {
		fmt.Println("Failed to unmarshal: ", string(buf))
		panic(err)
	}
	created.SetUpdatedAt(body.GetUpdatedAt())
	assert.Equal(t, created, body, "Body should be equal to updated model")
}

func testDelete[T models.BaseModel](t *testing.T, prefix string, created T) {
	fmt.Println("Testing delete for prefix:" + " " + prefix + " with id: " + strconv.Itoa(created.GetID()) + "")
	req, err := http.NewRequest(http.MethodDelete, testutil.TestApiUrl(prefix+"/"+strconv.Itoa(created.GetID())), nil)

	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusNoContent, res.StatusCode, "Expected status code 204")

	res, err = http.Get(testutil.TestApiUrl(prefix + "/" + strconv.Itoa(created.GetID())))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Expected status code 404")
}

func testReadAll[T models.BaseModel](t *testing.T, prefix string, createMany func(int)) {
	fmt.Println("Testing read all for prefix:" + " " + prefix + "")
	createMany(10)
	// Test pagination with size and page query parameters
	var size, page int
	size = 5
	page = 1

	res, err := http.Get(testutil.TestApiUrl(prefix + "?size=" + strconv.Itoa(size) + "&page=" + strconv.Itoa(page)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected status code 200")

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var body []T
	err = json.Unmarshal(buf, &body)
	if err != nil {
		fmt.Println("Failed to unmarshal: ", string(buf))
		panic(err)
	}

	// Check if the number of items retrieved is equal to size
	assert.Equal(t, size, len(body), "The number of items retrieved should be equal to the size parameter")

	// Additional checks:
	for _, item := range body {
		assert.NotZero(t, item.GetID(), "ID should not be zero")
		assert.NotZero(t, item.GetCreatedBy(), "CreatedBy should not be zero")
		assert.NotNil(t, item.GetCreatedAt(), "CreatedAt should not be nil")
		assert.NotNil(t, item.GetUpdatedAt(), "UpdatedAt should not be nil")
	}

	// Test another page
	page = 2
	res, err = http.Get(testutil.TestApiUrl(prefix + "?size=" + strconv.Itoa(size) + "&page=" + strconv.Itoa(page)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected status code 200")

	buf, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(buf, &body)
	if err != nil {
		panic(err)
	}

	// Check if the number of items retrieved is equal to size
	assert.Equal(t, size, len(body), "The number of items retrieved should be equal to the size parameter")

	// Additional checks:
	for _, item := range body {
		assert.NotZero(t, item.GetID(), "ID should not be zero")
		assert.NotZero(t, item.GetCreatedBy(), "CreatedBy should not be zero")
		assert.NotNil(t, item.GetCreatedAt(), "CreatedAt should not be nil")
		assert.NotNil(t, item.GetUpdatedAt(), "UpdatedAt should not be nil")
	}
}
