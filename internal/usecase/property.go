package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	repositoryInterface "github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
	"time"
)

type PropertyUseCase struct {
	DB                       *sqlx.DB
	Log                      *logrus.Logger
	Redis                    *redis.Client
	PropertyRepository       repositoryInterface.PropertyRepository
	PropertyRatingRepository repositoryInterface.PropertyRatingRepository
	DiscussRepository        repositoryInterface.DiscussRepository
	WishlistRepository       repositoryInterface.WishlistRepository
	UserRepository           repositoryInterface.UserRepository
	CartRepository           repositoryInterface.CartRepository
}

func NewPropertyUseCase(DB *sqlx.DB, log *logrus.Logger, redis *redis.Client,
	propertyRepository repositoryInterface.PropertyRepository,
	propertyRatingRepository repositoryInterface.PropertyRatingRepository,
	discussRepository repositoryInterface.DiscussRepository,
	wishlistRepository repositoryInterface.WishlistRepository,
	userRepository repositoryInterface.UserRepository,
	cartRepository repositoryInterface.CartRepository) interfaces.PropertyUseCase {
	return &PropertyUseCase{DB: DB, Log: log, Redis: redis, PropertyRepository: propertyRepository,
		PropertyRatingRepository: propertyRatingRepository,
		DiscussRepository:        discussRepository,
		WishlistRepository:       wishlistRepository,
		UserRepository:           userRepository,
		CartRepository:           cartRepository}
}

func (u *PropertyUseCase) GetAllWishlistsProperties(ctx context.Context, userID string) (*response.PropertiesWishlist, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.PropertiesWishlist)

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil {
		u.Log.Warnf("failed to read user details: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.WishlistRepository.GetMyWishlistsProperty(tx, userID, &res.Properties)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get wishlist products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if len(res.Properties) != 0 {
		keys := strings.Split((res.Properties)[0].ProductIDString, ",")
		discountProducts, err := u.Redis.MGet(ctx, keys...).Result()
		if err != nil {
			u.Log.Warnf("failed to get discount available products: %+v\n", err)
			return nil, ErrFailedToReadData
		}

		for i, discount := range discountProducts {
			if discount == nil {
				continue
			}

			d, err := strconv.Atoi(discount.(string))
			if err != nil {
				return nil, ErrParseStringToNumber
			}

			(res.Properties)[i].DiscountPrice = int64(d)
		}
	}

	var total int
	err = u.CartRepository.CountCart(tx, user.ID, &total)
	if err != nil {
		u.Log.Warnf("failed to get count cart: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	res.UserDetails.IsLoggedIn = user.ID != ""
	res.UserDetails.CountCarts = total
	res.UserDetails.PhotoProfile = user.PhotoUrl

	res.Pagination.Page = 0
	res.Pagination.TotalItems = int64(len(res.Properties))
	res.Pagination.TotalPages = 1

	return res, nil
}

func (u *PropertyUseCase) ManageWishlistProperties(ctx context.Context, userID string, req *request.ManageWishlistProperties) (*response.ManageWishlistProperties, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	wishlist := new(domain.PropertyWishlist)
	err = u.WishlistRepository.GetSpecificProperty(tx, userID, req.PropertyID, wishlist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get specific product from wishlist: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var isWishlisted bool
	if wishlist.ID != "" {
		err = u.WishlistRepository.DeleteProperty(tx, wishlist)
		if err != nil {
			u.Log.Warnf("failed to delete product wishlist: %+v\n", err)
			return nil, ErrFailedToClearData
		}
	} else {
		isWishlisted = true
		wishlist.ID = uuid.NewString()
		wishlist.UserID = userID
		wishlist.PropertyID = req.PropertyID
		wishlist.CreatedAt = time.Now().Local().UnixNano()

		err = u.WishlistRepository.CreateProperty(tx, wishlist)
		if err != nil {
			u.Log.Warnf("failed to store data product wishlist: %+v\n", err)
			return nil, ErrFailedToStoreData
		}
	}

	err = tx.Commit()
	if err != nil {
		u.Log.Warnf("failed to commit transaction: %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	res := &response.ManageWishlistProperties{
		PropertiesID: req.PropertyID,
		IsWishlisted: isWishlisted,
	}

	return res, nil
}

func (u *PropertyUseCase) GetProperties(ctx context.Context, userID, categoryName, sortBy, province string, page int) (*response.GetPropertiesWithPagination, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil {
		u.Log.Warnf("failed to read user details: %+v\n", err)
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

	res := new(response.GetPropertiesWithPagination)
	err = u.PropertyRepository.GetAllPropertiesWithoutPromo(tx, categoryName, userID, sortBy, notIN, province, page, &res.Properties.Data)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var total int
	err = u.CartRepository.CountCart(tx, userID, &total)
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

	res.Provinces = provinces
	res.Properties.Province = province

	res.UserDetails.IsLoggedIn = user.ID != ""
	res.UserDetails.CountCarts = total
	res.UserDetails.PhotoProfile = user.PhotoUrl

	res.Pagination.Page = page
	res.Pagination.TotalItems = int64(len(res.Properties.Data))
	res.Pagination.TotalPages = int64(math.Ceil(float64(res.Pagination.TotalItems) / 24))

	return res, nil
}

func (u *PropertyUseCase) GetPropertyDetails(ctx context.Context, userID, propertyID string) (*response.GetPropertyDetails, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil {
		u.Log.Warnf("failed to get user details: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	product := new(response.GetPropertyDetails)
	err = u.PropertyRepository.GetPropertyDetails(tx, propertyID, userID, product)
	if err != nil {
		u.Log.Warnf("failed to get product: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetPropertyHighlights(tx, propertyID, &product.Highlights)
	if err != nil {
		u.Log.Warnf("failed to get property highlights: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetPropertyMedia(tx, propertyID, &product.PhotoUrls)
	if err != nil {
		u.Log.Warnf("failed to get property media: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	err = u.PropertyRepository.GetPropertyDiscuss(tx, propertyID, &product.Discuss)
	if err != nil {
		u.Log.Warnf("failed to get property discuss: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	for i, d := range product.Discuss {

		rep := strings.Split(d.AnswersString, "$$$$$$$$$$$$$$$$$$,")
		product.Discuss[i].Answers = make([]response.PropertyDiscussReplies, len(rep))
		for j, r := range rep {
			reply := strings.Split(r, "_")

			isOwner := strings.ReplaceAll(reply[4], "$$$$$$$$$$$$$$$$$$", "")

			product.Discuss[i].Answers[j].ID = reply[0]
			product.Discuss[i].Answers[j].Content = reply[1]
			product.Discuss[i].Answers[j].Name = reply[2]
			product.Discuss[i].Answers[j].PhotoUrl = reply[3]
			product.Discuss[i].Answers[j].IsOwner = isOwner == "1"
		}
	}

	err = u.PropertyRepository.GetPropertyRatings(tx, propertyID, userID, &product.RatingsAndReviews.Data)
	if err != nil {
		u.Log.Warnf("failed to get property ratings: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	for i, r := range product.RatingsAndReviews.Data {
		photoUrlsSlice := strings.Split(r.PhotoUrlsString, ",")
		product.RatingsAndReviews.Data[i].PhotoUrls = make([]string, len(photoUrlsSlice))

		for j, photo := range photoUrlsSlice {
			product.RatingsAndReviews.Data[i].PhotoUrls[j] = photo
		}
	}

	discountPriceString, err := u.Redis.Get(ctx, fmt.Sprintf("discount_property_id_%s", product.ID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		u.Log.Warnf("failed to read discount price: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var discountPrice int
	if !errors.Is(err, redis.Nil) {
		discountPrice, err = strconv.Atoi(discountPriceString)
		if err != nil {
			return nil, ErrParseStringToNumber
		}
	}

	ratingBreakdown := new([]response.RatingBreakdown)
	err = u.PropertyRepository.RatingBreakdown(tx, propertyID, ratingBreakdown)
	if err != nil {
		u.Log.Warnf("failed to get rating breakdown: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	product.RatingsAndReviews.CountStarBreakDown = make([]int, 5)
	for i := 5; i >= 1; i-- {
		total := 0

		for _, rating := range *ratingBreakdown {
			if rating.Star == i {
				total = int(rating.Total)
				break
			}
		}

		product.RatingsAndReviews.CountStarBreakDown[5-i] = total
	}

	var total int
	err = u.CartRepository.CountCart(tx, user.ID, &total)
	if err != nil {
		u.Log.Warnf("failed to count cart %+v\n", err)
		return nil, ErrFailedToReadData
	}

	product.RatingsAndReviews.CountRatings = len(product.RatingsAndReviews.Data)
	product.DiscountPrice = int64(discountPrice)

	product.UserDetails.IsLoggedIn = user.ID != ""
	product.UserDetails.CountCarts = total
	product.UserDetails.PhotoProfile = user.PhotoUrl

	return product, nil
}
