# Agria Backend API
<p align="center">
  <img src="assets/agria-logo.png" width="108" alt="Agria Logo"/>
</p>

## Project Overview
**Agria** is an innovative application aimed at modernizing the agricultural sector in Indonesia by providing a platform for the rental and management of agricultural land. In the face of economic growth and urbanization challenges, Agria offers easy access for urban communities to engage in agriculture through features such as educational guides, assistance in finding customers and investors, as well as up-to-date information about land conditions and weather.
## Technology Stack:
- Go : https://github.com/golang/go
- MySQL : https://github.com/mysql/mysql-server
- Fiber : https://github.com/gofiber/fiber
- SQLX : https://github.com/jmoiron/sqlx
- JWT : https://github.com/golang-jwt/jwt
- Redis : https://github.com/redis/go-redis
## Folder and File Structure
```
.
├── Dockerfile
├── LICENSE
├── README.md
├── api
│   └── api-spec.json
├── cmd
│   └── app
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   ├── app.go
│   │   ├── env.go
│   │   ├── fiber.go
│   │   ├── logrus.go
│   │   ├── oauth.go
│   │   ├── redis.go
│   │   ├── sqlx.go
│   │   └── validator.go
│   ├── delivery
│   │   └── http
│   │       ├── handler
│   │       │   ├── cart.go
│   │       │   ├── error.go
│   │       │   ├── interfaces
│   │       │   │   ├── cart.go
│   │       │   │   ├── menu.go
│   │       │   │   ├── product.go
│   │       │   │   ├── property.go
│   │       │   │   ├── user.go
│   │       │   │   └── wishlist.go
│   │       │   ├── menu.go
│   │       │   ├── product.go
│   │       │   ├── property.go
│   │       │   ├── request
│   │       │   │   ├── cart.go
│   │       │   │   ├── user.go
│   │       │   │   └── wishlist.go
│   │       │   ├── response
│   │       │   │   ├── cart.go
│   │       │   │   ├── menu.go
│   │       │   │   ├── product.go
│   │       │   │   ├── response.go
│   │       │   │   ├── user.go
│   │       │   │   └── wishlist.go
│   │       │   ├── user.go
│   │       │   └── wishlist.go
│   │       ├── middleware
│   │       │   ├── error.go
│   │       │   └── middleware.go
│   │       ├── route
│   │       │   ├── admin.go
│   │       │   ├── ping.go
│   │       │   └── user.go
│   │       └── server.go
│   ├── domain
│   │   ├── address.go
│   │   ├── cart.go
│   │   ├── product.go
│   │   ├── property_wishlist.go
│   │   ├── refresh.go
│   │   ├── user.go
│   │   └── wishlist.go
│   ├── mock
│   │   ├── mockrepo
│   │   │   └── user_mock.go
│   │   └── mockusecase
│   │       └── user_mock.go
│   ├── pkg
│   │   ├── biteship
│   │   │   ├── biteship.go
│   │   │   └── var.go
│   │   ├── jwt
│   │   │   └── jwt.go
│   │   ├── oauth
│   │   │   └── oauth.go
│   │   ├── query
│   │   │   └── query.go
│   │   ├── validation
│   │   │   └── validation.go
│   │   └── verification
│   │       ├── mail.go
│   │       └── utils.go
│   ├── repository
│   │   ├── address.go
│   │   ├── cart.go
│   │   ├── discuss.go
│   │   ├── education.go
│   │   ├── interfaces
│   │   │   ├── address.go
│   │   │   ├── cart.go
│   │   │   ├── discuss.go
│   │   │   ├── education.go
│   │   │   ├── product.go
│   │   │   ├── product_media.go
│   │   │   ├── property.go
│   │   │   ├── property_rating.go
│   │   │   ├── rating.go
│   │   │   ├── rating_media.go
│   │   │   ├── refresh.go
│   │   │   ├── user.go
│   │   │   └── wishlist.go
│   │   ├── product.go
│   │   ├── product_media.go
│   │   ├── property.go
│   │   ├── property_rating.go
│   │   ├── query.go
│   │   ├── rating.go
│   │   ├── rating_media.go
│   │   ├── refresh.go
│   │   ├── user.go
│   │   └── wishlist.go
│   └── usecase
│       ├── cart.go
│       ├── error.go
│       ├── interfaces
│       │   ├── cart.go
│       │   ├── menu.go
│       │   ├── product.go
│       │   ├── property.go
│       │   ├── user.go
│       │   └── wishlist.go
│       ├── menu.go
│       ├── product.go
│       ├── property.go
│       ├── user.go
│       └── wishlist.go
├── test
│   ├── delivery
│   │   └── http
│   │       └── user_test.go
│   ├── repository
│   │   └── user_test.go
│   └── usecase
│       └── user_test.go
    
34 directories, 103 files
```

## Live Demo
Documentation: https://agria-api-spec.vercel.app/ \
Backend Server: https://agria-backend.iyh.me/v1