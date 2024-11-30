package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DSSD-Madison/gmu/models"
)

const MinQueryLengh = 3

func fetchSearchPage(c echo.Context) error {
	query := c.FormValue("query")
	if len(query) < MinQueryLengh {
		return c.String(http.StatusBadRequest, "Error 400, could not get query from request.")
	}
	return c.Render(http.StatusOK, "search-init", query)
}

func Search(c echo.Context) error {
	query := c.FormValue("query")
	// Attribute associated with initial load of search page, nonempty on first load
	init := c.FormValue("init")
	if len(query) < MinQueryLengh {
		return c.String(http.StatusBadRequest, "Error 400, could not get query from request.")
	}
	// If nonempty first load attribute, send OOB search bar with request
	if init != "" {
		c.Render(http.StatusOK, "search", query)
	}
	results := models.MakeQuery(query)

	return c.Render(http.StatusOK, "results", results)
}
