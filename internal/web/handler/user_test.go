package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	helper "main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	mock_interfaces "main.go/internal/usecase/mockUsecase"
	"main.go/internal/web/middleware"
)

// test for userlogin
func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	userUseCase := mock_interfaces.NewMockUserUseCase(ctrl)
	cartUseCase := mock_interfaces.NewMockCartUseCase(ctrl)
	walletUseCase := mock_interfaces.NewMockWalletUseCase(ctrl)
	referraluseCase := mock_interfaces.NewMockReferralUseCase(ctrl)
	UserHandler := NewUserHandler(userUseCase, cartUseCase, walletUseCase, referraluseCase)

	testData := []struct {
		name             string
		loginData        helper.LoginReq
		buildStub        func(userUseCase mock_interfaces.MockUserUseCase)
		expectedCode     int
		expectedResponse response.Response
		expectedError    error
	}{
		{
			name: "valid login",
			loginData: helper.LoginReq{
				Email:    "vishnusunil243@gmail.com",
				Password: "1234",
			},
			buildStub: func(userUseCase mock_interfaces.MockUserUseCase) {
				userUseCase.EXPECT().UserLogin(helper.LoginReq{
					Email:    "vishnusunil243@gmail.com",
					Password: "1234",
				}).Times(1).Return("validToken", nil)
			},
			expectedCode: 200,
			expectedResponse: response.Response{
				StatusCode: 200,
				Message:    "login successfull",
				Data:       nil,
				Errors:     nil,
			},
			expectedError: nil,
		}, {
			name: "invalid login",
			loginData: helper.LoginReq{
				Email:    "invalid@example.com",
				Password: "invalid",
			},
			buildStub: func(userUseCase mock_interfaces.MockUserUseCase) {
				userUseCase.EXPECT().UserLogin(helper.LoginReq{
					Email:    "invalid@example.com",
					Password: "invalid",
				}).Times(1).Return("", errors.New("invalid credentials"))
			},
			expectedCode: 400,
			expectedResponse: response.Response{
				StatusCode: 400,
				Message:    "Login failed",
				Data:       nil,
				Errors:     "invalid credentials",
			},
			expectedError: errors.New("invalid credentials"),
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userUseCase)
			engine := gin.Default()
			recorder := httptest.NewRecorder()
			engine.POST("/user/login", UserHandler.UserLogin)
			var body []byte
			body, err := json.Marshal(tt.loginData)
			assert.NoError(t, err)
			url := "/user/login"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, req)
			var actual response.Response
			err = json.Unmarshal(recorder.Body.Bytes(), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedResponse.Message, actual.Message)
		})
	}
}
func TestViewUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)

	userUseCase := mock_interfaces.NewMockUserUseCase(ctrl)
	cartUseCase := mock_interfaces.NewMockCartUseCase(ctrl)
	walletUseCase := mock_interfaces.NewMockWalletUseCase(ctrl)
	referralUseCase := mock_interfaces.NewMockReferralUseCase(ctrl)
	UserHandler := NewUserHandler(userUseCase, cartUseCase, walletUseCase, referralUseCase)
	testData := []struct {
		name             string
		userID           int64
		buildStub        func(userUsecase mock_interfaces.MockUserUseCase)
		expectedCode     int
		expectedResponse response.Response
		expectedData     response.UserProfile
		expectedError    error
	}{
		{
			name:   "valid profile",
			userID: 1,
			buildStub: func(userUsecase mock_interfaces.MockUserUseCase) {
				userUsecase.EXPECT().ViewUserProfile(gomock.Any()).Times(1).Return(response.UserProfile{
					Name:       "TestUser",
					Email:      "test@gmail.com",
					Mobile:     "1234567890",
					ReferralId: "pc999999999999/1u",
					Address: response.Address{
						House_number: "12",
						Street:       "Mg",
						City:         "malappuram",
						District:     "malappuram",
						Landmark:     "KC",
						Pincode:      679322,
						IsDefault:    true,
					},
				}, nil)
			},
			expectedCode: 200,
			expectedResponse: response.Response{
				StatusCode: 200,
				Message:    "user profile retrieved successfully",
				Data: response.UserProfile{
					Name:       "TestUser",
					Email:      "test@gmail.com",
					Mobile:     "1234567890",
					ReferralId: "pc999999999999/1u",
					Address: response.Address{
						House_number: "12",
						Street:       "Mg",
						City:         "malappuram",
						District:     "malappuram",
						Landmark:     "KC",
						Pincode:      679322,
						IsDefault:    true,
					},
				},
				Errors: nil,
			},
			expectedData: response.UserProfile{
				Name:       "TestUser",
				Email:      "test@gmail.com",
				Mobile:     "1234567890",
				ReferralId: "pc999999999999/1u",
				Address: response.Address{
					House_number: "12",
					Street:       "Mg",
					City:         "malappuram",
					District:     "malappuram",
					Landmark:     "KC",
					Pincode:      679322,
					IsDefault:    true,
				},
			},
			expectedError: nil,
		},
		{
			name:   "invalid profile",
			userID: 2,
			buildStub: func(userUsecase mock_interfaces.MockUserUseCase) {
				userUsecase.EXPECT().ViewUserProfile(gomock.Any()).Times(1).Return(response.UserProfile{}, errors.New("user not found"))
			},
			expectedCode: 400,
			expectedResponse: response.Response{
				StatusCode: 400,
				Message:    "error retrieving user profile",
				Data:       nil,
				Errors:     "user not found",
			},
			expectedData:  response.UserProfile{},
			expectedError: errors.New("user not found"),
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userUseCase)
			engine := gin.Default()
			recorder := httptest.NewRecorder()
			engine.GET("/user/userprofile/", middleware.TestUserAuth, UserHandler.ViewUserProfile)
			url := "/user/userprofile/"
			req := httptest.NewRequest(http.MethodGet, url, nil)
			engine.ServeHTTP(recorder, req)
			var actual response.Response
			err := json.Unmarshal(recorder.Body.Bytes(), &actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, recorder.Code)
			assert.Equal(t, tt.expectedResponse.Message, actual.Message)
		})
	}
}
