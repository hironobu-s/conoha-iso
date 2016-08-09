package webui

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"time"

	"github.com/hironobu-s/conoha-iso/command"
	"github.com/k0kubun/pp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v8"
)

type Template struct {
	templates *template.Template
}

var validate *validator.Validate
var ident *command.Identity

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	c.Logger().Debug("Render")

	var err error
	t.templates, err = template.ParseGlob("webui/template/*")
	if err != nil {
		return err
	}

	// Execute contents template
	buf := bytes.NewBufferString("")
	if err = t.templates.ExecuteTemplate(buf, name, data); err != nil {
		return err
	}

	// Execute layout template with contents.
	return t.templates.ExecuteTemplate(w, "layout", map[string]template.HTML{
		"Body": template.HTML(buf.String()),
	})
}

func RunServer(i *command.Identity) error {
	if err := i.Auth(); err != nil {
		return err
	}
	ident = i

	tpl := &Template{}

	// initialize validator
	config := &validator.Config{TagName: "test"}
	validate = validator.New(config)

	// initialize web framework
	e := echo.New()
	e.SetDebug(true)
	e.SetRenderer(tpl)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Assets
	e.Static("/static", "assets")

	// Routing
	e.GET("/", index)
	e.GET("/isos", isos)
	e.POST("/isos", insertIso)
	e.GET("/servers", servers)
	e.POST("/download", download)

	e.Run(standard.New(":12345"))
	return nil
}

func index(c echo.Context) error {
	params := map[string]string{}

	alert := popAlert(c)
	if alert != "" {
		params["alert"] = alert
	}
	notice := popNotice(c)
	if notice != "" {
		params["notice"] = notice
	}
	return c.Render(http.StatusOK, "top", params)
}

type IsoFormParams struct {
	Iso    string `validate:"required"`
	Server string `validate:"required"`
	Action string
}

func (p *IsoFormParams) FromFormParams(params map[string][]string) {
	var ok bool
	_, ok = params["server"]
	if ok {
		p.Server = params["server"][0]
	}

	_, ok = params["iso"]
	if ok {
		p.Iso = params["iso"][0]
	}
}

func insertIso(c echo.Context) error {
	params := IsoFormParams{}
	params.FromFormParams(c.FormParams())

	pp.Printf("%v\n", c.FormParams())
	pp.Printf("%v\n", params)

	isoId := c.Param("iso")
	serverId := c.Param("server")
	cp := command.NewCompute(ident)

	var err error
	server, err := cp.Server(serverId)
	if err != nil {
		c.Logger().Errorf("%v", err)
		return err
	}

	iso, err := cp.Iso(isoId)
	if err != nil {
		c.Logger().Errorf("%v", err)
		return err
	}

	if err = cp.Insert(server, iso); err != nil {
		c.Logger().Errorf("%v", err)
		return err
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

type DownloadFormParams struct {
	DownloadUrl string `validate:"required"`
}

func (p *DownloadFormParams) FromFormParams(params map[string][]string) {
	var ok bool
	_, ok = params["download_url"]
	if ok {
		p.DownloadUrl = params["download_url"][0]
	}
}

func download(c echo.Context) error {
	var err error
	params := &DownloadFormParams{}
	params.FromFormParams(c.FormParams())

	// validate
	config := &validator.Config{TagName: "validate"}
	v := validator.New(config)
	err = v.Struct(params)
	if err != nil {
		setAlert(c, err.Error())
	}

	// requrest downloading
	cp := command.NewCompute(ident)
	if err = cp.Download(params.DownloadUrl); err != nil {
		setAlert(c, err.Error())
	}

	if err == nil {
		setNotice(c, "Download request was accepted to succeed")
	}
	return c.Redirect(http.StatusMovedPermanently, "./")
}

func isos(c echo.Context) error {
	var err error

	cp := command.NewCompute(ident)
	isos, err := cp.Isos()
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, isos); err != nil {
		return err
	}
	return nil
}

func servers(c echo.Context) error {
	var err error

	cp := command.NewCompute(ident)
	servers, err := cp.Servers()
	if err != nil {
		return err
	}

	if err = c.JSON(http.StatusOK, servers); err != nil {
		return err
	}
	return nil
}

// ----- Flash Messages

func setAlert(c echo.Context, message string) {
	setFlash(c, "flash-alert", message)
}
func popAlert(c echo.Context) string {
	return popFlash(c, "flash-alert")
}
func setNotice(c echo.Context, message string) {
	setFlash(c, "flash-notice", message)
}
func popNotice(c echo.Context) string {
	return popFlash(c, "flash-notice")
}

func setFlash(c echo.Context, name string, message string) {
	cookie := &echo.Cookie{}
	cookie.SetName(name)
	cookie.SetValue(message)
	c.SetCookie(cookie)
}

func popFlash(c echo.Context, name string) (value string) {
	cc, err := c.Cookie(name)
	if err == nil {
		value = cc.Value()

		// remove flash cookie
		cookie := &echo.Cookie{}
		cookie.SetName(name)
		cookie.SetExpires(time.Now().Add(-1 * (100 * time.Second)))
		c.SetCookie(cookie)
	}

	return value
}
