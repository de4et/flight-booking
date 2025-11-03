package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/de4et/flight-booking/internal/logger"
	"github.com/de4et/flight-booking/internal/service"

	"github.com/gin-gonic/gin"
)

var ErrNoToken = errors.New("no token provided")

const tokenName = "token"

type SearchResultHandler struct {
	searchService *service.MultipleSearchService
}

func NewSearchResultHandler(searchService *service.MultipleSearchService) *SearchResultHandler {
	return &SearchResultHandler{
		searchService: searchService,
	}
}

func (handler *SearchResultHandler) Handle(c *gin.Context) {
	token := c.Query(tokenName)
	if len(token) == 0 {
		c.AbortWithError(http.StatusBadRequest, ErrNoToken)
		return
	}

	ctx := logger.WithContext(c, "token", token)
	ts, err := handler.searchService.SearchByToken(ctx, token)
	if err != nil {
		return
	}

	ctx = logger.WithContext(ctx, "trips", ts)
	slog.InfoContext(ctx, "Successfully recieved trips")

	c.JSON(http.StatusOK, ts.ToArray())
}
