package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	StatusCode int         `json:"statuscode"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message"`
	Success    bool        `json:"Success"`
}

func JSONResponse(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	res := APIResponse{
		StatusCode: statusCode,
		Data:       data,
		Message:    message,
		Success:    statusCode >= 200 && statusCode < 400,
	}
	return c.Status(statusCode).JSON(res)
}
