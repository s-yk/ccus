package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/s-yk/ccus/db"
)

func Records(ctx echo.Context) error {
	filter := make(map[string]string)
	p := ctx.QueryParam("p")
	if p != "" {
		filter["Project"] = p
	}
	n := ctx.QueryParam("n")
	if n != "" {
		filter["Nampespase"] = n
	}
	c := ctx.QueryParam("c")
	if c != "" {
		filter["Class"] = c
	}
	t := ctx.QueryParam("t")
	if t != "" {
		filter["Type"] = t
	}

	ts, err := db.GetData(&filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, ts)
}
