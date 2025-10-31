package controllers

import (
	"net/http"
	"steam-backend/models"
	"steam-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendController struct {
	friendService services.FriendService
}

func NewFriendController(server services.FriendService) *FriendController {
	return &FriendController{friendService: server}
}

func (ctrl *FriendController) GetFriendCount(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	count, err := ctrl.friendService.GetFriendCount(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get friendCount failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(count))
}

func (ctrl *FriendController) GetFriendList(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	list, err := ctrl.friendService.GetFriendList(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get friendList failed"))
		return
	}

	listDto := make([]models.UserDto, len(list))
	for i, user := range list {
		listDto[i] = models.UserDto{
			UserID:   user.UserID,
			UserName: user.UserName,
			Avatar:   user.Avatar,
			NickName: user.NickName,
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse(listDto))
}

func (ctrl *FriendController) SendInvitation(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	var req models.InvitationRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "request param error"))
		return
	}

	err := ctrl.friendService.SendInvitation(userID.(uint64), req.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "send invitaion failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "send successful"))
}

func (ctrl *FriendController) AcceptInvitation(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "invaild ID"))
		return
	}

	err = ctrl.friendService.AcceptInvitation(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "accept invitaion failed"))
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "accept invitaion"))
}

func (ctrl *FriendController) RefuseInvitation(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "invaild ID"))
		return
	}

	err = ctrl.friendService.RefuseInvitation(id, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "refuse invitation failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "refuse successful"))
}

func (ctrl *FriendController) GetReceivedInvitations(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	invitaions, err := ctrl.friendService.GetReceivedInvitations(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get received invitaion failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(invitaions))
}

func (ctrl *FriendController) GetSentInvitations(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.UnauthorizedResponse(nil, "unauthorized"))
		return
	}

	invitaions, err := ctrl.friendService.GetSentInvitations(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "get received invitaion failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(invitaions))
}

func (ctrl *FriendController) RemoveFriend(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.BadRequestResponse(nil, "unauthorized"))
		return
	}

	var req models.FriendRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "body error"))
		return
	}

	err := ctrl.friendService.RemoveFriend(userID.(uint64), req.FriendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil,
			"removeFriend failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMsg(nil, "remove successful"))
}

func (ctrl *FriendController) CheckFriendship(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.BadRequestResponse(nil, "unauthorized"))
		return
	}

	friendIDStr := c.Query("friendId")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BadRequestResponse(nil, "invaild friendId"))
		return
	}

	isFriend, err := ctrl.friendService.CheckFriendShip(userID.(uint64), friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerErrorResponse(nil, "chechFriendShip failed"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(isFriend))
}
