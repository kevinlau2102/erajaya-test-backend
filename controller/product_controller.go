package controller

import (
	"net/http"

	"erajaya-interview/dto"
	"erajaya-interview/service"
	"erajaya-interview/utils"

	"github.com/gin-gonic/gin"
)

type (
	ProductController interface {
		CreateProduct(ctx *gin.Context)
		GetAllProducts(ctx *gin.Context)
	}

	productController struct {
		productService service.ProductService
	}
)

func NewProductController(ps service.ProductService) ProductController {
	return &productController{
		productService: ps,
	}
}

func (c *productController) CreateProduct(ctx *gin.Context) {
	var product dto.ProductCreateRequest
	if err := ctx.ShouldBind(&product); err != nil {
		res := utils.BuildResponseFailed(http.StatusBadRequest, dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.productService.CreateProduct(ctx.Request.Context(), product)
	if err != nil {
		res := utils.BuildResponseFailed(http.StatusInternalServerError, dto.MESSAGE_FAILED_CREATE_PRODUCT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponseSuccess(http.StatusOK, dto.MESSAGE_SUCCESS_CREATE_PRODUCT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *productController) GetAllProducts(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(http.StatusBadRequest, dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.productService.GetAllProducts(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(http.StatusInternalServerError, dto.MESSAGE_FAILED_GET_LIST_PRODUCT, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	resp := utils.Response{
		Status:  http.StatusOK,
		Message: dto.MESSAGE_SUCCESS_GET_LIST_PRODUCT,
		Data:    result.Data,
		Meta:    result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, resp)
}
