package controllers

import (
	"net/http"
	"steam-backend/models"
	"steam-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{userService: service}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var reg models.JoinRequestDto
	if err := c.ShouldBindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "param error"))
		return
	}

	newUser, err := ctrl.userService.Register(&reg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "register failed"))
		return
	}

	userDTO := models.UserDto{
		UserID:   newUser.UserID,
		UserName: newUser.UserName,
		NickName: newUser.NickName,
		Avatar:   newUser.Avatar,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(userDTO))
}

func (ctrl *UserController) Login(c *gin.Context) {
	var loginDTO models.LoginRequestDto

	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "param error"))
		return
	}

	token, user, err := ctrl.userService.Login(&loginDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, err.Error()))
	}

	resposne := models.LoginResponseDto{
		Token:  token,
		UserID: user.UserID,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(resposne))
}

func (ctrl *UserController) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	user, err := ctrl.userService.GetUserInfo(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, models.NotFoundResponse(nil, "user not exists"))
		return
	}

	userDTO := models.UserDto{
		UserID:   user.UserID,
		UserName: user.UserName,
		NickName: user.NickName,
		Avatar:   user.Avatar,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(userDTO))
}

func (ctrl *UserController) CheckUsernameAvailable(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "username cant be nil"))
		return
	}

	available, err := ctrl.userService.ChechUserNameAvailable(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "check error"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(available))
}

func (ctrl *UserController) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "搜索关键词不能为空"))
		return
	}

	users, err := ctrl.userService.SearchUsers(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "服务器错误"))
		return
	}

	userDTOs := make([]models.UserDto, len(users))
	for i, user := range users {
		userDTOs[i] = models.UserDto{
			UserID:   user.UserID,
			UserName: user.UserName,
			Avatar:   user.Avatar,
			NickName: user.NickName,
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse(userDTOs))
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "无效的用户ID"))
		return
	}

	user, err := ctrl.userService.GetUserInfo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NotFoundResponse(nil, "用户不存在"))
		return
	}

	userDTO := models.UserDto{
		UserID:   user.UserID,
		UserName: user.UserName,
		Avatar:   user.Avatar,
		NickName: user.NickName,
	}

	c.JSON(http.StatusOK, models.SuccessResponse(userDTO))
}
