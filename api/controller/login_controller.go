package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

/*
login and signup flow will be changed
there is no need for email password pair
everything happens based on the google jwt
google JWT validation:
  - if ok -> generate access/refresh token pair
  - if first access -> signup flow
  - if returning access -> normal login flow
  - if not ok -> error message
*/
func (lc *LoginController) Login(c *gin.Context) {
	var request domain.GoogleLoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// -> apply for non google flow - should it be kept?
	// user, err := lc.LoginUsecase.GetUserByEmail(c, request.Email)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given email"})
	// 	return
	// }

	// if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
	// 	c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid credentials"})
	// 	return
	// }

	claims, err := lc.LoginUsecase.ExtractGoogleClaims(request.Credential)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := lc.LoginUsecase.GetUserByEmail(c, claims.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given email"})
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(&user, lc.Env.RefreshTokenSecret, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	cookieValue := fmt.Sprintf("%s---%s", accessToken, refreshToken)
	exp := time.Now().Add(time.Hour * time.Duration(lc.Env.RefreshTokenExpiryHour)).Unix()
	c.SetCookie(domain.AuthToken, cookieValue, int(exp), "/", "localhost", true, true)

	successResponse := domain.SuccessResponse{
		Message: "http only cookie set",
	}

	c.JSON(http.StatusOK, successResponse)
}
