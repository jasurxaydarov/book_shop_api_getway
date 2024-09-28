package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jasurxaydarov/book_shop_api_getway/api/handlers"
	"github.com/jasurxaydarov/book_shop_api_getway/api/middlewares"
	"github.com/jasurxaydarov/book_shop_api_getway/redis"
	"github.com/jasurxaydarov/book_shop_api_getway/service"

	"github.com/saidamir98/udevs_pkg/logger"
)

type Options struct {
	Service service.ServiceManagerI
	Log     logger.LoggerI
	Cache   redis.RedisRepoI
}

func Api(o Options) *gin.Engine {

	h := handlers.NewHandlers(o.Service, o.Log, o.Cache)

	engine := gin.Default()

	api := engine.Group("/api")
	us := api.Group("/us")

	us.Use(middlewares.AuthMiddlewareUser())
	{
		us.GET("/user/:id", h.GetUserById)
		us.POST("/update", h.UpdateUser)
		us.POST("/delete", h.DeleteUser)

		//order
		us.POST("/order", h.CreateOrder)
		us.GET("/order/:id", h.GetOrderById)

		//Author
		us.GET("/auth/:id", h.GetAuthById)

		// book
		us.GET("/book/:id", h.GetBookById)

		// orderItem
		us.POST("/order_item", h.CreateOrderItem)
		us.GET("/order_item/:id", h.GetOrderItemById)
		us.GET("/order_item_id/:id", h.GetOrderItemById)

	}

	adm := api.Group("/adm")

	adm.Use(middlewares.AuthMiddlewareAdmin())
	{
		// author
		adm.POST("/auth", h.CreateAuth)
		adm.GET("/auth/:id", h.GetAuthById)

		/////
		adm.POST("/update", h.UpdateUser)
		adm.POST("/delete", h.DeleteUser)

		////////
		adm.POST("/update_user", h.AdmUpdateUser)
		adm.POST("/delete_user", h.AdmDeleteUser)
		adm.POST("/get_users", h.GetUsers)





		//category
		adm.POST("/category", h.CreateCategory)
		adm.GET("/category/:id", h.GetCategoryById)

		//book
		adm.POST("/book", h.CreateBook)
		adm.GET("/book/:id", h.GetBookById)

	}

	all := api.Group("/all")

	{
		all.POST("/check-user", h.CheckUser) //completed
		all.POST("/sign-up", h.SignUp)       //completed
		all.POST("/sign-in", h.SigIn)        //completed

	}
	return engine

}
