package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	API struct {
		apiHandlers
		treeNode *Tree
		echo     *echo.Echo
	}

	apiHandlers interface {
		Start()

		search(c echo.Context) error
		insert(c echo.Context) error
		delete(c echo.Context) error
	}

	Response struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

func NewAPI(treeValues []int) (*API, error) {
	api := &API{}
	tree := &Tree{}
	api.treeNode = tree

	for i := 0; i < len(treeValues); i++ {
		err := tree.Insert(treeValues[i])
		if err != nil {
			return nil, err
		}
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("/api/v1")
	g.GET("/search", api.search)
	g.POST("/insert", api.insert)
	g.DELETE("/delete", api.delete)

	api.echo = e
	return api, nil
}

func (api *API) Start() {
	api.echo.Logger.Fatal(api.echo.Start(":1323"))
}

func (api *API) search(c echo.Context) error {
	value, err := strconv.Atoi(c.QueryParam("val"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	isFound := api.treeNode.Search(value)
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "OK",
		Data:    isFound,
	})
}

func (api *API) insert(c echo.Context) error {
	value := make(map[string]int)
	err := c.Bind(&value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = api.treeNode.Insert(value["val"])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "OK",
		Data:    nil,
	})
}

func (api *API) delete(c echo.Context) error {
	value, err := strconv.Atoi(c.QueryParam("val"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	err = api.treeNode.Delete(value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
