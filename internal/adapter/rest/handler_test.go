package rest_test

import (
	"bytes"
	"computerInventory/internal/adapter/rest"
	"computerInventory/internal/domain"
	"computerInventory/internal/testhelpers"
	"computerInventory/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest() (*gin.Engine, *testhelpers.MockNotifier) {
	gin.SetMode(gin.TestMode)
	repo := testhelpers.NewMockRepo()
	n := &testhelpers.MockNotifier{}

	service := usecase.NewService(repo, n)

	h := rest.NewHandler(service)

	r := gin.Default()
	h.RegisterRoutes(r)
	return r, n
}

func TestAddComputerHandler(t *testing.T) {
	r, _ := setupTest()

	computer := domain.Computer{
		MACAddress:           "00:11:22:33:44:66",
		ComputerName:         "testcomp",
		IPAddress:            "192.168.1.99",
		EmployeeAbbreviation: "abc",
		Description:          "test machine",
	}

	jsonData, _ := json.Marshal(computer)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/computers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	getResp := httptest.NewRecorder()
	getReq, _ := http.NewRequest("GET", "/computers/00:11:22:33:44:66", nil)
	r.ServeHTTP(getResp, getReq)
	assert.Equal(t, http.StatusOK, getResp.Code)

	var returned domain.Computer
	err := json.Unmarshal(getResp.Body.Bytes(), &returned)
	assert.NoError(t, err)
	assert.Equal(t, "testcomp", returned.ComputerName)
}

func TestGetAllComputersHandler(t *testing.T) {
	r, _ := setupTest()

	computers := []domain.Computer{
		{
			MACAddress:           "00:11:22:33:44:77",
			ComputerName:         "alpha",
			IPAddress:            "192.168.0.1",
			EmployeeAbbreviation: "abc",
		},
		{
			MACAddress:           "00:11:22:33:44:88",
			ComputerName:         "beta",
			IPAddress:            "192.168.0.2",
			EmployeeAbbreviation: "xyz",
		},
	}

	for _, comp := range computers {
		body, _ := json.Marshal(comp)
		req, _ := http.NewRequest("POST", "/computers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
	}

	req, _ := http.NewRequest("GET", "/computers", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result []domain.Computer
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetByEmployeeHandler(t *testing.T) {
	r, _ := setupTest()

	computer := domain.Computer{
		MACAddress:           "00:11:22:33:44:77",
		ComputerName:         "empMachine",
		IPAddress:            "192.168.1.20",
		EmployeeAbbreviation: "jdo",
		Description:          "assigned to employee jdoe",
	}

	jsonData, _ := json.Marshal(computer)
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/computers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	getResp := httptest.NewRecorder()
	getReq, _ := http.NewRequest("GET", "/employee/jdo/computers", nil)
	r.ServeHTTP(getResp, getReq)

	assert.Equal(t, http.StatusOK, getResp.Code)

	var returned []domain.Computer
	err := json.Unmarshal(getResp.Body.Bytes(), &returned)
	assert.NoError(t, err)
	assert.Len(t, returned, 1)
	assert.Equal(t, "empMachine", returned[0].ComputerName)
}

func TestNotificationFromHandler(t *testing.T) {

	r, n := setupTest()
	n.On("SendWarning", "mmu", mock.Anything).Once()

	// Create 3 computers for the same employee "mmu"
	for i := 1; i <= 3; i++ {
		mac := fmt.Sprintf("00:11:22:33:44:%02d", i)
		comp := domain.Computer{
			MACAddress:           mac,
			ComputerName:         fmt.Sprintf("comp%d", i),
			IPAddress:            fmt.Sprintf("192.168.1.%d", i),
			EmployeeAbbreviation: "mmu",
			Description:          "test",
		}
		body, _ := json.Marshal(comp)

		req := httptest.NewRequest("POST", "/computers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}
	n.AssertExpectations(t)
}

func TestNoNotificationFromHandler(t *testing.T) {

	r, n := setupTest()

	// Create 2 computers for the same employee "mmu"
	for i := 1; i <= 2; i++ {
		mac := fmt.Sprintf("00:11:22:33:44:%02d", i)
		comp := domain.Computer{
			MACAddress:           mac,
			ComputerName:         fmt.Sprintf("comp%d", i),
			IPAddress:            fmt.Sprintf("192.168.1.%d", i),
			EmployeeAbbreviation: "mmu",
			Description:          "test",
		}
		body, _ := json.Marshal(comp)

		req := httptest.NewRequest("POST", "/computers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}
	n.AssertNotCalled(t, "SendWarning", mock.Anything, mock.Anything)
	n.AssertExpectations(t)
}
