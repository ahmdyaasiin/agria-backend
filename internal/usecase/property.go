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
}

func NewPropertyUseCase(DB *sqlx.DB, log *logrus.Logger, redis *redis.Client,
	propertyRepository repositoryInterface.PropertyRepository,
	propertyRatingRepository repositoryInterface.PropertyRatingRepository,
	discussRepository repositoryInterface.DiscussRepository,
	wishlistRepository repositoryInterface.WishlistRepository) interfaces.PropertyUseCase {
	return &PropertyUseCase{DB: DB, Log: log, Redis: redis, PropertyRepository: propertyRepository,
		PropertyRatingRepository: propertyRatingRepository,
		DiscussRepository:        discussRepository,
		WishlistRepository:       wishlistRepository}
}

func (u *PropertyUseCase) GetAllWishlistsProperties(ctx context.Context, userID string) (*[]response.MyWishlistProperties, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	wishlists := new([]response.MyWishlistProperties)
	err = u.WishlistRepository.GetMyWishlistsProperty(tx, userID, wishlists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get wishlist products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if len(*wishlists) != 0 {
		keys := strings.Split((*wishlists)[0].ProductIDString, ",")
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

			(*wishlists)[i].DiscountPrice = int64(d)
		}
	}

	return wishlists, nil
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
	err = u.PropertyRepository.GetAllPropertiesWithoutPromo(tx, categoryName, userID, sortBy, notIN, province, page, &res.Properties)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get products: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	res.Page = page
	res.TotalItems = int64(len(res.Properties))
	res.TotalPages = int64(math.Ceil(float64(res.TotalItems) / 24))

	return res, nil
}

func (u *PropertyUseCase) GetPropertyDetails(ctx context.Context, userID, propertyID string) (*response.GetPropertyDetails, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
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

	err = u.PropertyRepository.GetPropertyRatings(tx, propertyID, userID, &product.RatingsAndReviews)
	if err != nil {
		u.Log.Warnf("failed to get property ratings: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	for i, r := range product.RatingsAndReviews {
		photoUrlsSlice := strings.Split(r.PhotoUrlsString, ",")
		product.RatingsAndReviews[i].PhotoUrls = make([]string, len(photoUrlsSlice))

		for j, photo := range photoUrlsSlice {
			product.RatingsAndReviews[i].PhotoUrls[j] = photo
		}
	}

	return product, nil
}
