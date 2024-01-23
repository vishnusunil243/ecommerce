package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

func TestUserSignup(t *testing.T) {
	tests := []struct {
		name           string
		input          helperStruct.UserReq
		expectedOutput response.UserData
		buildStub      func(mock sqlmock.Sqlmock)
		expectedErr    error
	}{
		{
			name: "successful creations",
			input: helperStruct.UserReq{
				Name:     "vishnu",
				Email:    "vishnusunil243@gmail.com",
				Mobile:   "8129987917",
				Password: "1234",
			},
			expectedOutput: response.UserData{
				Id:     1,
				Name:   "vishnu",
				Email:  "vishnusunil243@gmail.com",
				Mobile: "8129987917",
			},
			buildStub: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "mobile"}).AddRow(1, "vishnu", "vishnusunil243@gmail.com", "8129987917")
				mock.ExpectQuery("^INSERT INTO users (.+)$").WithArgs("vishnu", "vishnusunil243@gmail.com", "8129987917", "1234").WillReturnRows(rows)
			},
			expectedErr: nil,
		},
		{
			name: "duplicate user",
			input: helperStruct.UserReq{
				Name:     "vishnu",
				Email:    "vishnusunil243@gmail.com",
				Mobile:   "8129987917",
				Password: "1234",
			},
			expectedOutput: response.UserData{},
			buildStub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^INSERT INTO users (.+)$").WithArgs("vishnu", "vishnusunil243@gmail.com", "8129987917", "1234").WillReturnError(errors.New("email or phone number already used"))
			},
			expectedErr: errors.New("email or phone number already used"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error %s was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			if err != nil {
				t.Fatalf("an error %s was not expected when initializing a mock db session", err)
			}
			userRepository := NewUserRepo(gormDB)
			tt.buildStub(mock)
			actualOutput, actualErr := userRepository.UserSignUp(tt.input)
			if tt.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.Equal(t, tt.expectedErr, actualErr)
			}
			if !reflect.DeepEqual(tt.expectedOutput, actualOutput) {
				t.Errorf("got %v, but want %v", actualOutput, tt.expectedOutput)
			}
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations %s", err)
			}
		})
	}
}
