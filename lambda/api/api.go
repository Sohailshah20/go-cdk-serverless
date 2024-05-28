package api

import (
	"encoding/json"
	"lambda/database"
	"lambda/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.ResgisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request - empty fields",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	res, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server rrror",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if res {
		return events.APIGatewayProxyResponse{
			Body:       "user already exists",
			StatusCode: http.StatusConflict,
		}, err
	}
	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server rrror",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       "user saved successfully",
		StatusCode: http.StatusOK,
	}, err
}

func (api ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginUser types.ResgisterUser
	err := json.Unmarshal([]byte(request.Body), &loginUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	user, err := api.dbStore.GetUser(loginUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if !types.ValidatePassword(user.PasswordHash, loginUser.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "invalid user credentials",
			StatusCode: http.StatusBadRequest,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       "successfully logged in ",
		StatusCode: http.StatusOK,
	}, nil
}
