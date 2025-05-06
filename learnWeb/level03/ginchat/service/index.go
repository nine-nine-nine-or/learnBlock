package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserByID godoc
// @Summary      获取gy信息
// @Description  通过用户ID获取详细信息
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "成功返回用户列表"
// @Router       /users [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "高洋真帅",
	})
}

// GetUserByID godoc
// @Summary      获取所有用户信息
// @Description  获取所有用户信息
// @Tags        用户模块
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.UserBasic  "成功返回用户列表"
// @Router       /users/list [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	log.Printf("查询结果: %+v", data) // 确认数据不为空

	if len(data) == 0 {
		c.JSON(200, gin.H{"message": "暂无数据"})
		return
	}
	c.JSON(200, gin.H{"message": data})
}

// CreateUser
// @Description  新增用户信息
// @Summary      新增用户信息
// @Tags        用户模块
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       name     formData string true "用户名"
// @Param       password formData string true "密码"
// @Param       phone    formData string true "手机号"
// @Success      200  {string}  json{“code”,"message"}  "成功返回用户列表"
// @Router       /users/create [post]
func CreateUser(c *gin.Context) {

	user := models.UserBasic{
		Name:     c.PostForm("name"),
		Password: c.PostForm("password"),
		Phone:    c.PostForm("phone"),
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	//创建用户时候判断是否已经创建过
	log.Printf("用户信息为: %+v", user)
	data := models.GetUserByName(user.Name)

	if data.Name != "" {
		c.JSON(400, gin.H{"message": "用户已存在"})
		return
	}
	//对存入数据库的密码加密
	password := utils.MakePassword(user.Password, salt)
	log.Printf("用户加密后的密码为: %v", password)
	user.Password = password
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{"message": "新增用户成功"})
}

// DeleteUser
// @Description  删除用户信息
// @Summary      删除用户信息
// @Tags        用户模块
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       id     path int true "用户名"
// @Success      200  {string}  json{“code”,"message"}  "成功返回用户列表"
// @Router       /users/delete/{id} [get]
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	models.DeleteUser(id)
	c.JSON(200, gin.H{"message": "删除用户成功"})
}

// UpdateUser
// @Description  修改用户信息
// @Summary      修改用户信息
// @Tags        用户模块
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       name     formData string true "用户名"
// @Param       password formData string true "密码"
// @Param       phone    formData string true "手机号"
// @Param       id    formData int true "id"
// @Success      200  {string}  json{“code”,"message"}  "成功返回"
// @Router       /users/update [post]
func UpdateUser(c *gin.Context) {

	user := models.UserBasic{
		Name:     c.PostForm("name"),
		Password: c.PostForm("password"),
		Phone:    c.PostForm("phone"),
	}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	log.Printf("用户信息为: %+v", user)
	models.UpdateUser(&user)
	c.JSON(200, gin.H{"message": "修改用户成功"})
}

// LoginIn
// @Description  登录到主界面
// @Summary      登录方法
// @Tags        通用模块
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       name     formData string true "用户名"
// @Param       password formData string true "密码"
// @Success      200  {string}  json{“code”,"message"}  "成功返回用户列表"
// @Router       /users/loginIn [post]
func LoginIn(c *gin.Context) {

	user := models.UserBasic{
		Name:     c.PostForm("name"),
		Password: c.PostForm("password"),
	}
	//通过数据库里的加密密码来比对
	data := models.GetUserByName(user.Name)
	flag := utils.DoPassword(user.Password, data.Salt, data.Password)
	if !flag {
		c.JSON(400, gin.H{"message": "用户密码不对，请重新输入"})
		return
	}
	c.JSON(200, gin.H{"message": "登录成功"})
}

// 防止伪造的跨域请求，设置为true就是允许所有的跨域请求
// 如果不设置，应该就不能调用了
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(c *gin.Context) {
	wr, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
	}
	defer func(wr *websocket.Conn) {
		err := wr.Close()
		if err != nil {
			log.Print(err)
		}
	}(wr)
	MsgHandle(wr, c)
}

func MsgHandle(wr *websocket.Conn, c *gin.Context) {

	subscribe, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		log.Println("subscribe:---------------", err)
		return
	}
	log.Println("subscribe:", subscribe)
	tem := time.Now().Format("2006-01-02 15:04:05")
	sprintf := fmt.Sprintf("[wr][%s]:%s", tem, subscribe)
	err = wr.WriteMessage(websocket.TextMessage, []byte(sprintf))
	if err != nil {
		log.Println("writeMessage:---------------", err)
		return
	}
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "upgrade failed"})
		return
	}

	// 假设前端传来的 channel 名称是 "chat_room_1"
	channel := "chat_room_1"

	// 开启一个 goroutine 来监听 Redis 消息并推送给客户端
	utils.ListenRedisChannel(c.Request.Context(), channel, conn)

	// 可选：处理从客户端发来的消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received from client: %s", p)

		// 回复客户端
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}
