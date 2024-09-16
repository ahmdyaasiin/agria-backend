package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/jwt"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/verification"
	repositoryInterface "github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserUseCase struct {
	DB                *sqlx.DB
	Log               *logrus.Logger
	Redis             *redis.Client
	UserRepository    repositoryInterface.UserRepository
	AddressRepository repositoryInterface.AddressRepository
	RefreshRepository repositoryInterface.RefreshRepository
}

func NewUserUseCase(DB *sqlx.DB, log *logrus.Logger, redis *redis.Client, userRepository repositoryInterface.UserRepository, addressRepository repositoryInterface.AddressRepository, refreshRepository repositoryInterface.RefreshRepository) interfaces.UserUseCase {
	return &UserUseCase{DB: DB, Log: log, Redis: redis, UserRepository: userRepository, AddressRepository: addressRepository, RefreshRepository: refreshRepository}
}

func (u *UserUseCase) FacebookCallBack(ctx context.Context, profile *response.FetchFacebookProfile) (*response.OAuth, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrFRCreateDatabaseTransaction
	}

	res := new(response.OAuth)
	user := &domain.User{
		Email: profile.Email,
	}

	err = u.UserRepository.Read(tx, "email", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get user data: %+v\n", err)
		return nil, ErrFRFailedToReadData
	}

	if user.ID != "" {
		if !user.IsFacebook {
			u.Log.Warnf("not registered with facebook: %s\n", user.Email)
			return nil, ErrFRNotFacebookUser
		}

		accessToken, err := jwt.CreateToken(user.Username, false)
		if err != nil {
			u.Log.Warnf("failed to create access token: %+v\n", err)
			return nil, ErrFRCreateToken
		}

		refresh := &domain.Refresh{
			UserID: user.ID,
		}
		err = u.RefreshRepository.Read(tx, "user_iD", refresh)
		if err != nil {
			u.Log.Warnf("failed to get refresh data: %+v\n", err)
			return nil, ErrFailedToReadData
		}

		refreshToken, err := jwt.CreateToken(user.Username, true)
		if err != nil {
			u.Log.Warnf("failed to create refresh token: %+v\n", err)
			return nil, ErrCreateToken
		}

		refresh.Token = refreshToken
		err = u.RefreshRepository.Update(tx, refresh)
		if err != nil {
			u.Log.Warnf("failed to store data (refresh): %+v\n", err)
			return nil, ErrFailedToStoreData
		}

		err = tx.Commit()
		if err != nil {
			u.Log.Warnf("failed to commit transaction: %+v\n", err)
			return nil, ErrFailedToStoreData
		}

		res.IsRegistered = true
		res.AccessToken = accessToken
		res.RefreshToken = refresh.Token
	} else {

		user.ID = uuid.NewString()
		user.Name = profile.Name
		user.UrlPhoto = profile.Picture.Data.URL
		user.Status = "identity-card-verification-needed"
		user.IsFacebook = true

		userMarshal, _ := json.Marshal(user)
		err = u.Redis.Set(ctx, fmt.Sprintf("%s_temp", user.Email), userMarshal, 5*time.Minute).Err()
		if err != nil {
			u.Log.Warnf("failed to store data to redis: %+v\n", err)
			return nil, ErrFRFailedToStoreData
		}

		res.Token = user.ID
	}

	return res, nil
}

func (u *UserUseCase) GoogleCallBack(ctx context.Context, profile *response.FetchGoogleProfile) (*response.OAuth, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrFRCreateDatabaseTransaction
	}

	res := new(response.OAuth)
	user := &domain.User{
		Email: profile.Email,
	}

	err = u.UserRepository.Read(tx, "email", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get user data: %+v\n", err)
		return nil, ErrFRFailedToReadData
	}

	if user.ID != "" {
		if !user.IsGoogle {
			u.Log.Warnf("not registered with google: %s\n", user.Email)
			return nil, ErrFRNotGoogleUser
		}

		accessToken, err := jwt.CreateToken(user.Username, false)
		if err != nil {
			u.Log.Warnf("failed to create access token: %+v\n", err)
			return nil, ErrFRCreateToken
		}

		refresh := &domain.Refresh{
			UserID: user.ID,
		}
		err = u.RefreshRepository.Read(tx, "user_iD", refresh)
		if err != nil {
			u.Log.Warnf("failed to get refresh data: %+v\n", err)
			return nil, ErrFailedToReadData
		}

		refreshToken, err := jwt.CreateToken(user.Username, true)
		if err != nil {
			u.Log.Warnf("failed to create refresh token: %+v\n", err)
			return nil, ErrCreateToken
		}

		refresh.Token = refreshToken
		err = u.RefreshRepository.Update(tx, refresh)
		if err != nil {
			u.Log.Warnf("failed to store data (refresh): %+v\n", err)
			return nil, ErrFailedToStoreData
		}

		err = tx.Commit()
		if err != nil {
			u.Log.Warnf("failed to commit transaction: %+v\n", err)
			return nil, ErrFailedToStoreData
		}

		res.IsRegistered = true
		res.AccessToken = accessToken
		res.RefreshToken = refresh.Token
	} else {

		user.ID = uuid.NewString()
		user.Name = profile.Name
		user.UrlPhoto = profile.Picture
		user.Status = "identity-card-verification-needed"
		user.IsGoogle = true

		userMarshal, _ := json.Marshal(user)
		err = u.Redis.Set(ctx, fmt.Sprintf("%s_temp", user.Email), userMarshal, 10*time.Minute).Err()
		if err != nil {
			u.Log.Warnf("failed to store data to redis: %+v\n", err)
			return nil, ErrFRFailedToStoreData
		}

		res.Token = user.ID
	}

	return res, nil
}

func (u *UserUseCase) RegisterWithOAuth(ctx context.Context, req *request.FinishRegisterOAuth) (*response.FinishRegister, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	userRedis := new(domain.User)

	userString, err := u.Redis.Get(ctx, fmt.Sprintf("%s_temp", req.Email)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		u.Log.Warnf("failed to get data from redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if errors.Is(err, redis.Nil) {
		u.Log.Warnf("key not found in redis: %s\n", req.Email)
		return nil, ErrInvalidToken
	}

	err = json.Unmarshal([]byte(userString), userRedis)
	if err != nil {
		u.Log.Warnf("failed to unmarshal: %+v\n", err)
		return nil, ErrFailedToUnMarshal
	}

	if userRedis.ID != req.Token {
		u.Log.Warnf("token doesn't match: %s\n", req.Email)
		return nil, ErrInvalidToken
	}

	phoneNumber := strings.Replace(req.PhoneNumber, "+", "", 1)
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = "62" + phoneNumber[1:]
	}

	user := &domain.User{
		Email:       req.Email,
		Username:    req.Username,
		PhoneNumber: phoneNumber,
	}

	err = u.UserRepository.CheckUserExists(tx, user)
	if err != nil {
		u.Log.Warnf("duplicate user: %+v\n", err)
		return nil, ErrDuplicateUser
	}

	user = userRedis

	err = u.Redis.Del(ctx, fmt.Sprintf("%s_temp", user.Email)).Err()
	if err != nil {
		u.Log.Warnf("failed to delete key redis: %+v\n", err)
		return nil, ErrFailedToClearData
	}

	now := time.Now().Local().UnixNano()

	user.ID = uuid.NewString()
	user.Username = req.Username
	user.PhoneNumber = phoneNumber
	user.CreatedAt = now
	user.UpdatedAt = now

	err = u.UserRepository.Create(tx, user)
	if err != nil {
		u.Log.Warnf("failed to store data (user): %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	var addressName string
	if strings.Contains(user.Name, " ") {
		addressName = strings.Split(user.Name, " ")[0]
	} else {
		addressName = user.Name
	}

	addressName += "'s Primary Address"
	address := &domain.Address{
		ID:          uuid.NewString(),
		Name:        addressName,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		IsPrimary:   true,
		PhoneNumber: phoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserID:      user.ID,
	}

	err = u.AddressRepository.Create(tx, address)
	if err != nil {
		u.Log.Warnf("failed to store data (address): %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	res := new(response.FinishRegister)

	jwtAccessToken, err := jwt.CreateToken(user.Username, false)
	if err != nil {
		u.Log.Warnf("failed to create access token: %+v\n", err)
		return nil, ErrFRCreateToken
	}

	jwtRefreshToken, err := jwt.CreateToken(user.Username, true)
	if err != nil {
		u.Log.Warnf("failed to create refresh token: %+v\n", err)
		return nil, ErrFRCreateToken
	}

	refresh := &domain.Refresh{
		ID:            uuid.NewString(),
		Token:         jwtRefreshToken,
		CreatedAt:     now,
		UpdatedAt:     now,
		LastRefreshAt: now,
		UserID:        user.ID,
	}

	err = u.RefreshRepository.Create(tx, refresh)
	if err != nil {
		u.Log.Warnf("failed to store data (refresh): %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	err = tx.Commit()
	if err != nil {
		u.Log.Warnf("failed to commit transaction: %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	res.AccessToken = jwtAccessToken
	res.RefreshToken = jwtRefreshToken

	return res, nil
}

func (u *UserUseCase) RegisterWithEmailPassword(ctx context.Context, req *request.Register) error {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return ErrCreateDatabaseTransaction
	}

	phoneNumber := strings.Replace(req.PhoneNumber, "+", "", 1)
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = "62" + phoneNumber[1:]
	}

	user := &domain.User{
		Email:       req.Email,
		Username:    req.Username,
		PhoneNumber: phoneNumber,
	}

	err = u.UserRepository.CheckUserExists(tx, user)
	if err != nil {
		u.Log.Warnf("duplicate user: %+v\n", err)
		return ErrDuplicateUser
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Warnf("failed to generate password: %+v\n", err)
		return ErrFailedToGeneratePassword
	}

	now := time.Now().Local().UnixNano()

	user.ID = uuid.NewString()
	user.Name = req.Name
	user.Password = string(password)
	user.Status = "email-verification-needed"
	user.UrlPhoto = "https://example.com/default-profile-picture.jpg"
	user.CreatedAt = now
	user.UpdatedAt = now

	err = u.UserRepository.Create(tx, user)
	if err != nil {
		u.Log.Warnf("failed to store data (user): %+v\n", err)
		return ErrFailedToStoreData
	}

	var addressName string
	if strings.Contains(user.Name, " ") {
		addressName = strings.Split(user.Name, " ")[0]
	} else {
		addressName = user.Name
	}

	addressName += "'s Primary Address"
	address := &domain.Address{
		ID:          uuid.NewString(),
		Name:        addressName,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		IsPrimary:   true,
		PhoneNumber: phoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserID:      user.ID,
	}

	err = u.AddressRepository.Create(tx, address)
	if err != nil {
		u.Log.Warnf("failed to store data (address): %+v\n", err)
		return ErrFailedToStoreData
	}

	sixCode := verification.GenerateVerificationCode()
	err = u.Redis.Set(ctx, fmt.Sprintf("%s_verification_code_register", user.Email), sixCode, 5*time.Minute).Err()
	if err != nil {
		u.Log.Warnf("failed to generate verification code: %+v\n", err)
		return ErrFailedToGenerateCode
	}

	err = verification.SendEmail(user.Email, "Verification Code", fmt.Sprintf("Verification code: %s", sixCode))
	if err != nil {
		u.Log.Warnf("failed to send verification code to email: %+v\n", err)
		return ErrFailedToSendEmail
	}

	err = tx.Commit()
	if err != nil {
		u.Log.Warnf("failed to commit transaction: %+v\n", err)
		return ErrFailedToStoreData
	}

	return nil
}

func (u *UserUseCase) SendVerificationCodeForRegister(ctx context.Context, req *request.PostRegister) error {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return ErrCreateDatabaseTransaction
	}

	user := &domain.User{
		Email: req.Email,
	}
	err = u.UserRepository.Read(tx, "email", user)
	if err != nil {
		u.Log.Warnf("failed to read data (user): %+v\n", err)
		return ErrFailedToReadData
	}

	if user.Status != "email-verification-needed" {
		u.Log.Warnf("failed to verify account (already verified): %s\n", user.Email)
		return ErrFailedToVerifyAccount
	}

	err = u.Redis.Get(ctx, fmt.Sprintf("%s_verification_code_register", user.Email)).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		u.Log.Warnf("failed to read data from redis: %+v\n", err)
		return ErrFailedToReadData
	}

	if err == nil {
		u.Log.Warnf("limit exceeded for send email: %s\n", user.Email)
		return ErrSendEmailLimitExceeded
	}

	sixCode := verification.GenerateVerificationCode()
	if err != nil {
		u.Log.Warnf("failed to generate verification code: %+v\n", err)
		return ErrFailedToGenerateCode
	}

	err = verification.SendEmail(user.Email, "Verification Code", fmt.Sprintf("Verification code: %s", sixCode))
	if err != nil {
		u.Log.Warnf("failed to send verification code to email: %+v\n", err)
		return ErrFailedToSendEmail
	}

	return nil
}

func (u *UserUseCase) VerifySixCode(ctx context.Context, req *request.FinishRegister) (*response.FinishRegister, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	user := &domain.User{
		Email: req.Email,
	}
	err = u.UserRepository.Read(tx, "email", user)
	if err != nil {
		u.Log.Warnf("failed to read data (user): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if user.Status != "email-verification-needed" {
		u.Log.Warnf("failed to verify account (already verified): %s\n", user.Email)
		return nil, ErrFailedToVerifyAccount
	}

	sixCode, err := u.Redis.Get(ctx, fmt.Sprintf("%s_verification_code_register", user.Email)).Result()
	if err != nil {
		u.Log.Warnf("failed to read data from redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if sixCode != req.Code {
		u.Log.Warnf("invalid token: %s\n", user.Email)
		return nil, ErrInvalidToken
	}

	user.Status = "identity-card-verification-needed"
	err = u.UserRepository.Update(tx, user)
	if err != nil {
		u.Log.Warnf("failed to update data (user): %+v\n", err)
		return nil, ErrFailedToUpdateData
	}

	err = u.Redis.Del(ctx, fmt.Sprintf("%s_verification_code_register", user.Email)).Err()
	if err != nil {
		u.Log.Warnf("failed to delete data from redis: %+v\n", err)
		return nil, ErrFailedToClearData
	}

	now := time.Now().Local().UnixNano()
	res := new(response.FinishRegister)

	jwtAccessToken, err := jwt.CreateToken(user.Username, false)
	if err != nil {
		u.Log.Warnf("failed to create access token: %+v\n", err)
		return nil, ErrFRCreateToken
	}

	jwtRefreshToken, err := jwt.CreateToken(user.Username, true)
	if err != nil {
		u.Log.Warnf("failed to create refresh token: %+v\n", err)
		return nil, ErrFRCreateToken
	}

	refresh := &domain.Refresh{
		ID:            uuid.NewString(),
		Token:         jwtRefreshToken,
		CreatedAt:     now,
		UpdatedAt:     now,
		LastRefreshAt: now,
		UserID:        user.ID,
	}

	err = u.RefreshRepository.Create(tx, refresh)
	if err != nil {
		u.Log.Warnf("failed to store data (refresh): %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	err = tx.Commit()
	if err != nil {
		u.Log.Warnf("failed to commit transaction: %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	res.AccessToken = jwtAccessToken
	res.RefreshToken = jwtRefreshToken

	return res, nil
}

func (u *UserUseCase) Login(ctx context.Context, req *request.Login) (*response.FinishRegister, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	user := &domain.User{
		Email: req.Email,
	}
	err = u.UserRepository.Read(tx, "email", user)
	if err != nil {
		u.Log.Warnf("failed to read data (user): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if user.IsFacebook || user.IsGoogle {
		return nil, ErrLoginTypeOAuth
	}

	if user.Status == "email-verification-needed" {
		return nil, ErrNeedEmailVerification
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrWrongPassword
	}

	res := new(response.FinishRegister)

	accessToken, err := jwt.CreateToken(user.Username, false)
	if err != nil {
		u.Log.Warnf("failed to create access token: %+v\n", err)
		return nil, ErrFRCreateToken
	}

	refresh := &domain.Refresh{
		UserID: user.ID,
	}
	err = u.RefreshRepository.Read(tx, "user_iD", refresh)
	if err != nil {
		u.Log.Warnf("failed to get refresh data: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	refreshToken, err := jwt.CreateToken(user.Username, true)
	if err != nil {
		u.Log.Warnf("failed to create refresh token: %+v\n", err)
		return nil, ErrCreateToken
	}

	refresh.Token = refreshToken
	err = u.RefreshRepository.Update(tx, refresh)
	if err != nil {
		u.Log.Warnf("failed to store data (refresh): %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	err = tx.Commit()
	if err != nil {
		u.Log.Warnf("failed to commit transaction: %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	res.AccessToken = accessToken
	res.RefreshToken = refresh.Token

	return res, nil
}
