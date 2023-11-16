package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"main.go/internal/web/handler"
	"main.go/internal/web/middleware"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, productHandler *handler.ProductHandler, superadminHandler *handler.SuperAdminHandler, carrtHandler *handler.CartHandler) *ServerHTTP {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	home := engine.Group("/home")
	{
		home.GET("/", productHandler.ListAllProducts)
		home.GET("/:productItem_id", productHandler.DisplayProductItem)
		home.GET("/brands", productHandler.ListAllBrands)
		home.GET("/brands/:brand_id", productHandler.DisplayBrand)
		home.GET("/categories", productHandler.ListAllCategories)
		home.GET("/categories/:id", productHandler.DisplayCategory)
	}
	user := engine.Group("/user")
	{
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/login", userHandler.UserLogin)
		user.PATCH("/forgotpassword", userHandler.ForgotPassword)
		user.Use(middleware.UserAuth)
		{
			user.POST("/logout", handler.UserLogout)
			userProfile := user.Group("/userprofile")
			{
				userProfile.GET("/", userHandler.ViewUserProfile)
				userProfile.PATCH("/mobile/edit", userHandler.UpdateMobile)
				userProfile.PATCH("/changepassword", userHandler.ChangePassword)
				address := userProfile.Group("/address")
				{
					address.POST("/add", userHandler.AddAddress)
					address.PATCH("/:address_id/edit", userHandler.UpdateAddress)
				}
			}
			cart := user.Group("/cart")
			{
				cart.GET("/", carrtHandler.ListCart)
				cart.POST("/:product_item_id/addtocart", carrtHandler.AddToCart)
				cart.DELETE("/:product_item_id/removefromcart", carrtHandler.RemoveFromCart)
			}
		}
		user.Use(middleware.UserIsBlocked)
	}
	admin := engine.Group("/admin")
	{
		admin.POST("/login", adminHandler.AdminLogin)
		admin.Use(middleware.AdminAuth)
		{
			admin.POST("/logout", adminHandler.AdminLogout)
			users := admin.Group("/users")
			{
				users.GET("/", adminHandler.ListAllUsers)
				users.GET("/:user_id", adminHandler.DisplayUser)
				users.PATCH("/:user_id/report", adminHandler.ReportUser)
			}
			category := admin.Group("/category")
			{
				category.POST("/create", productHandler.CreateCategory)
				category.PATCH("/update/:id", productHandler.UpdateCategory)
				category.DELETE("/delete/:id", productHandler.DeleteCategory)
				category.GET("/", productHandler.ListAllCategories)
				category.GET("/:id", productHandler.DisplayCategory)
			}
			brand := admin.Group("/brand")
			{
				brand.POST("/create", productHandler.CreateBrand)
				brand.PATCH("/update/:id", productHandler.UpdateBrand)
				brand.DELETE("/delete/:id", productHandler.DeleteBrand)
				brand.GET("/", productHandler.ListAllBrands)
				brand.GET("/:brand_id", productHandler.DisplayBrand)
			}
			product := admin.Group("/product")
			{
				product.POST("/create", productHandler.AddProduct)
				product.PATCH("/update/:product_id", productHandler.UpdateProduct)
				product.DELETE("/delete/:product_id", productHandler.DeleteProduct)
				product.GET("/", productHandler.ListAllProducts)
				product.GET("/:product_id", productHandler.DisplayProduct)
			}
			productItem := admin.Group("/productitem")
			{
				productItem.POST("/create", productHandler.AddProductItem)
				productItem.POST("/create/uploadimage/:productItem_id", productHandler.UploadImage)
				productItem.PATCH("/update/:productItem_id", productHandler.UpdateProductItem)
				productItem.DELETE("/delete/:productItem_id", productHandler.DeleteProductItem)
				productItem.GET("/", productHandler.ListAllProductItems)
				productItem.GET("/:productItem_id", productHandler.DisplayProductItem)
				productItem.DELETE("/:productItem_id/deleteimage", productHandler.DeleteImage)
			}

		}
	}
	superAdmin := engine.Group("/superadmin")
	{
		superAdmin.POST("/login", superadminHandler.SuperLogin)
		superAdmin.Use(middleware.SuperAdminAuth)
		{
			superAdmin.POST("/logout", superadminHandler.SuperLogout)
			admin := superAdmin.Group("/admin")
			{
				admin.POST("/create", superadminHandler.CreateAdmin)
				admin.GET("/", superadminHandler.ListAllAdmins)
				admin.GET("/:admin_id", superadminHandler.DisplayAdmin)
				admin.PATCH("/:admin_id/block", superadminHandler.BlockAdmin)
				admin.PATCH("/:admin_id/unblock", superadminHandler.UnBlockAdminManually)
			}
			user := superAdmin.Group("/user")
			{
				user.GET("/", adminHandler.ListAllUsers)
				user.GET("/:user_id", adminHandler.DisplayUser)
				user.PATCH("/:user_id/block", superadminHandler.BlockUser)
				user.PATCH("/:user_id/unblock", superadminHandler.UnBlockUserManually)
			}
		}
	}

	return &ServerHTTP{engine: engine}
}
func (sh *ServerHTTP) Start() {
	sh.engine.Run(":8080")
}
