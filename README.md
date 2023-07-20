# boilerplate

boilerplate for using aiteung package. Please init your apps, and replace package import with the name iteung/... in the main.go, controller.go dan url.go

```sh
go mod init yourpackageurl
go mod tidy
```

## main

main.go file

```go
package main

import (
"log"

"github.com/gocroot/gocroot/config"

"github.com/aiteung/musik"
"github.com/gofiber/fiber/v2/middleware/cors"

"github.com/whatsauth/whatsauth"

"github.com/gocroot/gocroot/url"

"github.com/gofiber/fiber/v2"
)

func main() {
go whatsauth.RunHub()
site := fiber.New(config.Iteung)
site.Use(cors.New(config.Cors))
url.AuthSession(site)
url.Web(site)
log.Fatal(site.Listen(musik.Dangdut()))
}


```

## url

url.go file inside url folder, act as package

```go
package url

import (
"github.com/gocroot/gocroot/controller"
"github.com/gocroot/gocroot/middleware"
"github.com/gofiber/fiber/v2"
"github.com/gofiber/websocket/v2"
)

func Web(page *fiber.App) {
page.Use(middleware.Middleware)
page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
page.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth

page.Get("/", controller.Homepage) //ujicoba panggil package musik
page.Get("/presensi", controller.GetPresensiBulanIni)
}

func AuthSession(page *fiber.App) {
page.Post("/login", controller.LoginUser)
page.Post("/register", controller.RegisterUser)
page.Get("/logout", func(c *fiber.Ctx) error {
c.ClearCookie("token")
return nil
})
}


```

## controller

controller.go inside controller folder act as package

```go
package controller

import (
"context"
"errors"
"github.com/aiteung/musik"
"github.com/aiteung/presensi"
"github.com/gocroot/gocroot/config"
"github.com/gocroot/gocroot/models"
"github.com/gocroot/gocroot/pisato"
"github.com/gocroot/gocroot/utils"
"github.com/google/uuid"
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
var userLogin models.UserLogin
var userInfo whatsauth.LoginInfo

db := config.Ulbimongoconn.Collection("user")

if err := c.BodyParser(&userLogin); err != nil {
return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
}

if err := db.FindOne(context.TODO(), userLogin.UserName).Decode(userInfo); err == nil {
return utils.ErrorLogRes(http.StatusInternalServerError, err, "error", c)
}

if userLogin.UserName == "" {
return utils.ErrorLogRes(http.StatusBadRequest, errors.New("user not found"), "error", c)
}

if err := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(userInfo.Password)); err != nil {
return utils.ErrorLogRes(http.StatusInternalServerError, err, "error", c)
}

maker, err := pisato.NewPasetoMaker(config.PrivateKey)
if err != nil {
return utils.ErrorLogRes(http.StatusBadRequest, err, "error", c)
}

token, err := maker.CreateToken(userLogin.UserName, 4*time.Hour)
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

```

## Config Folder Package

cors.go

```go
package config

import (
"os"
"strings"

"github.com/gofiber/fiber/v2/middleware/cors"
)

var origins = []string{
"https://adorableproject.github.io",
"https://auth.ulbi.ac.id",
"https://sip.ulbi.ac.id",
"https://euis.ulbi.ac.id",
"https://home.ulbi.ac.id",
"https://alpha.ulbi.ac.id",
"https://dias.ulbi.ac.id",
"https://iteung.ulbi.ac.id",
"https://whatsauth.github.io",
" http://127.0.0.1:80",
}

var Internalhost string = os.Getenv("INTERNALHOST") + ":" + os.Getenv("PORT")

var Cors = cors.Config{
AllowOrigins:     strings.Join(origins[:], ","),
AllowHeaders:     "Origin",
ExposeHeaders:    "Content-Length",
AllowCredentials: true,
}
```

token.go

```go
package config

import "os"

var PublicKey string = os.Getenv("PUBLICKEY")
var PrivateKey string = os.Getenv("PRIVATEKEY")

```

db.go

```go
package config

import (
 "os"

 "github.com/aiteung/atdb"
 "github.com/whatsauth/whatsauth"
)

var IteungIPAddress string = os.Getenv("ITEUNGBEV1")

var MongoString string = os.Getenv("MONGOSTRING")

var MariaStringAkademik string = os.Getenv("MARIASTRINGAKADEMIK")

var DBUlbimariainfo = atdb.DBInfo{
 DBString: MariaStringAkademik,
 DBName:   "db_ulbi",
}

var Ulbimariaconn = atdb.MariaConnect(DBUlbimariainfo)

var Usertables = [4]whatsauth.LoginInfo{mhs, dosen, user, user1}

var mhs = whatsauth.LoginInfo{
 Userid:   "MhswID",
 Password: "Password",
 Phone:    "Telepon",
 Username: "Login",
 Uuid:     "simak_mst_mahasiswa",
 Login:    "2md5",
}

var dosen = whatsauth.LoginInfo{
 Userid:   "NIDN",
 Password: "Password",
 Phone:    "Handphone",
 Username: "Login",
 Uuid:     "simak_mst_dosen",
 Login:    "2md5",
}

var user = whatsauth.LoginInfo{
 Userid:   "user_id",
 Password: "user_password",
 Phone:    "phone",
 Username: "user_name",
 Uuid:     "simak_besan_users",
 Login:    "2md5",
}

var user1 = whatsauth.LoginInfo{
 Userid:   "user_id",
 Password: "user_password",
 Phone:    "user_phone",
 Username: "user_name",
 Uuid:     "besan_users",
 Login:    "2md5",
}

```

## handler

maker.go


```go
package handler

import "time"

type Maker interface {
CreateToken(username string, duration time.Duration) (string, error)

VerifyToken(token string) (*Payload, error)
}
```


payload.go


```go
package handler

import (
"errors"
"github.com/google/uuid"
"time"
)

type Payload struct {
Id        uuid.UUID `json:"id"`
UserName  string    `json:"user_name"`
IssuedAt  time.Time `json:"issued_at"`
ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
tokenId, err := uuid.NewRandom()

if err != nil {
return nil, err
}

payload := &Payload{
Id:        tokenId,
UserName:  username,
IssuedAt:  time.Now(),
ExpiredAt: time.Now().Add(duration),
}

return payload, nil
}

func (payload *Payload) Valid() error {
if time.Now().After(payload.ExpiredAt) {
return errors.New("token is expired")
}
return nil
}

```

## middleware

middleware.go


```go
package middleware

import (
"github.com/gocroot/gocroot/config"
"github.com/gocroot/gocroot/pisato"
"github.com/gofiber/fiber/v2"
"net/http"
)

func Middleware(c *fiber.Ctx) error {
token := c.Cookies("token")

if token == "" {
return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
"status": "Unauthorized",
})
}

maker, err := pisato.NewPasetoMaker(config.PrivateKey)
if err != nil {
return c.Status(http.StatusBadRequest).JSON(err)
}

_, err = maker.VerifyToken(token)
if err != nil {
return c.Status(http.StatusUnauthorized).JSON(err)
}

return c.Next()
}
```

## pisato

pisato.go


```go
package pisato

import (
"errors"
"github.com/aead/chacha20poly1305"
"github.com/gocroot/gocroot/handler"
"github.com/o1egl/paseto"
"time"
)

type PasetoMaker struct {
paseto *paseto.V2
key    []byte
}

func NewPasetoMaker(key string) (handler.Maker, error) {
if len(key) != chacha20poly1305.KeySize {
return nil, errors.New("error because symmetric doesn't contain 32 length")
}

maker := &PasetoMaker{
paseto: paseto.NewV2(),
key:    []byte(key),
}

return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
payload, err := handler.NewPayload(username, duration)
if err != nil {
return "", err
}

encrypt, err := maker.paseto.Encrypt(maker.key, payload, nil)
if err != nil {
return "", err
}

return encrypt, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*handler.Payload, error) {
payload := &handler.Payload{}

if err := maker.paseto.Decrypt(token, maker.key, payload, nil); err != nil {
return nil, err
}

if err := payload.Valid(); err != nil {
return nil, err
}

return payload, nil
}

```


# Endpoints

## Unprotected
1. `POST` /login :
Post with {"user_name":"{your user name}","user_pass":"{your password}"}
User must log in to this endpoint before access protected endpoints

2. `POST` /register :
Post with {"user_name":"{your user name}","user_pass":"{your password}","phone":"{your phone number}"}
User must register their account before access log in endpoint

3. `GET` /logout :
If user want to log out, they must get request in this endpoint

## Protected
1. `POST` /api/whatsauth/request
2. `GET` /ws/whatsauth/qr
3. `GET` /
4. `GET` /presensi