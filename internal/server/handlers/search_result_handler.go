package handlers

import (
	"fmt"

	"flight-booking/internal/service"

	"github.com/gin-gonic/gin"
)

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
		return
	}

	ts, err := handler.searchService.SearchByToken(c, token)
	if err != nil {
		return
	}

	fmt.Println(ts)
}
