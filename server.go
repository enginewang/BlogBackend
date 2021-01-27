package BlogBackend

import (
	"BlogBackend/controller/article"
	"BlogBackend/controller/auth"
	"BlogBackend/controller/comment"
	"BlogBackend/controller/upload"
	"BlogBackend/controller/user"
	"echo/middleware"
	"github.com/labstack/echo"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

type Server struct {
	Addr string
	e    *echo.Echo
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
		e:    echo.New(),
	}
}

func (s *Server) Init() (err error) {
	s.e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<h1>Welcome to Echo!</h1>
			<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
		`)
	})
	s.e.Static("/image", "/root/go/src/BlogBackend/upload")
	s.e.Static("/file", "/root/go/src/BlogBackend/file")
	//g := s.e.Group("")
	e := s.e
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	err = article.NewController().Initialize(*e)
	err = user.NewController().Initialize(*e)
	err = comment.NewController().Initialize(*e)
	err = auth.NewController().Initialize(*e)
	err = upload.NewController().Initialize(*e)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) StartServer ()  {
	//s.e.Logger.Fatal(s.e.Start(":1323"))
	s.e.Logger.Fatal(s.e.StartAutoTLS(":443"))
}
