package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/service"
)

type (
	PostHandler struct {
		PostService service.PostService
		UserService service.UserService
	}

	NewPostReq struct {
		Content string `json:"content"`
	}
)

func NewPostHandler(userService service.UserService, postService service.PostService) *PostHandler {
	return &PostHandler{
		postService,
		userService,
	}
}

func (p *PostHandler) New(ctx *gin.Context) {
	var body NewPostReq

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	iuser, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user", "message": "missing user"})
		ctx.Abort()
		return
	}

	user := iuser.(*models.User)

	post, err := p.PostService.New(body.Content, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    post,
		"message": "Post Created",
	})

}

func (p *PostHandler) Post(ctx *gin.Context) {
	uuid := ctx.Param("post_uuid")

	if _, exist := ctx.Get("user"); !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user", "message": "missing user"})
		ctx.Abort()
		return
	}

	post, err := p.PostService.Post(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no post found",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number", "message": "Invalid page"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number", "message": "Invalid limit"})
		return
	}

	comments, err := p.PostService.Comments(post, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no post found",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	post.Comments = comments

	ctx.JSON(http.StatusOK, gin.H{
		"data":    post,
		"mesasge": "Post found",
	})
}

func (p *PostHandler) Comment(ctx *gin.Context) {
	uuid := ctx.Param("post_uuid")

	var body NewPostReq

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	iuser, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user", "message": "missing user"})
		ctx.Abort()
		return
	}

	user := iuser.(*models.User)

	post, err := p.PostService.Post(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	updatedPost, err := p.PostService.Comment(body.Content, post, user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    updatedPost,
		"message": "Post commented",
	})

}

func (p *PostHandler) Comments(ctx *gin.Context) {
	uuid := ctx.Param("post_uuid")

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number", "message": "Invalid page"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number", "message": "Invalid limit"})
		return
	}

	post, err := p.PostService.Post(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
			"data":  err.Error(),
		})
		ctx.Abort()
		return
	}

	comments, _ := p.PostService.Comments(post, pageInt, limitInt)
	post.Comments = comments

	ctx.JSON(http.StatusOK, gin.H{
		"data":    post,
		"mesasge": "Post commented",
	})

}
