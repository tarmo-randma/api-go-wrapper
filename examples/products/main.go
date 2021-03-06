package main

import (
	"context"
	"fmt"
	"github.com/erply/api-go-wrapper/internal/common"
	"github.com/erply/api-go-wrapper/pkg/api"
	sharedCommon "github.com/erply/api-go-wrapper/pkg/api/common"
	"github.com/erply/api-go-wrapper/pkg/api/products"
	"time"
)

func main() {
	apiClient, err := api.BuildClient()
	common.Die(err)

	resp, err := DeleteProductBulk(apiClient)
	common.Die(err)
	fmt.Printf("DeleteProductBulk:\n%+v\n", resp)

	err = DeleteProduct(apiClient)
	common.Die(err)

	saveProd, err := SaveProduct(apiClient)
	common.Die(err)
	fmt.Printf("SaveProduct:\n%+v\n", saveProd)

	saveProds, err := SaveProductsBulk(apiClient)
	common.Die(err)
	fmt.Printf("SaveProductsBulk:\n%+v\n", saveProds)

	prodGroups, err := GetProductGroups(apiClient)
	common.Die(err)
	fmt.Printf("GetProductGroups:\n%+v\n", prodGroups)

	GetProductsBulk(apiClient)

	prods, err := GetProductsInParallel(apiClient)
	common.Die(err)
	fmt.Printf("GetProductsInParallel:\n%+v\n", prods)

	prodStock, err := GetProductStock(apiClient)
	common.Die(err)
	fmt.Printf("GetProductStock:\n%+v\n", prodStock)

	prodStockFile, err := GetProductStockFile(apiClient)
	common.Die(err)
	fmt.Printf("GetProductStockFile:\n%+v\n", prodStockFile)

	prodStockFileBulk, err := GetProductStockFileBulk(apiClient)
	common.Die(err)
	fmt.Printf("GetProductStockFileBulk:\n%+v\n", prodStockFileBulk)

	res, err := SaveAssortment(apiClient)
	common.Die(err)
	fmt.Printf("SaveAssortment:\n%+v\n", res)

	resBulk, err := SaveAssortmentBulk(apiClient)
	common.Die(err)
	fmt.Printf("SaveAssortmentBulk:\n%+v\n", resBulk)

	AddAssortmentProducts(apiClient)

	AddAssortmentProductsBulk(apiClient)

	EditAssortmentProducts(apiClient)

	EditAssortmentProductsBulk(apiClient)

	RemoveAssortmentProducts(apiClient)

	RemoveAssortmentProductsBulk(apiClient)

	SaveProductCategory(apiClient)

	SaveProductCategoryBulk(apiClient)

	SaveBrand(apiClient)

	SaveBrandBulk(apiClient)

	SaveProductPriorityGroup(apiClient)

	SaveProductPriorityGroupBulk(apiClient)

	SaveProductGroup(apiClient)

	SaveProductGroupBulk(apiClient)

	DeleteProductGroup(apiClient)

	DeleteProductGroupBulk(apiClient)

	GetProductPriorityGroupBulk(apiClient)

	GetProductPriorityGroupsInParallel(apiClient)

	GetProductGroupBulk(apiClient)

	GetProductGroupsInParallel(apiClient)

	GetProductCategoriesBulk(apiClient)

	GetProductCategoriesInParallel(apiClient)
}

func GetProductsBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	bulkFilters := []map[string]interface{}{
		{
			"recordsOnPage":            "10",
			"pageNo":                   "1",
			"getAllLanguages":          "1",
			"getFIFOCost":              "1",
			"getMatrixVariations":      "1",
			"getPackageInfo":           "1",
			"getPriceCalculationSteps": 1,
			"getPackagingMaterials":    "1",
			"getParameters":            "1",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bulkResp, err := prodCli.GetProductsBulk(ctx, bulkFilters, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(bulkResp))
}

func GetProductsInParallel(cl *api.Client) ([]products.Product, error) {
	productsDataProvider := products.NewListingDataProvider(cl.ProductManager)

	lister := sharedCommon.NewLister(
		sharedCommon.ListingSettings{
			MaxRequestsCountPerSecond: 5,
			StreamBufferLength:        10,
			MaxItemsPerRequest:        300,
			MaxFetchersCount:          10,
		},
		productsDataProvider,
		func(sleepTime time.Duration) {
			time.Sleep(sleepTime)
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	prodsChan := lister.Get(ctx, map[string]interface{}{
		"changedSince": time.Date(2021, 2, 15, 0, 0, 0, 0, time.UTC).Unix(),
	})

	prods := make([]products.Product, 0)
	for prod := range prodsChan {
		if prod.Err != nil {
			return prods, prod.Err
		}
		prods = append(prods, prod.Payload.(products.Product))
	}

	return prods, nil
}

func GetProductGroups(cl *api.Client) ([]products.ProductGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	groups, err := cl.ProductManager.GetProductGroups(ctx, map[string]string{
		"productGroupID": "2",
	})

	return groups, err
}

func GetProductStock(cl *api.Client) ([]products.GetProductStock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	productStock, err := cl.ProductManager.GetProductStock(ctx, map[string]string{
		"warehouseID": "1",
	})

	return productStock, err
}

func GetProductStockFile(cl *api.Client) ([]products.GetProductStockFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	productStockFile, err := cl.ProductManager.GetProductStockFile(ctx, map[string]string{
		"warehouseID": "1",
	})

	return productStockFile, err
}

func GetProductStockFileBulk(cl *api.Client) (stockFiles []products.GetProductStockFile, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	bulkResp, err := cl.ProductManager.GetProductStockFileBulk(ctx, []map[string]interface{}{
		{
			"warehouseID": "1",
		},
		{
			"warehouseID": "2",
		},
	}, map[string]string{})
	if err != nil {
		return
	}

	for _, bulkItem := range bulkResp.BulkItems {
		stockFiles = append(stockFiles, bulkItem.GetProductStockFiles...)
	}

	return stockFiles, err
}

func SaveProductsBulk(cl *api.Client) (saveProdResult []products.SaveProductResult, err error) {
	prodCli := cl.ProductManager

	bulkFilters := []map[string]interface{}{
		{
			"groupID": "4",
			"code":    "123",
		},
		{
			"groupID": "4",
			"code":    "124",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bulkResp, err := prodCli.SaveProductBulk(ctx, bulkFilters, map[string]string{})
	if err != nil {
		return
	}

	for _, bulkItem := range bulkResp.BulkItems {
		saveProdResult = append(saveProdResult, bulkItem.Products...)
	}

	return
}

func SaveProduct(cl *api.Client) (saveProdResult products.SaveProductResult, err error) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"groupID": "4",
		"code":    "127",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	saveProdResult, err = prodCli.SaveProduct(ctx, filter)
	if err != nil {
		return
	}

	return
}

func DeleteProduct(cl *api.Client) (err error) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"productID": "85656",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = prodCli.DeleteProduct(ctx, filter)
	if err != nil {
		return
	}

	return
}

func DeleteProductBulk(cl *api.Client) (bulkResp products.DeleteProductResponseBulk, err error) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"productID": "85654",
		},
		{
			"productID": "85655",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bulkResp, err = prodCli.DeleteProductBulk(ctx, filter, map[string]string{})
	if err != nil {
		return
	}

	return
}

func SaveAssortment(cl *api.Client) (res products.SaveAssortmentResult, err error) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"name": "some assortment",
		"code": "123",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err = prodCli.SaveAssortment(ctx, filter)
	if err != nil {
		return
	}

	return
}

func SaveAssortmentBulk(cl *api.Client) (res products.SaveAssortmentResponseBulk, err error) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"name":            "onetwothree",
			"code":            "126",
			"attributeName1":  "one",
			"attributeType1":  "string",
			"attributeValue1": "mine",
		},
		{
			"name": "onefour",
			"code": "127",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err = prodCli.SaveAssortmentBulk(ctx, filter, map[string]string{})
	if err != nil {
		return
	}

	return
}

func AddAssortmentProducts(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"productIDs":   "1",
		"assortmentID": "4",
		"status":       "ACTIVE",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.AddAssortmentProducts(ctx, filter)
	common.Die(err)

	fmt.Printf("AddAssortmentProducts:\n%+v\n", res)
}

func AddAssortmentProductsBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"productIDs":   "1",
			"assortmentID": "4",
			"status":       "ACTIVE",
		},
		{
			"productIDs":   "2",
			"assortmentID": "4",
			"status":       "ACTIVE",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.AddAssortmentProductsBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Printf("AddAssortmentProductsBulk:\n%+v\n", res)
}

func EditAssortmentProducts(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"productIDs":   "1",
		"assortmentID": "4",
		"status":       "NOT_FOR_SALE",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.EditAssortmentProducts(ctx, filter)
	common.Die(err)

	fmt.Printf("EditAssortmentProducts:\n%+v\n", res)
}

func EditAssortmentProductsBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"productIDs":   "1",
			"assortmentID": "4",
			"status":       "NO_LONGER_ORDERED",
		},
		{
			"productIDs":   "2",
			"assortmentID": "4",
			"status":       "NO_LONGER_ORDERED",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.EditAssortmentProductsBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Printf("EditAssortmentProductsBulk:\n%+v\n", res)
}

func RemoveAssortmentProducts(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"productIDs":   "1",
		"assortmentID": "4",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.RemoveAssortmentProducts(ctx, filter)
	common.Die(err)

	fmt.Printf("RemoveAssortmentProducts:\n%+v\n", res)
}

func RemoveAssortmentProductsBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"productIDs":   "2",
			"assortmentID": "4",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.RemoveAssortmentProductsBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Printf("RemoveAssortmentProductsBulk:\n%+v\n", res)
}

func SaveProductCategory(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"name": "some prod category",
		"code": "123",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveProductCategory(ctx, filter)
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveProductCategoryBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"name": "onetwothree",
		},
		{
			"name": "onefour",
		},
		{
			"name": "twor",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveProductCategoryBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveBrand(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"name": "some brand",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveBrand(ctx, filter)
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveBrandBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"name": "onetwothree",
		},
		{
			"name": "onefour",
		},
		{
			"name": "twor",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveBrandBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveProductPriorityGroup(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"name": "some prio group",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveProductPriorityGroup(ctx, filter)
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveProductPriorityGroupBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"name": "onetwothree",
		},
		{
			"name": "onefour",
		},
		{
			"name": "twor",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveProductPriorityGroupBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func GetProductPriorityGroupBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"recordsOnPage": 5,
			"pageNo":        1,
		},
		{
			"recordsOnPage": 5,
			"pageNo":        2,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.GetProductPriorityGroupBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveProductGroup(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"name": "some product group",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveProductGroup(ctx, filter)
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func SaveProductGroupBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"name": "onetwothree",
		},
		{
			"name": "onefour",
		},
		{
			"name": "twor",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.SaveProductGroupBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func DeleteProductGroup(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := map[string]string{
		"productGroupID": "105",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := prodCli.DeleteProductGroup(ctx, filter)
	common.Die(err)
}

func DeleteProductGroupBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"productGroupID": "106",
		},
		{
			"productGroupID": "107",
		},
		{
			"productGroupID": "108",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bulkResp, err := prodCli.DeleteProductGroupBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(bulkResp))
}

func GetProductPriorityGroupsInParallel(cl *api.Client) {
	dataProvider := products.NewPrioGroupListingDataProvider(cl.ProductManager)

	lister := sharedCommon.NewLister(
		sharedCommon.ListingSettings{
			MaxRequestsCountPerSecond: 5,
			StreamBufferLength:        10,
			MaxItemsPerRequest:        300,
			MaxFetchersCount:          10,
		},
		dataProvider,
		func(sleepTime time.Duration) {
			time.Sleep(sleepTime)
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	prioGroupChan := lister.Get(ctx, map[string]interface{}{
		"recordsOnPage": 10,
		"pageNo":        1,
	})

	prioGroups := make([]products.ProductPriorityGroup, 0)
	for prodPrioGroup := range prioGroupChan {
		if prodPrioGroup.Err != nil {
			common.Die(prodPrioGroup.Err)
		}
		prioGroups = append(prioGroups, prodPrioGroup.Payload.(products.ProductPriorityGroup))
	}

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(prioGroups))
}

func GetProductGroupBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"recordsOnPage": 5,
			"pageNo":        3,
		},
		{
			"recordsOnPage": 5,
			"pageNo":        4,
		},
		{
			"recordsOnPage": 5,
			"pageNo":        5,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.GetProductGroupsBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func GetProductGroupsInParallel(cl *api.Client) {
	dataProvider := products.NewProductGroupsListingDataProvider(cl.ProductManager)

	lister := sharedCommon.NewLister(
		sharedCommon.ListingSettings{
			MaxRequestsCountPerSecond: 5,
			StreamBufferLength:        10,
			MaxItemsPerRequest:        300,
			MaxFetchersCount:          10,
		},
		dataProvider,
		func(sleepTime time.Duration) {
			time.Sleep(sleepTime)
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	prioGroupChan := lister.Get(ctx, map[string]interface{}{
		"recordsOnPage": 10,
		"pageNo":        1,
	})

	groups := make([]products.ProductGroup, 0)
	for prodGroup := range prioGroupChan {
		if prodGroup.Err != nil {
			common.Die(prodGroup.Err)
		}
		groups = append(groups, prodGroup.Payload.(products.ProductGroup))
	}

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(groups))
}

func GetProductCategoriesBulk(cl *api.Client) {
	prodCli := cl.ProductManager

	filter := []map[string]interface{}{
		{
			"recordsOnPage": 5,
			"pageNo":        1,
		},
		{
			"recordsOnPage": 5,
			"pageNo":        2,
		},
		{
			"recordsOnPage": 5,
			"pageNo":        3,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := prodCli.GetProductCategoriesBulk(ctx, filter, map[string]string{})
	common.Die(err)

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(res))
}

func GetProductCategoriesInParallel(cl *api.Client) {
	dataProvider := products.NewProductCategoriesListingDataProvider(cl.ProductManager)

	lister := sharedCommon.NewLister(
		sharedCommon.ListingSettings{
			MaxRequestsCountPerSecond: 5,
			StreamBufferLength:        10,
			MaxItemsPerRequest:        300,
			MaxFetchersCount:          10,
		},
		dataProvider,
		func(sleepTime time.Duration) {
			time.Sleep(sleepTime)
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	prioGroupChan := lister.Get(ctx, map[string]interface{}{
		"recordsOnPage": 10,
		"pageNo":        1,
	})

	categories := make([]products.ProductCategory, 0)
	for prodGroup := range prioGroupChan {
		if prodGroup.Err != nil {
			common.Die(prodGroup.Err)
		}
		categories = append(categories, prodGroup.Payload.(products.ProductCategory))
	}

	fmt.Println(common.ConvertSourceToJsonStrIfPossible(categories))
}
