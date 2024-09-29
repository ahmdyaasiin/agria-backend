package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	repositoryInterface "github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type MenuUseCase struct {
	DB                  *sqlx.DB
	Log                 *logrus.Logger
	Redis               *redis.Client
	UserRepository      repositoryInterface.UserRepository
	CartRepository      repositoryInterface.CartRepository
	PropertyRepository  repositoryInterface.PropertyRepository
	ProductRepository   repositoryInterface.ProductRepository
	EducationRepository repositoryInterface.EducationRepository
}

func NewMenuUseCase(DB *sqlx.DB, log *logrus.Logger, redis *redis.Client, userRepository repositoryInterface.UserRepository,
	cartRepository repositoryInterface.CartRepository, propertyRepository repositoryInterface.PropertyRepository,
	productRepository repositoryInterface.ProductRepository, educationRepository repositoryInterface.EducationRepository) interfaces.MenuUseCase {
	return &MenuUseCase{DB: DB, Log: log, Redis: redis, UserRepository: userRepository, CartRepository: cartRepository, PropertyRepository: propertyRepository,
		ProductRepository: productRepository, EducationRepository: educationRepository}
}

func (u *MenuUseCase) Homepage(ctx context.Context, userID string) (*response.Homepage, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.Homepage)

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get user detail: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ttlProducts, err := u.Redis.TTL(ctx, "promo_products").Result()
	if err != nil {
		u.Log.Warnf("failed to get expired products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ttlProperty, err := u.Redis.TTL(ctx, "promo_products").Result()
	if err != nil {
		u.Log.Warnf("failed to get expired property: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ids, err := u.Redis.Keys(ctx, "discount_property_id_*").Result()
	if err != nil {
		u.Log.Warnf("failed to read keys in redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var uuids []string
	prefix := "discount_property_id_"

	for _, id := range ids {

		uuid := strings.TrimPrefix(id, prefix)
		uuids = append(uuids, fmt.Sprintf("'%s'", uuid))
	}

	notIN := fmt.Sprintf("(%s)", strings.Join(uuids, ","))

	err = u.PropertyRepository.GetAllPropertiesWithoutPromo(tx, "", user.ID, "newest", notIN, "all", 0, &res.Properties.Data)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetAllPropertiesWithPromo(tx, "", user.ID, "newest", notIN, "all", 0, &res.PropertyPromo.Properties)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ids, err = u.Redis.Keys(ctx, "discount_product_id_*").Result()
	if err != nil {
		u.Log.Warnf("failed to read keys in redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	prefix = "discount_product_id_"

	for _, id := range ids {

		uuid := strings.TrimPrefix(id, prefix)
		uuids = append(uuids, fmt.Sprintf("'%s'", uuid))
	}

	notIN = fmt.Sprintf("(%s)", strings.Join(uuids, ","))

	err = u.ProductRepository.GetAllProductsWithoutPromo(tx, "", user.ID, "newest", notIN, 0, &res.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.ProductRepository.GetAllProductsWithPromo(tx, "", user.ID, "newest", notIN, 0, &res.ProductsPromo.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var total int
	err = u.CartRepository.CountCart(tx, user.ID, &total)
	if err != nil {
		u.Log.Warnf("failed to count cart: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.EducationRepository.GetAllEducation(tx, &res.Educations)
	if err != nil {
		u.Log.Warnf("failed to get educations: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	res.UserDetails.CountCarts = total
	res.UserDetails.IsLoggedIn = user.ID != ""
	res.UserDetails.PhotoProfile = user.PhotoUrl

	res.Properties.Province = "All"
	res.ProductsPromo.TimeLifeInSeconds = int64(ttlProducts.Seconds()) % 86400
	res.PropertyPromo.TimeLifeInSeconds = int64(ttlProperty.Seconds()) % 86400

	return res, nil
}

func (u *MenuUseCase) Market(ctx context.Context, userID string) (*response.Market, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.Market)

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get user detail: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ttlProducts, err := u.Redis.TTL(ctx, "promo_products").Result()
	if err != nil {
		u.Log.Warnf("failed to get expired products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ttlProperty, err := u.Redis.TTL(ctx, "promo_products").Result()
	if err != nil {
		u.Log.Warnf("failed to get expired property: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ids, err := u.Redis.Keys(ctx, "discount_property_id_*").Result()
	if err != nil {
		u.Log.Warnf("failed to read keys in redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var uuids []string
	prefix := "discount_property_id_"

	for _, id := range ids {

		uuid := strings.TrimPrefix(id, prefix)
		uuids = append(uuids, fmt.Sprintf("'%s'", uuid))
	}

	notIN := fmt.Sprintf("(%s)", strings.Join(uuids, ","))

	err = u.PropertyRepository.GetAllPropertiesWithoutPromo(tx, "", user.ID, "newest", notIN, "all", 0, &res.Properties.Data)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetAllPropertiesWithPromo(tx, "", user.ID, "newest", notIN, "all", 0, &res.PropertyPromo.Properties)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ids, err = u.Redis.Keys(ctx, "discount_product_id_*").Result()
	if err != nil {
		u.Log.Warnf("failed to read keys in redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	prefix = "discount_product_id_"

	for _, id := range ids {

		uuid := strings.TrimPrefix(id, prefix)
		uuids = append(uuids, fmt.Sprintf("'%s'", uuid))
	}

	notIN = fmt.Sprintf("(%s)", strings.Join(uuids, ","))

	err = u.ProductRepository.GetAllProductsWithoutPromo(tx, "", user.ID, "newest", notIN, 0, &res.Products.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.ProductRepository.GetAllProductsWithPromo(tx, "", user.ID, "newest", notIN, 0, &res.ProductsPromo.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var total int
	err = u.CartRepository.CountCart(tx, user.ID, &total)
	if err != nil {
		u.Log.Warnf("failed to count cart: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var provinces []string
	err = u.PropertyRepository.GetState(tx, &provinces)
	if err != nil {
		u.Log.Warnf("failed to get state")
		return nil, ErrFailedToReadData
	}

	res.Properties.Pagination.Page = 0
	res.Properties.Pagination.TotalItems = int64(len(res.Properties.Data))
	res.Properties.Pagination.TotalPages = 1

	res.Products.Pagination.Page = 0
	res.Products.Pagination.TotalItems = int64(len(res.Products.Products))
	res.Products.Pagination.TotalPages = 1

	res.Properties.Province = "All"
	res.Properties.Provinces = provinces
	res.ProductsPromo.TimeLifeInSeconds = int64(ttlProducts.Seconds()) % 86400
	res.PropertyPromo.TimeLifeInSeconds = int64(ttlProperty.Seconds()) % 86400

	res.UserDetails.CountCarts = total
	res.UserDetails.IsLoggedIn = user.ID != ""
	res.UserDetails.PhotoProfile = user.PhotoUrl

	return res, nil
}
