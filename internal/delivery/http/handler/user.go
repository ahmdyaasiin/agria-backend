package handler

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/oauth"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/verification"
	usecaseInterface "github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type UserHandler struct {
	Log           *logrus.Logger
	Validator     *validator.Validate
	FacebookOAuth *oauth2.Config
	GoogleOAuth   *oauth2.Config
	UserUserCase  usecaseInterface.UserUseCase
}

func NewUserHandler(log *logrus.Logger, validator *validator.Validate, facebookOAuth *oauth2.Config, googleOAuth *oauth2.Config, userUC usecaseInterface.UserUseCase) interfaces.UserHandler {
	return &UserHandler{Log: log, Validator: validator, FacebookOAuth: facebookOAuth, GoogleOAuth: googleOAuth, UserUserCase: userUC}
}

func (h *UserHandler) URLOAuthFacebook(ctx fiber.Ctx) error {
	url := h.FacebookOAuth.AuthCodeURL(verification.GenerateRandomString(36))

	return ctx.Redirect().Status(fiber.StatusMovedPermanently).To(url)
}

func (h *UserHandler) URLOAuthGoogle(ctx fiber.Ctx) error {
	url := h.GoogleOAuth.AuthCodeURL(verification.GenerateRandomString(36))

	return ctx.Redirect().Status(fiber.StatusMovedPermanently).To(url)
}

func (h *UserHandler) FacebookOAuthCallback(ctx fiber.Ctx) error {

	resp := new(response.OAuth)
	var fetch *response.FetchFacebookProfile

	code := ctx.Query("code")
	token, err := h.FacebookOAuth.Exchange(ctx.Context(), code)
	if err != nil {
		h.Log.Warnf("Failed to exchange facebook code: %+v \n", err)

		resp.Error = true
		resp.ErrorMessage = ErrFRExchangeOAuthCode
		goto redirect
	}

	fetch, err = oauth.FetchFacebookProfile(token.AccessToken)
	if err != nil {
		h.Log.Warnf("Failed to fetch facebook profile: %+v \n", err)

		resp.Error = true
		resp.ErrorMessage = ErrFRFetchOAuthProfile
		goto redirect
	}

	resp, err = h.UserUserCase.FacebookCallBack(ctx.Context(), fetch)
	if err != nil {
		resp.Error = true
		resp.ErrorMessage = err.Error()
		goto redirect
	}

	if resp.RefreshToken != "" {
		ctx.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			Path:     "/",
			Expires:  time.Now().Local().Add(7 * 24 * time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
		})
	}

redirect:
	url := oauth.DetermineRedirectURL(resp)
	return ctx.Redirect().Status(fiber.StatusMovedPermanently).To(url)
}

func (h *UserHandler) GoogleOAuthCallback(ctx fiber.Ctx) error {

	resp := new(response.OAuth)
	var fetch *response.FetchGoogleProfile

	code := ctx.Query("code")
	token, err := h.GoogleOAuth.Exchange(ctx.Context(), code)
	if err != nil {
		h.Log.Warnf("Failed to exchange facebook code: %+v \n", err)

		resp.Error = true
		resp.ErrorMessage = ErrFRExchangeOAuthCode
		goto redirect
	}

	fetch, err = oauth.FetchGoogleProfile(token.AccessToken)
	if err != nil {
		h.Log.Warnf("Failed to fetch facebook profile: %+v \n", err)

		resp.Error = true
		resp.ErrorMessage = ErrFRFetchOAuthProfile
		goto redirect
	}

	resp, err = h.UserUserCase.GoogleCallBack(ctx.Context(), fetch)
	if err != nil {
		resp.Error = true
		resp.ErrorMessage = err.Error()
		goto redirect
	}

	if resp.RefreshToken != "" {
		ctx.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			Path:     "/",
			Expires:  time.Now().Local().Add(7 * 24 * time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
		})
	}

redirect:
	url := oauth.DetermineRedirectURL(resp)
	return ctx.Redirect().Status(fiber.StatusMovedPermanently).To(url)
}

func (h *UserHandler) RegisterWithOAuth(ctx fiber.Ctx) error {

	req := new(request.FinishRegisterOAuth)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	res, err := h.UserUserCase.RegisterWithOAuth(ctx.Context(), req)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Local().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusCreated).JSON(response.Final{
		Message: "Account created successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusCreated,
			Message: http.StatusText(fiber.StatusCreated),
		},
	})
}

func (h *UserHandler) PreRegister(ctx fiber.Ctx) error {
	req := new(request.PreRegister)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "Email and username are available and not duplicated",
		Data:    nil,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *UserHandler) RegisterWithEmailPassword(ctx fiber.Ctx) error {
	req := new(request.Register)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	err := h.UserUserCase.RegisterWithEmailPassword(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "The code has been sent to your email.",
		Data:    nil,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *UserHandler) SendVerificationCodeForRegister(ctx fiber.Ctx) error {
	req := new(request.PostRegister)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	err := h.UserUserCase.SendVerificationCodeForRegister(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "The code has been sent to your email.",
		Data:    nil,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *UserHandler) RegisterComplete(ctx fiber.Ctx) error {
	req := new(request.FinishRegister)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	res, err := h.UserUserCase.VerifySixCode(ctx.Context(), req)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Local().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusCreated).JSON(response.Final{
		Message: "Account created successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusCreated,
			Message: http.StatusText(fiber.StatusCreated),
		},
	})
}

func (h *UserHandler) Login(ctx fiber.Ctx) error {
	req := new(request.Login)

	if err := ctx.Bind().JSON(req); err != nil {
		h.Log.Warnf("failed to bind request: %+v\n", err)
		return ErrBindRequest
	}

	if err := h.Validator.Struct(req); err != nil {
		h.Log.Warnf("failed to validate request: %+v\n", err)
		return err
	}

	res, err := h.UserUserCase.Login(ctx.Context(), req)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Local().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusCreated).JSON(response.Final{
		Message: "Account login successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}

func (h *UserHandler) RenewAccessToken(ctx fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	res, err := h.UserUserCase.RenewAccessToken(ctx.Context(), refreshToken)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.Final{
		Message: "Renew access token successfully",
		Data:    res,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusCreated,
			Message: http.StatusText(fiber.StatusCreated),
		},
	})
}

func (h *UserHandler) Logout(ctx fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	err := h.UserUserCase.Logout(ctx.Context(), refreshToken)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Local().Add(-1 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(response.Final{
		Message: "You have successfully logged out",
		Data:    nil,
		Errors:  nil,
		Status: response.Status{
			Code:    fiber.StatusOK,
			Message: http.StatusText(fiber.StatusOK),
		},
	})
}
