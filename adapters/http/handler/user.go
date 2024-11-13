package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/service"
)

type (
	UserHandler struct {
		UserService service.UserService
		PostService service.PostService
	}
	NewUserReq struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
)

func NewUserHandler(userService service.UserService, postService service.PostService) *UserHandler {
	return &UserHandler{
		userService,
		postService,
	}
}

func (u *UserHandler) New(ctx *gin.Context) {

	var body NewUserReq

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	user := &models.User{UUID: "", Username: body.Username, Email: body.Email}
	user, err := u.UserService.New(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"message": "User Created",
	})

}

func (u *UserHandler) User(ctx *gin.Context) {
	iuser, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user", "message": "missing user"})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    iuser.(*models.User),
		"mesasge": "Users feed",
	})
}

func (u *UserHandler) Follow(ctx *gin.Context) {
	toFollow := ctx.Param("user_uuid")

	user, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user", "message": "missing user"})
		ctx.Abort()
		return
	}

	updatedUser, err := u.UserService.Follow(user.(*models.User), toFollow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow user", "message": "no user to follow found"})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedUser, "message": "user followed"})

}

func (u *UserHandler) Feed(ctx *gin.Context) {

	user, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user", "message": "missing user"})
		ctx.Abort()
		return
	}
	feed, err := u.UserService.Feed(user.(*models.User))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Feed not found for user",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	castedUser := user.(*models.User)
	castedUser.Feed = feed

	ctx.JSON(http.StatusOK, gin.H{
		"data":    castedUser,
		"mesasge": "Users feed",
	})
}
