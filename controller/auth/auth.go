package auth

import (
	"BlogBackend/db"
	"BlogBackend/model"
	_ "crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("secret"),
})

const BaseURL = "/api/auth"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}
func (c *Controller) Initialize(e echo.Echo) (err error) {
	e.POST(BaseURL+"/adminLogin", adminLogin)
	//e.POST(BaseURL + "/userRegisterPre", userRegisterPre)
	//e.POST(BaseURL + "/userRegister", userRegister)
	e.POST(BaseURL+"/userLogin", Login)



	return nil
}

func adminLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var result model.User
	err := collection.Find(bson.M{"username": username}).One(&result)
	fmt.Println(result)
	if result.Role != 0 {
		return c.String(http.StatusUnauthorized, "该用户没有管理员权限！")
	}
	if err != nil {
		return c.String(http.StatusUnauthorized, "用户名不存在!")
	}
	if password != result.Password {
		return c.String(http.StatusUnauthorized, "密码不正确!")
	} else {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}

/*
func userRegisterPre(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var newUser model.User
	err = c.Bind(&newUser)
	if err != nil{
		return err
	}
	return c.JSON(http.StatusOK, newArticleTag)
}
*/

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var result model.User
	err := collection.Find(bson.M{"username": username}).One(&result)
	fmt.Println(result)
	if err != nil {
		return c.String(http.StatusUnauthorized, "用户名不存在!")
	}
	if password != result.Password {
		return c.String(http.StatusUnauthorized, "密码不正确!")
	} else {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["admin"] = false
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}




