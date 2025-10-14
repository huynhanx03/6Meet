package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ResponseWithPagination struct {
	Code       int                 `json:"code"`
	Message    string              `json:"message"`
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

// ToErrorResponse converts error to error response
func ToErrorResponse(input interface{}) string {
    var messages []string

    switch v := input.(type) {
    case string:
        messages = []string{v}
    case []string:
        messages = v
    case error:
        messages = []string{v.Error()}
    default:
        messages = []string{"Unknown error"}
    }

    return strings.Join(messages, ", ")
}


// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(getHTTPStatusCode(code), ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}