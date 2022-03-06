package shop

import (
	"fmt"
	"login_register_demo/config"
	"strconv"

	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
)

type GoodsCategory struct {
	CategoryId          int    `xorm:"not null pk comment('种类id') INT(11)"`
	CategoryName        string `xorm:"not null comment('种类名称') VARCHAR(255)"`
	CategoryIcon        string `xorm:"not null comment('种类图标') VARCHAR(255)"`
	CategoryDescription string `xorm:"not null comment('种类描述') TEXT"`
	CategoryShowStatus  int    `xorm:"not null comment('是否展示该种类') TINYINT(1)"`
	ProductCount        int    `xorm:"not null comment('该种类下产品数量') INT(11)"`
	ShopId              int    `xorm:"not null comment('店铺id') index INT(11)"`
	ParentId            int    `xorm:"not null comment('父分类id值，0表示最高级分类') INT(11)"`
}

func GetShopCategoryName(c *gin.Context) {
	var shopCategoryName []interface{}
	var shopCategory []GoodsCategory
	set := mapset.NewSet()
	shopId := c.Request.URL.Query().Get("shop_id")
	err := config.Engine.Table("goods_category").Where("shop_id=?", shopId).Find(&shopCategory)
	if err != nil {
		panic(err.Error())
	} else {
		for i := 0; i < len(shopCategory); i++ {
			set.Add(shopCategory[i].CategoryName)
		}
		for val := range set.Iterator().C {
			shopCategoryName = append(shopCategoryName, val)
		}
		c.JSON(200, gin.H{
			"shop_categories": shopCategoryName,
		})
	}

}
func InsertShopCategory(c *gin.Context) {
	shopId, err := strconv.Atoi(c.Request.PostForm.Get("shop_id"))
	if err != nil {
		panic(err.Error())
	}
	categoryId, err := strconv.Atoi(c.Request.PostForm.Get("category_id"))
	if err != nil {
		panic(err.Error())
	}
	categoryName := c.Request.PostForm.Get("category_name")
	categoryIcon := c.Request.PostForm.Get("category_icon")
	categoryDescription := c.Request.PostForm.Get("category_description")
	categoryShowStatus, err := strconv.Atoi(c.Request.PostForm.Get("category_show_status"))
	if err != nil {
		panic(err.Error())
	}
	productCount, err := strconv.Atoi(c.Request.PostForm.Get("product_count"))
	if err != nil {
		panic(err.Error())
	}
	parentId, err := strconv.Atoi(c.Request.PostForm.Get("parent_id"))
	if err != nil {
		panic(err.Error())
	}
	shopCategory := new(GoodsCategory)
	result, err := config.Engine.Where("category_name=?", categoryName).Get(shopCategory)
	fmt.Println(result)
	if err != nil {
		fmt.Println(err)
	}
	if shopId != shopCategory.ShopId {
		// 无此种分类
		shopCategory.CategoryDescription = categoryDescription
		shopCategory.CategoryIcon = categoryIcon
		shopCategory.CategoryId = categoryId
		shopCategory.CategoryName = categoryName
		shopCategory.CategoryShowStatus = categoryShowStatus
		shopCategory.ShopId = shopId
		shopCategory.ProductCount = productCount
		shopCategory.ParentId = parentId
		affected, err := config.Engine.Insert(shopCategory)
		if err != nil {
			fmt.Println(err)
		}
		if affected != 1 {
			c.JSON(200, gin.H{
				"success": false,
			})
		} else {
			c.JSON(200, gin.H{
				"success":       true,
				"shop_id":       shopId,
				"category_name": categoryName,
			})
		}
	} else {
		fmt.Println("Already has one exsit item!")
		c.JSON(200, gin.H{
			"code":    400,
			"success": false,
		})
	}
}
