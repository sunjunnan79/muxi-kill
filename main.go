package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	g := gin.Default()
	dao := NewDAO(InitDB())
	cache := NewCache()
	item, err := dao.GetItemByID(1)
	if err == gorm.ErrRecordNotFound {
		err := dao.SaveItem(&Item{
			Model: gorm.Model{ID: 1},
			Name:  "muxi",
			Num:   10,
		})
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	if item.Num <= 0 {
		item.Num = 10
		err := dao.SaveItem(item)
		if err != nil {
			return
		}
	}

	g.POST("/kill", func(c *gin.Context) {

		var req struct {
			Id uint `json:"id" binding:"required"`
		}

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(422, gin.H{
				"msg": "请求参数错误,请重试",
			})
			return
		}

		//获取分布式锁
		m, err := cache.LockResource(fmt.Sprintf("item:%d", req.Id))
		if err != nil {
			c.JSON(429, gin.H{
				"msg": "商品太火爆请稍后再试",
			})
			return
		}
		defer func() {
			//关闭分布式锁
			err = cache.UnlockResource(m)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}()

		item, err := dao.GetItemByID(req.Id)
		if err != nil {
			c.JSON(404, gin.H{
				"msg": "不存在的商品",
			})
			return
		}

		if item.Num <= 0 {
			c.JSON(200, gin.H{
				"msg": "商品已售罄",
			})
			return
		}

		item.Num = item.Num - 1

		err = dao.SaveItem(item)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": "系统发生错误",
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "秒杀下单成功",
		})

	})
	g.Run()

}
