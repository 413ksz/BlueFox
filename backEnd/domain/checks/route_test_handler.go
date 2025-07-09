package checks

import (
	"net/http"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
)

type TestResponse struct {
	Message string `json:"message"`
}

func TestHandler(w http.ResponseWriter, r *http.Request) (*models.ApiResponse[TestResponse], *models.CustomError) {
	apiReasponse := models.NewApiResponse[TestResponse](nil, http.StatusBadRequest, "Bad request")

	// Create an instance of our response struct
	response := TestResponse{
		Message: "Hello from the /api/test endpoint! Routing works!",
	}

	apiReasponse.WithData("routing works as expected", &models.ResponseData[TestResponse]{Items: []TestResponse{response}}, http.StatusOK)
	return apiReasponse, nil
}
