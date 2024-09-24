package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	"github.com/ahmdyaasiin/agria-backend/internal/pkg/biteship"
	repositoryInterface "github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
	"time"
)

type ProductUseCase struct {
	//
	DB                     *sqlx.DB
	Log                    *logrus.Logger
	Redis                  *redis.Client
	AddressRepository      repositoryInterface.AddressRepository
	ProductRepository      repositoryInterface.ProductRepository
	ProductMediaRepository repositoryInterface.ProductMediaRepository
	RatingRepository       repositoryInterface.RatingRepository
}

func NewProductUseCase(DB *sqlx.DB,
	log *logrus.Logger,
	redis *redis.Client,
	addressRepository repositoryInterface.AddressRepository,
	productRepository repositoryInterface.ProductRepository,
	productMediaRepository repositoryInterface.ProductMediaRepository,
	ratingRepository repositoryInterface.RatingRepository) interfaces.ProductUseCase {
	return &ProductUseCase{
		DB:                     DB,
		Log:                    log,
		Redis:                  redis,
		AddressRepository:      addressRepository,
		ProductRepository:      productRepository,
		ProductMediaRepository: productMediaRepository,
		RatingRepository:       ratingRepository,
	}
}

func (u *ProductUseCase) GetProducts(ctx context.Context, userID, categoryName, sortBy string, page int) (*response.GetProductWithPagination, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	ids, err := u.Redis.Keys(ctx, "discount_product_id_*").Result()
	if err != nil {
		u.Log.Warnf("failed to read keys in redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var uuids []string
	prefix := "discount_product_id_"

	for _, id := range ids {

		uuid := strings.TrimPrefix(id, prefix)
		uuids = append(uuids, fmt.Sprintf("'%s'", uuid))
	}

	notIN := fmt.Sprintf("(%s)", strings.Join(uuids, ","))

	res := new(response.GetProductWithPagination)
	err = u.ProductRepository.GetAllProductsWithoutPromo(tx, categoryName, userID, sortBy, notIN, page, &res.Products)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	res.Page = page
	res.TotalItems = int64(len(res.Products))
	res.TotalPages = int64(math.Ceil(float64(res.TotalItems) / 24))

	return res, nil
}

func (u *ProductUseCase) GetProductDetails(ctx context.Context, userID, productID string) (*response.GetProductDetails, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	product := new(response.GetProductDetails)
	err = u.ProductRepository.GetDetailsProduct(tx, productID, userID, product)
	if err != nil {
		u.Log.Warnf("failed to get product: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.ProductMediaRepository.GetProductMedia(tx, productID, &product.PhotoUrls)
	if err != nil {
		u.Log.Warnf("failed to get product media: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.RatingRepository.GetProductReviews(tx, productID, userID, "newest", 1, &product.Reviews)
	if err != nil {
		u.Log.Warnf("failed to get product reviews: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	for i, r := range product.Reviews {
		photoUrls := strings.Split(r.PhotoReviewUrlsString, ",")
		product.Reviews[i].PhotoReviewUrls = photoUrls
	}

	// set default to rektorat ub wkwk
	latitude, longitude := -6.213231948641893, 106.79724408707149

	if userID != "" {
		primaryAddress := &domain.Address{
			UserID: userID,
		}

		err = u.AddressRepository.Read(tx, "user_iD", primaryAddress)
		if err != nil {
			return nil, ErrFailedToReadData
		}

		latitude = primaryAddress.Latitude
		longitude = primaryAddress.Longitude
	}

	resString, err := u.Redis.Get(ctx, fmt.Sprintf("shipping_information_for_%f_%f", latitude, longitude)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		u.Log.Warnf("failed to read data from redis: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	res := biteship.ResponseShippingForProductDetail{}
	if errors.Is(err, redis.Nil) {
		items := []biteship.Items{
			{
				Name:     product.Name,
				Value:    product.Price,
				Weight:   product.UnitWeight,
				Quantity: 1,
			},
		}
		res, err = biteship.ShippingForProductDetails(latitude, longitude, items)
		if err != nil {
			u.Log.Warnf("failed to get shipping from biteship: %+v\n", err)
			return nil, ErrCalculateShipping
		}

		resMarshal, err := json.Marshal(res)
		if err != nil {
			u.Log.Warnf("failed to marshal: %+v\n", err)
			return nil, ErrFailedToMarshal
		}

		err = u.Redis.Set(ctx, fmt.Sprintf("shipping_information_for_%f_%f", latitude, longitude), resMarshal, 10*time.Minute).Err()
		if err != nil {
			u.Log.Warnf("failed to store data to redis: %+v\n", err)
			return nil, ErrFailedToStoreData
		}
	} else {
		err = json.Unmarshal([]byte(resString), &res)
		if err != nil {
			u.Log.Warnf("failed to unmarshal: %+v\n", err)
			return nil, ErrFailedToUnMarshal
		}
	}

	product.PriceRange = res.CostRange
	product.TimeRange = res.EstimatedArrived

	return product, nil
}

func (u *ProductUseCase) GetProductReviews(ctx context.Context, userID, productID, sortBy string, page int) (*response.ReviewDetails, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.ReviewDetails)
	err = u.RatingRepository.GetProductReviews(tx, productID, userID, sortBy, page, &res.Reviews)
	if err != nil {
		u.Log.Warnf("failed to read data (product reviews): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	ratingBreakdown := new([]response.RatingBreakdown)
	err = u.RatingRepository.RatingBreakdown(tx, productID, ratingBreakdown)
	if err != nil {
		u.Log.Warnf("failed to get rating breakdown: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	res.RatingBreakdown = make([]int64, 5)
	for i := 5; i >= 1; i-- {
		total := int64(0)

		for _, rating := range *ratingBreakdown {
			if rating.Star == i {
				total = rating.Total
				break
			}
		}

		res.RatingBreakdown[5-i] = total
	}

	res.Page = page
	res.TotalItems = int64(len(res.Reviews))
	res.TotalPages = int64(math.Ceil(float64(res.TotalItems) / 5))

	return res, nil
}
