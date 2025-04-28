package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"time"
)

type StructA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
	var b StructB
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
}

func GetDataC(c *gin.Context) {
	var b StructC
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})
}

func GetDataD(c *gin.Context) {
	var b StructD
	c.Bind(&b)
	c.JSON(200, gin.H{
		"x": b.NestedAnonyStruct,
		"d": b.FieldD,
	})
}

type City struct {
	Name   string
	Area   float64
	Person string
}

func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func main() {
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"UnixToTime": UnixToTime,
	})

	router.LoadHTMLGlob("templates/**/*")
	router.GET("/getb", GetDataB)
	router.GET("/getc", GetDataC)
	router.GET("/getd", GetDataD)
	router.GET("/gson01/:name", func(context *gin.Context) {
		context.JSON(200, map[string]interface{}{
			"name": context.Param("name"),
			"age":  20,
		})
	})

	router.GET("/gson01", func(context *gin.Context) {
		context.JSON(200, map[string]interface{}{
			"name": context.Query("name"),
			"age":  20,
		})
	})
	router.GET("/gson02", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"name": context.Query("name"),
		})
	})

	router.GET("/city", func(context *gin.Context) {
		city := &City{
			Name:   "信阳",
			Area:   100,
			Person: "999",
		}
		context.HTML(http.StatusOK, "index.html", gin.H{
			"title": "地理页面",
			"city":  city,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"title": "gyzs",
			"msg":   "我是gy",
			"score": 100,
			"hobby": []string{"吃饭", "睡觉", "学习", "上班"},
			"newsList": []interface{}{
				&City{
					Name:   "gy",
					Area:   100,
					Person: "999",
				},
				&City{
					Name:   "cq",
					Area:   100,
					Person: "666",
				},
			},
			"testSlice": []string{},
			"news": &City{
				Name:   "zs",
				Area:   100,
				Person: "888",
			},
			"date": 1629423555,
		})
	})

	router.Run()
}

//func main() {
//	router := gin.Default()
//	//r.GET("/", func(c *gin.Context) {
//	//	c.JSON(200, gin.H{
//	//		"message": "Hello, Gin!",
//	//	})
//	//})
//	//
//	//r.Run(":8080") // 默认监听 :8080
//	router.GET("/someJSON", func(c *gin.Context) {
//		data := map[string]interface{}{
//			"lang": "GO语言",
//			"tag":  "<br>",
//		}
//
//		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
//		c.AsciiJSON(http.StatusOK, data)
//	})
//
//	// 监听并在 0.0.0.0:8080 上启动服务
//	router.Run(":8080")
//}
