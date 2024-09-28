package service

import (
	"fmt"

	"github.com/jasurxaydarov/book_shop_api_getway/genproto/product_service"
	"github.com/jasurxaydarov/book_shop_api_getway/genproto/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Service() ServiceManagerI {

	userService, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	productService, err := grpc.NewClient("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	serviseManager := &serviceManager{
		userService:    user_service.NewUserServiceClient(userService),
		productService: product_service.NewProductServiceClient(productService),
	}

	return serviseManager
}

type ServiceManagerI interface {
	GetUserSevice() user_service.UserServiceClient
	GetProductSevice() product_service.ProductServiceClient
}

type serviceManager struct {
	userService    user_service.UserServiceClient
	productService product_service.ProductServiceClient
}

func (s *serviceManager) GetUserSevice() user_service.UserServiceClient {

	return s.userService
}

func (s *serviceManager) GetProductSevice() product_service.ProductServiceClient {

	return s.productService
}
