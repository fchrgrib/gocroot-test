package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/aiteung/musik"
	"github.com/aiteung/presensi"
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/pisato"
	"github.com/gocroot/gocroot/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/whatsauth/whatsauth"
)

func Homepage(c *fiber.Ctx) error {
	ipaddr := musik.GetIPaddress()
	return c.JSON(ipaddr)
}

func GetPresensiBulanIni(c *fiber.Ctx) error {
	ps := presensi.GetPresensiCurrentMonth(config.Ulbimongoconn)
	return c.JSON(ps)
}

func WsWhatsAuthQR(c *websocket.Conn) {
	whatsauth.RunSocket(c, config.PublicKey, config.Usertables[:], config.Ulbimariaconn)
}

func PostWhatsAuthRequest(c *fiber.Ctx) error {
	if string(c.Request().Host()) == config.Internalhost {
		var req whatsauth.WhatsauthRequest
		err := c.BodyParser(&req)
		if err != nil {
			return err
		}
		ntfbtn := whatsauth.RunModule(req, config.PrivateKey, config.Usertables[:], config.Ulbimariaconn)
		return c.JSON(ntfbtn)
	} else {
		var ws whatsauth.WhatsauthStatus
		ws.Status = string(c.Request().Host())
		return c.JSON(ws)
	}
}

func LoginUser(c *fiber.Ctx) error {
	var userLogin whatsauth.LoginInfo
	var userInfo whatsauth.LoginInfo

	db := config.Ulbimongoconn.Collection("user")

	if err := c.BodyParser(&userLogin); err != nil {
		return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
	}

	if err := db.FindOne(context.TODO(), bson.M{"user_name": userLogin.Username}).Decode(&userInfo); err != nil {
		return utils.ErrorLogRes(http.StatusInternalServerError, err, "error", c)
	}

	if userInfo.Username == "" {
		return utils.ErrorLogRes(http.StatusBadRequest, errors.New("user not found"), "error", c)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(userLogin.Password)); err != nil {
		fmt.Println("2")
		return utils.ErrorLogRes(http.StatusInternalServerError, err, "error", c)
	}

	maker, err := pisato.NewPasetoMaker(config.PrivateKey)
	if err != nil {
		return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
	}

	token, err := maker.CreateToken(userLogin.Username, 4*time.Hour)
	if err != nil {
		return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(4 * time.Hour),
		Secure:   false,
		HTTPOnly: true,
	})
	return utils.ErrorLogRes(http.StatusOK, nil, "ok", c)
}

func RegisterUser(c *fiber.Ctx) error {
	var userInfo whatsauth.LoginInfo

	db := config.Ulbimongoconn.Collection("user")

	if err := c.BodyParser(&userInfo); err != nil {
		return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
	}

	var checkUserName whatsauth.LoginInfo
	if _ = db.FindOne(context.TODO(), userInfo.Username).Decode(&checkUserName); checkUserName.Username != "" {
		return utils.ErrorLogRes(http.StatusBadRequest, errors.New("user already exist"), "error", c)
	}

	userInfo.Uuid = uuid.New().String()
	userInfo.Userid = uuid.New().String()
	pass, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	userInfo.Password = string(pass)

	if _, err := db.InsertOne(context.TODO(), userInfo); err != nil {
		return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
	}
	return utils.ErrorLogRes(http.StatusOK, nil, "ok", c)
}
