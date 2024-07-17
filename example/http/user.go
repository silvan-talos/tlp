package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/silvan-talos/tlp/example/user"
	"github.com/silvan-talos/tlp/log"
	"github.com/silvan-talos/tlp/logging"
	"github.com/silvan-talos/tlp/transaction"
)

type userHandler struct {
	us *user.Service
}

func (uh *userHandler) addRoutes(r *gin.RouterGroup) {
	r.GET("/:id", uh.getUser)
	r.POST("/", uh.createUser)
	r.PATCH("/:id", uh.updateUser)
	r.DELETE("/:id", uh.deleteUser)
}

func (uh *userHandler) getUser(c *gin.Context) {
	tx, ctx := transaction.DefaultTracer().StartTransaction(c.Request.Context(),
		"get user",
		"request",
		logging.NewAttr("requestPath", c.Request.URL.Path),
		logging.NewAttr("requestMethod", c.Request.Method),
		logging.NewAttr("callerIP", c.Request.RemoteAddr),
		logging.NewAttr("userAgent", c.Request.UserAgent()),
	)
	defer tx.End()

	stringID := c.Param("id")
	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong user id"})
		return
	}

	u, err := uh.us.GetUser(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func (uh *userHandler) createUser(c *gin.Context) {
	tx, ctx := transaction.DefaultTracer().StartTransaction(c.Request.Context(),
		"create user",
		"request",
		logging.NewAttr("requestPath", c.Request.URL.Path),
		logging.NewAttr("requestMethod", c.Request.Method),
		logging.NewAttr("callerIP", c.Request.RemoteAddr),
		logging.NewAttr("userAgent", c.Request.UserAgent()),
	)
	defer tx.End()

	var req UserCreation
	if err := c.ShouldBind(&req); err != nil {
		log.Warn(ctx, "create user: bind failed", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	u := user.User{
		Name: req.Name,
		Age:  req.Age,
	}
	id, err := uh.us.CreateUser(ctx, u)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type UserCreation struct {
	Name string `binding:"required,lte=100"`
	Age  int    `binding:"required,gt=0"`
}

func (uh *userHandler) updateUser(c *gin.Context) {
	tx, ctx := transaction.DefaultTracer().StartTransaction(c.Request.Context(),
		"update user",
		"request",
		logging.NewAttr("requestPath", c.Request.URL.Path),
		logging.NewAttr("requestMethod", c.Request.Method),
		logging.NewAttr("callerIP", c.Request.RemoteAddr),
		logging.NewAttr("userAgent", c.Request.UserAgent()),
	)
	defer tx.End()

	stringID := c.Param("id")
	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong user id"})
		return
	}
	var req UserUpdate
	if err := c.ShouldBind(&req); err != nil {
		log.Warn(ctx, "update user: bind failed", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	u := user.User{
		Name: req.Name,
		Age:  req.Age,
	}
	err = uh.us.UpdateUser(ctx, id, u)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

type UserUpdate struct {
	Name string `binding:"omitempty,lte=100"`
	Age  int    `binding:"omitempty,gt=0"`
}

func (uh *userHandler) deleteUser(c *gin.Context) {
	tx, ctx := transaction.DefaultTracer().StartTransaction(c.Request.Context(),
		"delete user",
		"request",
		logging.NewAttr("requestPath", c.Request.URL.Path),
		logging.NewAttr("requestMethod", c.Request.Method),
		logging.NewAttr("callerIP", c.Request.RemoteAddr),
		logging.NewAttr("userAgent", c.Request.UserAgent()),
	)
	defer tx.End()

	stringID := c.Param("id")
	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong user id"})
		return
	}

	err = uh.us.DeleteUser(ctx, id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
