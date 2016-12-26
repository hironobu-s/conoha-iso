package webui

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hironobu-s/conoha-iso/command"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v8"
)

func RunServer(address string, ident *command.Identity) (err error) {
	tpl := &Template{}

	identHandler := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("ident", ident)
			return next(c)
		}
	}

	// initialize web framework
	e := echo.New()
	e.Renderer = tpl

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(identHandler)

	// Assets
	e.Static("/static", "assets")

	// Routing
	e.GET("/", index)
	e.POST("/download", download)
	e.POST("/insert", insert)
	e.POST("/eject", eject)
	e.GET("/isos", isos)
	e.GET("/servers", servers)

	// parse listen address
	pair := strings.Split(address, ":")
	if len(pair) != 2 {
		return fmt.Errorf("Invalid listen address[%s].", address)
	}

	ip := net.ParseIP(pair[0])
	_, err = strconv.Atoi(pair[1])
	if ip == nil || err != nil {
		return fmt.Errorf("Invalid listen address[%s].", address)
	}

	e.Logger.Printf("Running on http://%s/", address)
	e.Logger.Fatal(e.Start(address))
	return nil
}

// -----------------------------------------------------------------

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
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

// ---------------------------------------------

type IndexTemplateParams struct {
	Error         error
	CheckedServer string
	CheckedIso    string
	DownloadUrl   string
	Notice        string
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", IndexTemplateParams{
		Notice: popNotice(c),
	})
}

func formError(c echo.Context, validationErrors error) error {
	var err error

	ve, ok := validationErrors.(validator.ValidationErrors)
	if ok {
		b := bytes.NewBufferString("")
		for _, fe := range ve {
			if fe.Tag == "url" {
				b.WriteString(fmt.Sprintf("Invalid URL format[%s]. ", fe.Name))
			} else if fe.Tag == "required" {
				b.WriteString(fmt.Sprintf("Required[%s]. ", fe.Name))
			} else {
				b.WriteString(fmt.Sprintf("Unknown error[%s]. ", fe.Name))
			}

		}
		err = fmt.Errorf("%s", b.String())

	} else {
		err = validationErrors
	}

	return c.Render(http.StatusOK, "index", IndexTemplateParams{
		Error:         err,
		CheckedServer: fmt.Sprintf(", checked:'%s'", c.FormValue("server")),
		CheckedIso:    fmt.Sprintf(", checked:'%s'", c.FormValue("iso")),
		DownloadUrl:   c.FormValue("download_url"),
		Notice:        popNotice(c),
	})
}

func eject(c echo.Context) error {
	params := struct {
		ServerId string `validate:"required,uuid"`
	}{
		ServerId: c.FormValue("server"),
	}

	config := &validator.Config{TagName: "validate"}
	v := validator.New(config)
	err := v.Struct(params)
	if err != nil {
		return formError(c, err)
	}

	ident, ok := c.Get("ident").(*command.Identity)
	if !ok {
		return formError(c, fmt.Errorf(`Can not convert 'ident' to "*command.Identity".`))
	}

	cp := command.NewCompute(ident)
	server, err := cp.Server(params.ServerId)
	if err != nil {
		return formError(c, err)
	}

	if err = cp.Eject(server); err != nil {
		return formError(c, err)
	}

	setNotice(c, fmt.Sprintf("ISO image has been ejected successfully from VPS. [%s]", server.Metadata.InstanceNameTag))

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func insert(c echo.Context) error {
	params := struct {
		IsoId    string `validate:"required,uuid"`
		ServerId string `validate:"required,uuid"`
	}{
		IsoId:    c.FormValue("iso"),
		ServerId: c.FormValue("server"),
	}

	config := &validator.Config{TagName: "validate"}
	v := validator.New(config)
	err := v.Struct(params)
	if err != nil {
		return formError(c, err)
	}

	ident, ok := c.Get("ident").(*command.Identity)
	if !ok {
		return formError(c, fmt.Errorf(`Can not convert 'ident' to "*command.Identity".`))
	}

	cp := command.NewCompute(ident)
	server, err := cp.Server(params.ServerId)
	if err != nil {
		return formError(c, err)
	}

	iso, err := cp.Iso(params.IsoId)
	if err != nil {
		return formError(c, err)
	}

	if err = cp.Insert(server, iso); err != nil {
		return formError(c, err)
	}

	setNotice(c, fmt.Sprintf(`An ISO image has been inserted to VPS. [%s => %s]`, iso.Name, server.Metadata.InstanceNameTag))

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func download(c echo.Context) error {
	var err error
	params := struct {
		DownloadUrl string `validate:"required,url"`
	}{
		DownloadUrl: c.FormValue("download_url"),
	}

	// validate
	config := &validator.Config{TagName: "validate"}
	v := validator.New(config)
	err = v.Struct(params)
	if err != nil {
		return formError(c, err)
	}

	// requrest downloading
	ident, ok := c.Get("ident").(*command.Identity)
	if !ok {
		return formError(c, fmt.Errorf(`Can not convert 'ident' to "*command.Identity".`))
	}

	cp := command.NewCompute(ident)
	if err = cp.Download(params.DownloadUrl); err != nil {
		return formError(c, err)
	}

	setNotice(c, "Download request was accepted successfully.")
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func isos(c echo.Context) error {
	var err error

	ident, ok := c.Get("ident").(*command.Identity)
	if !ok {
		return fmt.Errorf(`Can not convert 'ident' to "*command.Identity".`)
	}

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

	ident, ok := c.Get("ident").(*command.Identity)
	if !ok {
		return fmt.Errorf(`Can not convert 'ident' to "*command.Identity".`)
	}

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
func setNotice(c echo.Context, message string) {
	setFlash(c, "flash-notice", message)
}
func popNotice(c echo.Context) string {
	return popFlash(c, "flash-notice")
}

func setFlash(c echo.Context, name string, message string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = message
	c.SetCookie(cookie)
}

func popFlash(c echo.Context, name string) (value string) {
	cc, err := c.Cookie(name)
	if err == nil {
		value = cc.Value

		// remove flash cookie
		cookie := new(http.Cookie)
		cookie.Name = name
		cookie.Expires = time.Now().Add(-1 * (100 * time.Second))
		c.SetCookie(cookie)
	}

	return value
}
