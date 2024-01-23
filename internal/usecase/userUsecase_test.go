package usecase

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	mock_interfaces "main.go/internal/repository/mockRepository"
)

type eqCreateParamsMatcher struct {
	arg      helperStruct.UserReq
	password string
}

func (e eqCreateParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(helperStruct.UserReq)
	if !ok {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(arg.Password), []byte(e.password)); err != nil {
		return false
	}
	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}
func (e eqCreateParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}
func EqCreateParams(arg helperStruct.UserReq, password string) gomock.Matcher {
	return eqCreateParamsMatcher{arg, password}
}
func TestUserSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mock_interfaces.NewMockUserRepository(ctrl)
	UserUseCase := NewUserUsecase(userRepo)
	testData := []struct {
		name           string
		input          helperStruct.UserReq
		buildStub      func(userRepo mock_interfaces.MockUserRepository)
		expectedOutput response.UserData
		expectedError  error
	}{
		{
			name: "new user",
			input: helperStruct.UserReq{
				Name:     "vishnu",
				Email:    "vishnusunil243@gmail.com",
				Mobile:   "8129987917",
				Password: "1234",
			},
			buildStub: func(userRepo mock_interfaces.MockUserRepository) {
				userRepo.EXPECT().UserSignUp(
					EqCreateParams(helperStruct.UserReq{
						Name:     "vishnu",
						Email:    "vishnusunil243@gmail.com",
						Mobile:   "8129987917",
						Password: "1234",
					}, "1234")).Times(1).Return(response.UserData{
					Id:     1,
					Name:   "vishnu",
					Email:  "vishnusunil243@gmail.com",
					Mobile: "8129987917",
				}, nil)
			},
			expectedOutput: response.UserData{
				Id:     1,
				Name:   "vishnu",
				Email:  "vishnusunil243@gmail.com",
				Mobile: "8129987917",
			},
			expectedError: nil,
		},
		{
			name: "already exists",
			input: helperStruct.UserReq{
				Name:     "vishnu",
				Email:    "vishnusunil243@gmail.com",
				Mobile:   "8129987917",
				Password: "1234",
			},
			buildStub: func(userRepo mock_interfaces.MockUserRepository) {
				userRepo.EXPECT().UserSignUp(EqCreateParams(helperStruct.UserReq{
					Name:     "vishnu",
					Email:    "vishnusunil243@gmail.com",
					Mobile:   "8129987917",
					Password: "1234",
				}, "1234")).Times(1).Return(response.UserData{}, errors.New("user already exists"))
			},
			expectedOutput: response.UserData{},
			expectedError:  errors.New("user already exists"),
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userRepo)
			actualUser, err := UserUseCase.UserSignup(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, actualUser, tt.expectedOutput)
		})
	}
}
