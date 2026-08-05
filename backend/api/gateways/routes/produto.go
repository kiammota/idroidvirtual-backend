package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

// ProdutoRoutes registra o RouterGroup de produtos na aplicação Gin
func ProdutoRoutes(router *gin.RouterGroup, productController *controllers.ProductController) {
	productsGroup := router.Group("/products")
	{
		// GET /products - Lista produtos
		productsGroup.GET("", productController.ListAll)

		// POST /products - Cria um produto
		productsGroup.POST("", productController.Create)

		// GET /products/category/:category - Lista produtos por categoria
		productsGroup.GET("/category/:category", productController.GetByCategory)

		// POST /products/:id - Atualiza um produto completo
		productsGroup.POST("/:id", productController.Update)

		// DELETE /products/:id - Remove um produto
		productsGroup.DELETE("/:id", productController.Delete)
	}
}
