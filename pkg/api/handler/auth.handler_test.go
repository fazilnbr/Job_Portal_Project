package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fazilnbr/project-workey/pkg/config"
	"github.com/fazilnbr/project-workey/pkg/repository"
	"github.com/fazilnbr/project-workey/pkg/response"
	"github.com/fazilnbr/project-workey/pkg/usecase"
	"github.com/fazilnbr/project-workey/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

// var Login = domain.Login{
// 	UserName: utils.RandomString(6),
// 	Password: fmt.Sprint(utils.RandomInt(10000, 99999)),
// }

var (
	Login           = []byte(`{UserName: "anu", Password: "12345"}`)
	gormDB, _       = utils.MockGormDB()
	authRepoMock    = repository.NewUserRepo(gormDB)
	authService     = usecase.NewUserService(authRepoMock)
	authServiceMock = NewAuthHandler(nil, nil, authService, nil, nil, config.Config{})
)

func TestLogin(t *testing.T) {

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(Login))

	gin := gin.New()
	res := httptest.NewRecorder()

	gin.POST("/signup", authServiceMock.UserSignUp)

	gin.ServeHTTP(res, req)

	status := res.Code

	if status != http.StatusOK {
		t.Errorf("Handler return wrong status code got : %v want : %v ", res.Code, http.StatusOK)
	}

	var newUser response.Response

	json.NewDecoder(io.Reader(res.Body)).Decode(&newUser)

	assert.NotNil(t, newUser)
	assert.Equal(t, Login, newUser.Data)

}

func TestLoginte(t *testing.T) {

	gormDB, _ := utils.MockGormDB()
	authRepoMock := repository.NewUserRepo(gormDB)
	authService := usecase.NewUserService(authRepoMock)
	authServiceMock := NewAuthHandler(nil, nil, authService, nil, nil, config.Config{})

	t.Run("test normal case login 1", func(t *testing.T) {

		// authServiceMock.On("Login", mock.AnythingOfType("string")).Return(nil)

		gin := gin.New()
		rec := httptest.NewRecorder()

		// authHandler := v1.NewAuthHandler(authServiceMock, infra.New("../../config/config.json"))
		gin.POST("/signup", authServiceMock.UserSignUp)

		body, err := json.Marshal(Login)
		// bodybite, err := json.MarshalIndent(domain.Login{UserName: "anu", Password: "12345"}, "", " ")
		assert.NoError(t, err)
		// bodystring := string(bodybite)
		// body := fmt.Sprint(bodystring)
		// assert.Equal(t, Login, body)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(body)))
		gin.ServeHTTP(rec, req)

		var newUser response.Response
		err = json.Unmarshal(rec.Body.Bytes(), &newUser)
		assert.NoError(t, err)

		exp := response.Response{
			Status:  false,
			Message: "Failed to create user",
			Errors:  []interface{}{"Username already exists"},
			Data:    nil,
		}

		t.Run("test fail response", func(t *testing.T) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
			assert.Equal(t, exp, newUser)
		})

		// exp := string(time.Now().Add(time.Hour * 2).Format(time.RFC3339))
		exp = response.Response{
			Status:  true,
			Message: "SUCCESS",
			Errors:  "",
			Data:    Login,
		}

		t.Run("test success response", func(t *testing.T) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, exp, newUser)
		})

	})
}
