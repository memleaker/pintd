package api

import (
	"net/http"
	"pintd/model"
	"pintd/plog"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndirectCfgNew(c *gin.Context) {
	cfg := model.IndirectConfig{}

	// parse JSON to structure.
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求格式",
		})
		return
	}

	//todo check values. 检查长度, IP, 端口

	// check config is or not exist.
	ok, err := model.RepeatConfig(&cfg)
	if ok {
		c.JSON(http.StatusConflict, gin.H{
			"success": 0,
			"msg":     "配置已经存在",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 0,
			"msg":     "服务器发生错误!",
		})

		plog.Println("model.RepeatConfig() : " + err.Error())
		return
	}

	// new config.
	if err = model.NewIndirectCfg(&cfg); err != nil {
		// create failed.
		plog.Println("model.NewIndirectCfg() : " + err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 0,
			"msg":     "服务器发生错误!",
		})
		return
	}

	// create success.
	c.JSON(http.StatusCreated, gin.H{
		"success": 1,
		"msg":     "新建配置成功",
	})
}

func IndirectCfgDel(c *gin.Context) {
	cfg := model.IndirectConfig{}

	// parse JSON to structure.
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求格式",
		})
		return
	}

	// delete config.
	ok, err := model.DelIndirectCfg(&cfg.Protocol, &cfg.ListenPort)
	if !ok {
		// delete failed.
		if err != nil {
			plog.Println("model.DelIndirectCfg() : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "服务器发生错误!",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "删除配置失败",
			})
		}

		return
	}

	// delete success.
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"msg":     "删除配置成功",
	})
}

func IndirectCfgShow(c *gin.Context) {
	// get url args.
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page <= 0 || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误的请求参数",
		})
		return
	}

	// get rows number.
	rows, err := model.GetIndirectTblRows()
	if err != nil {
		plog.Println("model.GetIndirectTblRows() : " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 0,
			"msg":     "服务器发生错误!",
		})
		return
	}

	// get config.
	cfgs, err := model.GetIndirectCfg(page, limit)
	if err != nil {
		plog.Println("model.GetIndirectCfg() : " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 0,
			"msg":     "服务器发生错误!",
		})
		return
	}

	// get config success.
	c.JSON(http.StatusOK, gin.H{
		"code":  0, // code:0, means no error.
		"count": rows,
		"data":  &cfgs,
	})
}

func IndirectCfgEdit(c *gin.Context) {
	cfg := model.IndirectConfig{}

	// get url args.
	field := c.Query("field")

	// parse JSON.
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求格式",
		})
		return
	}

	// because url arg 'field' is equal to cfg's tags.
	// so compare with tags and get value.
	name := ""
	tags := reflect.TypeOf(&cfg).Elem()
	ref := reflect.ValueOf(cfg)

	// get field name by tag name.
	for i := 0; i < tags.NumField(); i++ {
		if tags.Field(i).Tag.Get("json") == field {
			name = tags.Field(i).Name
			break
		}
	}

	// not found.
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求格式",
		})
		return
	}

	// found field, get value.
	val := ref.FieldByName(name).String()

	// todo check value is valid ?

	// update config.
	ok, err := model.UpdateIndirectCfg(&field, &val, &cfg)
	if !ok {
		// update failed.
		if err != nil {
			plog.Println("model.UpdateIndirectCfg() : " + err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "服务器发生错误!",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "编辑列失败",
			})
		}

		return
	}

	// update success.
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"msg":     "编辑列成功",
	})
}
