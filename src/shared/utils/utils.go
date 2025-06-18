package utils

import (
	"encoding/json"
	"frcofilippi/pedimeapp/shared/logger"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type ApiResponse struct {
	Message   string `json:"message"`
	IsError   bool   `json:"isError"`
	RequestId string `json:"request_id"`
	Data      any    `json:"data,omitempty"`
}

func SendErrorResponse(w http.ResponseWriter, r *http.Request, httpCode int, errorMessage string, errToLog error) {

	requestId, _ := r.Context().Value(middleware.RequestIDKey).(string)

	w.Header().Add("Content-type", "application/json")
	w.Header().Add("x-request-id", requestId)
	w.WriteHeader(httpCode)

	errorResponse := &ApiResponse{
		Message:   errorMessage,
		IsError:   true,
		RequestId: requestId,
	}

	marshalledJsonResponse, err := json.MarshalIndent(&errorResponse, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshalledJsonResponse)

	logHanlderError(errorMessage, requestId, errToLog, errorResponse, r.URL.String())
}

func logHanlderError(errorDesc string, rid string, err error, response *ApiResponse, path string) {
	logger.GetLogger().Error(errorDesc, zap.String("request-id", rid), zap.String("path", path), zap.Error(err), zap.Any("response", response))
}

func logHandlerSuccess(rid string, path string, response *ApiResponse) {
	logger.GetLogger().Info("Request processed successfully", zap.String("request-id", rid), zap.String("path", path), zap.Any("response", response))
}

func SendApiResponse(w http.ResponseWriter, r *http.Request, httpCode int, message string, data any) {
	requestId, _ := r.Context().Value(middleware.RequestIDKey).(string)

	w.Header().Add("Content-type", "application/json")
	w.Header().Add("x-request-id", requestId)
	w.WriteHeader(httpCode)

	apiResponse := &ApiResponse{
		Message:   message,
		IsError:   false,
		RequestId: requestId,
		Data:      data,
	}

	marshalledJsonResponse, err := json.MarshalIndent(&apiResponse, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshalledJsonResponse)
	logHandlerSuccess(requestId, r.URL.String(), apiResponse)
}
