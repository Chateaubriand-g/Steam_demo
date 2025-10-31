package controllers

import (
	"errors"
	"net/http"
	"steam-backend/models"
	"steam-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppController struct {
	appService services.AppService
}

func NewAppController(appService services.AppService) *AppController {
	return &AppController{
		appService: appService,
	}
}

func (ctrl *AppController) GetRecommendations(c *gin.Context) {
	//c.DefaultQuery从url中查询参数，若未提供则使用默认值
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit >= 30 {
		limit = 30
	}

	recommendations, err := ctrl.appService.GetRecommendations(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get recommendations failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(recommendations))
}

func (ctrl *AppController) GetSpecials(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}
	if limit > 30 {
		limit = 30
	}

	specials, err := ctrl.appService.GetSpecials(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get specials fialed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(specials))
}

func (ctrl *AppController) GetSearchSuggestions(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "keyword is nil"))
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 5
	}
	if limit > 5 {
		limit = 5
	}

	suggestions, err := ctrl.appService.GetSearchSuggestions(keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get suggestions failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(suggestions))
}

func (ctrl *AppController) GetAppByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64) //十进制，64位
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "invaild ID"))
		return
	}

	app, err := ctrl.appService.GetAppByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.NotFoundResponse(nil, "id not exists"))
			return
		} else {
			c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "search failed"))
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse(app))
}
