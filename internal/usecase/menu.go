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
	"math"
	"strconv"
	"strings"
	"time"
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

	ttlProperty, err := u.Redis.TTL(ctx, "promo_property").Result()
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

	err = u.PropertyRepository.GetAllPropertiesWithoutPromo(tx, "", user.ID, "newest", notIN, "all", 1, 6, &res.Properties.Data)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetAllPropertiesWithPromo(tx, "", user.ID, "newest", notIN, "all", 0, &res.PropertyPromo.Properties)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var redisKeys []string
	for _, product := range res.PropertyPromo.Properties {
		redisKey := fmt.Sprintf("discount_property_id_%s", product.ID)
		redisKeys = append(redisKeys, redisKey)
	}

	discounts, err := u.Redis.MGet(ctx, redisKeys...).Result()
	if err != nil {
		u.Log.Warnf("failed to get discount price: %+v\n", err.Error())
		return nil, ErrFailedToReadData
	}

	for i, discount := range discounts {
		d, err := strconv.Atoi(discount.(string))
		if err != nil {
			continue
		}

		if discount != nil {
			res.PropertyPromo.Properties[i].DiscountPrice = int64(d)
		}
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

	err = u.ProductRepository.GetAllProductsWithoutPromo(tx, "", user.ID, "newest", notIN, 1, 6, &res.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.ProductRepository.GetAllProductsWithPromo(tx, "", user.ID, "newest", notIN, 0, &res.ProductsPromo.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	redisKeys = []string{}
	for _, product := range res.ProductsPromo.Products {
		redisKey := fmt.Sprintf("discount_product_id_%s", product.ID)
		redisKeys = append(redisKeys, redisKey)
	}

	discounts, err = u.Redis.MGet(ctx, redisKeys...).Result()
	if err != nil {
		u.Log.Warnf("failed to get discount price (product): %+v\n", err.Error())
		return nil, ErrFailedToReadData
	}

	for i, discount := range discounts {
		d, err := strconv.Atoi(discount.(string))
		if err != nil {
			continue
		}

		if discount != nil {
			res.ProductsPromo.Products[i].DiscountPrice = int64(d)
		}
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

	for i, e := range res.Educations {
		t := time.Unix(0, e.CreatedAt)
		f := t.Format("2 January 2006")

		res.Educations[i].CreatedAtString = f
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

	err = u.PropertyRepository.GetAllPropertiesWithoutPromo(tx, "", user.ID, "newest", notIN, "all", 1, 18, &res.Properties.Data)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products22: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetAllPropertiesWithPromo(tx, "", user.ID, "newest", notIN, "all", 0, &res.PropertyPromo.Properties)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products2: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var redisKeys []string
	for _, product := range res.PropertyPromo.Properties {
		redisKey := fmt.Sprintf("discount_property_id_%s", product.ID)
		redisKeys = append(redisKeys, redisKey)
	}

	discounts, err := u.Redis.MGet(ctx, redisKeys...).Result()
	if err != nil {
		u.Log.Warnf("failed to get discount price: %+v\n", err.Error())
		return nil, ErrFailedToReadData
	}

	for i, discount := range discounts {
		d, err := strconv.Atoi(discount.(string))
		if err != nil {
			continue
		}

		if discount != nil {
			res.PropertyPromo.Properties[i].DiscountPrice = int64(d)
		}
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

	err = u.ProductRepository.GetAllProductsWithoutPromo(tx, "", user.ID, "newest", notIN, 1, 18, &res.Products.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products1: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.ProductRepository.GetAllProductsWithPromo(tx, "", user.ID, "newest", notIN, 0, &res.ProductsPromo.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products3: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	redisKeys = []string{}
	for _, product := range res.ProductsPromo.Products {
		redisKey := fmt.Sprintf("discount_product_id_%s", product.ID)
		redisKeys = append(redisKeys, redisKey)
	}

	discounts, err = u.Redis.MGet(ctx, redisKeys...).Result()
	if err != nil {
		u.Log.Warnf("failed to get discount price (product): %+v\n", err.Error())
		return nil, ErrFailedToReadData
	}

	for i, discount := range discounts {
		d, err := strconv.Atoi(discount.(string))
		if err != nil {
			continue
		}

		if discount != nil {
			res.ProductsPromo.Products[i].DiscountPrice = int64(d)
		}
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

	res.Properties.Pagination.Page = 1
	res.Properties.Pagination.TotalItems = int64(len(res.Properties.Data))
	res.Properties.Pagination.TotalPages = 1

	res.Products.Pagination.Page = 1
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

func (u *MenuUseCase) Education(ctx context.Context, userID string) (*response.Education, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.Education)

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get user detail: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	idMainArticle, err := u.Redis.Get(ctx, "education_main_article").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		u.Log.Warnf("failed to get data from redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.EducationRepository.MainArticle(tx, idMainArticle, &res.MainArticle)
	if err != nil {
		u.Log.Warnf("failed to get data (main article): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	t := time.Unix(0, res.MainArticle.CreatedAt)
	f := t.Format("2 January 2006")

	res.MainArticle.CreatedAtString = f

	idMustRead, err := u.Redis.Get(ctx, "education_must_read").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		u.Log.Warnf("failed to get data from redis (must read): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var ids []string
	idsMustRead := strings.Split(idMustRead, ",")
	for _, id := range idsMustRead {

		ids = append(ids, fmt.Sprintf("'%s'", id))
	}

	err = u.EducationRepository.MustRead(tx, strings.Join(ids, ","), &res.MustRead)
	if err != nil {
		u.Log.Warnf("failed to get data (must read): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	for i, e := range res.MustRead {
		t := time.Unix(0, e.CreatedAt)
		f := t.Format("2 January 2006")

		res.MustRead[i].CreatedAtString = f
	}

	ids = append(ids, fmt.Sprintf("'%s'", idMainArticle))

	err = u.EducationRepository.Latest(tx, strings.Join(ids, ","), &res.Latest)
	if err != nil {
		u.Log.Warnf("failed to get data (latest): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ids = []string{}
	for i, e := range res.Latest {
		t := time.Unix(0, e.CreatedAt)
		f := t.Format("2 January 2006")

		res.Latest[i].CreatedAtString = f

		ids = append(ids, fmt.Sprintf("'%s'", e.ID))
	}

	err = u.EducationRepository.ExceptionWithRandom(tx, strings.Join(ids, ","), &res.DiscoverMore.Data)
	if err != nil {
		u.Log.Warnf("failed to get data (discover more): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	for i, e := range res.DiscoverMore.Data {
		t := time.Unix(0, e.CreatedAt)
		f := t.Format("2 January 2006")

		res.DiscoverMore.Data[i].CreatedAtString = f
	}

	res.DiscoverMore.Pagination.Page = 1
	res.DiscoverMore.Pagination.TotalItems = int64(len(res.DiscoverMore.Data))
	res.DiscoverMore.Pagination.TotalPages = int64(math.Ceil(float64(res.DiscoverMore.Pagination.TotalItems / 6)))

	return res, nil
}

func (u *MenuUseCase) EducationDetails(ctx context.Context, userID, educationID string) (*response.EducationDetails, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.EducationDetails)

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get user detail: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var total int
	err = u.CartRepository.CountCart(tx, user.ID, &total)
	if err != nil {
		u.Log.Warnf("failed to count cart: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.EducationRepository.EducationDetails(tx, educationID, &res.Data)
	if err != nil {
		u.Log.Warnf("failed to get education details: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.EducationRepository.ExceptionWithRandom(tx, "'as'", &res.RelatedArticle.Data)
	if err != nil {
		u.Log.Warnf("failed to related articles: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	t := time.Unix(0, res.Data.CreatedAt)
	f := t.Format("2 January 2006")

	res.Data.CreatedAtString = f

	res.RelatedArticle.Pagination.Page = 1
	res.RelatedArticle.Pagination.TotalItems = int64(len(res.RelatedArticle.Data))
	res.RelatedArticle.Pagination.TotalPages = int64(math.Ceil(float64(res.RelatedArticle.Pagination.TotalItems / 4)))

	res.UserDetails.CountCarts = total
	res.UserDetails.IsLoggedIn = user.ID != ""
	res.UserDetails.PhotoProfile = user.PhotoUrl

	return res, nil
}
