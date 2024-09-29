package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/request"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"github.com/ahmdyaasiin/agria-backend/internal/domain"
	repositoryInterface "github.com/ahmdyaasiin/agria-backend/internal/repository/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/usecase/interfaces"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type WishlistUseCase struct {
	DB                 *sqlx.DB
	Log                *logrus.Logger
	Redis              *redis.Client
	WishlistRepository repositoryInterface.WishlistRepository
	UserRepository     repositoryInterface.UserRepository
	CartRepository     repositoryInterface.CartRepository
}

func NewWishlistUseCase(DB *sqlx.DB, log *logrus.Logger, redis *redis.Client, wishlistRepository repositoryInterface.WishlistRepository, userRepository repositoryInterface.UserRepository, cartRepository repositoryInterface.CartRepository) interfaces.WishlistUseCase {
	return &WishlistUseCase{DB: DB, Log: log, Redis: redis, WishlistRepository: wishlistRepository, UserRepository: userRepository, CartRepository: cartRepository}
}

func (u *WishlistUseCase) GetAllWishlists(ctx context.Context, userID string) (*response.ProductWishlist, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	res := new(response.ProductWishlist)

	user := &domain.User{
		ID: userID,
	}
	err = u.UserRepository.Read(tx, "iD", user)
	if err != nil {
		u.Log.Warnf("failed to get user details: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var total int
	err = u.CartRepository.CountCart(tx, user.ID, &total)
	if err != nil {
		u.Log.Warnf("failed to get count error: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	wishlists := new([]response.MyWishlist)
	err = u.WishlistRepository.GetMyWishlists(tx, userID, wishlists)
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

	res.UserDetails.CountCarts = total
	res.UserDetails.IsLoggedIn = user.ID != ""
	res.UserDetails.PhotoProfile = user.PhotoUrl
	res.Products = *wishlists
	res.Pagination.TotalItems = int64(len(*wishlists))
	res.Pagination.Page = 1
	res.Pagination.TotalPages = 1

	return res, nil
}

func (u *WishlistUseCase) ManageWishlist(ctx context.Context, userID string, req *request.ManageWishlist) (*response.ManageWishlist, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	wishlist := new(domain.Wishlist)
	err = u.WishlistRepository.GetSpecificProduct(tx, userID, req.ProductID, wishlist)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to get specific product from wishlist: %+v\n", err)
		return nil, ErrFailedToReadData
	}

	var isWishlisted bool
	if wishlist.ID != "" {
		err = u.WishlistRepository.Delete(tx, wishlist)
		if err != nil {
			u.Log.Warnf("failed to delete product wishlist: %+v\n", err)
			return nil, ErrFailedToClearData
		}
	} else {
		isWishlisted = true
		wishlist.ID = uuid.NewString()
		wishlist.UserID = userID
		wishlist.ProductID = req.ProductID
		wishlist.CreatedAt = time.Now().Local().UnixNano()

		err = u.WishlistRepository.Create(tx, wishlist)
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

	res := &response.ManageWishlist{
		ProductID:    req.ProductID,
		IsWishlisted: isWishlisted,
	}

	return res, nil
}
