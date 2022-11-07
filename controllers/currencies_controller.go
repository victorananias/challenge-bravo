package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/victorananias/challenge-bravo/contracts/requests"
	"github.com/victorananias/challenge-bravo/contracts/responses"
	"github.com/victorananias/challenge-bravo/models"
	"github.com/victorananias/challenge-bravo/repositories"
	"github.com/victorananias/challenge-bravo/services"
)

func CreateCurrencyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var currencyRequest requests.CreateCurrencyRequest

	err := json.NewDecoder(request.Body).Decode(&currencyRequest)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	exchangesRepository := repositories.NewCurrenciesRepository()

	currency := models.Currency{
		Code:                currencyRequest.Code,
		Value:               currencyRequest.Value,
		BackingCurrencyCode: currencyRequest.BackingCurrencyCode,
	}
	err = exchangesRepository.CreateOrUpdate(currency)
	if err != nil {
		respondWithJson(responseWriter, responses.NewMessageResponse(err.Error()), http.StatusBadRequest)
		return
	}

	respondWithJson(responseWriter, responses.NewMessageResponse("currency created"), http.StatusCreated)
}

func DeleteCurrencyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	code := strings.TrimPrefix(request.URL.Path, "/")
	exchangesRepository := repositories.NewCurrenciesRepository()

	err := exchangesRepository.DeleteByCurrencyCode(code)
	if err != nil {
		respondWithJson(responseWriter, responses.NewMessageResponse(err.Error()), http.StatusBadRequest)
		return
	}

	respondWithJson(responseWriter, responses.NewMessageResponse("currency deleted"), http.StatusOK)
}

func ConversionsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	sourceCurrency := request.URL.Query().Get("from")
	targetCurrency := request.URL.Query().Get("to")
	amount, err := strconv.ParseFloat(request.URL.Query().Get("amount"), 32)

	if err != nil {
		respondWithJson(responseWriter, responses.NewMessageResponse("error parsing amount"), http.StatusBadRequest)
		return
	}

	currenciesService := services.NewCurrenciesService()

	result, err := currenciesService.ConvertCurrencies(amount, sourceCurrency, targetCurrency)
	if err != nil {
		respondWithJson(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJson(responseWriter, responses.NewConversionsResponse(result, targetCurrency), http.StatusOK)
}

func respondWithJson(responseWriter http.ResponseWriter, i interface{}, status int) {
	jsonResponse, _ := json.Marshal(i)
	responseWriter.WriteHeader(status)
	_, err := responseWriter.Write(jsonResponse)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
	}
}