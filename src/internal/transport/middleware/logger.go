package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			log.Printf("%+v", err)

			if c.Response() == nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "Internal Server Error",
				})
			}

			return nil
		}

		log.Printf("%s %s %d", c.Request().Method, c.Request().RequestURI, c.Response().Status)

		return nil
	}
}
