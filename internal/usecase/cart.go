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

type CartUseCase struct {
	DB             *sqlx.DB
	Log            *logrus.Logger
	Redis          *redis.Client
	CartRepository repositoryInterface.CartRepository
}

func NewCartUseCase(DB *sqlx.DB, log *logrus.Logger, redis *redis.Client, cartRepository repositoryInterface.CartRepository) interfaces.CartUseCase {
	return &CartUseCase{DB: DB, Log: log, Redis: redis, CartRepository: cartRepository}
}

func (u *CartUseCase) GetMyCart(ctx context.Context, userID string) (*response.MyCart, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	products := new(response.MyCart)
	err = u.CartRepository.GetMyCartAvailable(tx, userID, &products.AvailableProducts)
	if err != nil {
		u.Log.Warnf("failed to get products in cart (available): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if len(products.AvailableProducts) != 0 {
		keys := strings.Split(products.AvailableProducts[0].ProductIDString, ",")
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

			products.AvailableProducts[i].DiscountPrice = int64(d)
		}
	}

	err = u.CartRepository.GetMyCartUnavailable(tx, userID, &products.UnavailableProducts)
	if err != nil {
		u.Log.Warnf("failed to get products in cart (unavailable): %+v\n", err)
		return nil, ErrParseStringToNumber
	}

	if len(products.UnavailableProducts) != 0 {
		keys := strings.Split(products.UnavailableProducts[0].ProductIDString, ",")
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

			products.UnavailableProducts[i].DiscountPrice = int64(d)
		}
	}

	return products, nil
}

func (u *CartUseCase) ManageCart(ctx context.Context, userID string, req *request.ManageCart) (*response.ManageCart, error) {
	tx, err := u.DB.Beginx()
	defer tx.Rollback()
	if err != nil {
		u.Log.Warnf("create transaction: %+v\n", err)
		return nil, ErrCreateDatabaseTransaction
	}

	cart := new(domain.Cart)
	err = u.CartRepository.GetMyCart(tx, userID, req.ProductID, cart)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.Log.Warnf("failed to read data (cart): %+v\n", err)
		return nil, ErrFailedToReadData
	}

	if req.Quantity > cart.Quantity {
		return nil, ErrOutOfStock
	}

	now := time.Now().Local().UnixNano()
	cart.Quantity = req.Quantity
	cart.UpdatedAt = now

	if cart.ID != "" {

		if cart.Quantity == 0 {
			//delete
			err = u.CartRepository.Delete(tx, cart)
			if err != nil {
				u.Log.Warnf("failed to delete product in cart: %+v\n", err)
				return nil, ErrFailedToClearData
			}
		} else {
			//update
			err = u.CartRepository.Update(tx, cart)
			if err != nil {
				u.Log.Warnf("failed to update product in cart: %+v\n", err)
				return nil, ErrFailedToUpdateData
			}
		}
	} else if cart.Quantity != 0 {
		cart.ID = uuid.NewString()
		cart.CreatedAt = now
		cart.UserID = userID
		cart.ProductID = req.ProductID

		// create
		err = u.CartRepository.Create(tx, cart)
		if err != nil {
			u.Log.Warnf("failed to store product to cart: %+v\n", err)
			return nil, ErrFailedToStoreData
		}
	}

	err = tx.Commit()
	if err != nil {
		u.Log.Warnf("failed to commit transaciton: %+v\n", err)
		return nil, ErrFailedToStoreData
	}

	res := &response.ManageCart{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	return res, nil
}
