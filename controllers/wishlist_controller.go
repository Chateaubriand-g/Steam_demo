package controllers

import (
	"net/http"
	"steam-backend/models"
	"steam-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistController struct {
	wishlistService services.WishlistService
}

func NewWishlistController(wishlistService services.WishlistService) *WishlistController {
	return &WishlistController{
		wishlistService: wishlistService,
	}
}

func (ctrl *WishlistController) GetWishlistSize(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "未授权访问"))
		return
	}

	size, err := ctrl.wishlistService.GetWishlistSize(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "获取愿望清单大小失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(size))
}

func (ctrl *WishlistController) GetWishlist(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "未授权访问"))
		return
	}

	wishlist, err := ctrl.wishlistService.GetWishlist(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "获取愿望清单失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(wishlist))
}

func (ctrl *WishlistController) AddToWishlist(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "未授权访问"))
		return
	}

	var req models.WishlistRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "请求参数错误"))
		return
	}

	err := ctrl.wishlistService.AddToWishlist(userID.(uint64), req.AppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "已添加到愿望清单"))
}

func (ctrl *WishlistController) RemoveFromWishlist(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "未授权访问"))
		return
	}

	appIDStr := c.Param("appId")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "无效的应用ID"))
		return
	}

	err = ctrl.wishlistService.RemoveFromWishlist(userID.(uint64), appID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "已从愿望清单移除"))
}

func (ctrl *WishlistController) IsInWishlist(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "未授权访问"))
		return
	}

	appIDStr := c.Query("appId")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "无效的应用ID"))
		return
	}

	isInWishlist, err := ctrl.wishlistService.IsInWishlist(userID.(uint64), appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "检查是否在愿望清单中失败"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(isInWishlist))
}

func (ctrl *WishlistController) SortWishlist(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "未授权访问"))
		return
	}

	var req models.SortWishlistRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "请求参数错误"))
		return
	}

	err := ctrl.wishlistService.SortWishlist(userID.(uint64), req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "愿望清单已排序"))
}
