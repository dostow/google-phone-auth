package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/mgutz/logxi/v1"
	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

var logger = log.New("google-phone-auth")

// AddonConfig defines configuration of this handler
var AddonConfig = map[string]interface{}{
	"name":        "gpa",
	"title":       "Google phone authentication",
	"description": "An addon for authenticating via google phone auth",
	"properties": map[string]interface{}{
		"secret": map[string]interface{}{"type": "string", "description": "api secret key"},
	},
	"required":             []string{"secret"},
	"type":                 "object",
	"additionalProperties": false,
}

// AuthenticationContext context
type AuthenticationContext interface {
	Get(key string) (value interface{}, exists bool)
	MustGet(string) interface{}
	Param(string) string
}

// AuthenticationRouterFunction auth route function
type AuthenticationRouterFunction func(AuthenticationContext) (interface{}, error)

// AuthenticationStore describes a store
type AuthenticationStore interface {
	Save(string, map[string]interface{}) (string, error)
	Update(string, string, map[string]interface{}) error
	Search(string, map[string]interface{}) ([]map[string]interface{}, error)
}

func sendCode(c AuthenticationContext) (interface{}, error) {
	var err error
	phone := c.Param("phone")
	recaptcha := c.Param("recaptcha")

	authConfig := c.MustGet("authConfig").(map[string]interface{})
	headers := c.MustGet("headers").(http.Header)
	store := c.MustGet("authStore").(AuthenticationStore)
	ctx := context.Background()

	secret := authConfig["secret"].(string)
	identitytoolkitService, err := identitytoolkit.NewService(ctx, option.WithAPIKey(secret))
	if err != nil {
		return nil, err
	}

	req := identitytoolkitService.Relyingparty.SendVerificationCode(&identitytoolkit.IdentitytoolkitRelyingpartySendVerificationCodeRequest{
		PhoneNumber:    phone,
		RecaptchaToken: recaptcha,
	})

	req.Context(ctx)
	response, err := req.Do()
	if err != nil {
		return nil, err
	}
	_, err = store.Save("authRequest", map[string]interface{}{
		"ip":        headers.Get("X-FORWARDED-FOR"),
		"code":      response.SessionInfo,
		"phone":     phone,
		"recaptcha": recaptcha,
		"status":    "pending",
	})
	if err != nil {
		logger.Warn("failed saving state", "err", err.Error())
		return nil, err
	}

	// Check for errors (usually driver not valid)
	if err != nil {
		logger.Error("error %v", err)
		return nil, err
	}
	return nil, fmt.Errorf("code not generated")
}

func verifyCode(c AuthenticationContext) (interface{}, error) {
	ctx := context.Background()
	authConfig := c.MustGet("authConfig").(map[string]interface{})
	headers := c.MustGet("headers").(http.Header)
	store := c.MustGet("authStore").(AuthenticationStore)

	secret := authConfig["secret"].(string)
	phone := c.Param("phone")
	identitytoolkitService, err := identitytoolkit.NewService(ctx, option.WithAPIKey(secret))
	if err != nil {
		logger.Warn(fmt.Sprintf("error %v", err))
		return nil, err
	}

	res, err := store.Search("authRequest", map[string]interface{}{
		"phone":  phone,
		"ip":     headers.Get("X-FORWARDED-FOR"),
		"status": "pending",
	})
	if err == nil {
		phoneAuth := res[0]
		id := phoneAuth["id"].(string)
		code := phoneAuth["code"].(string)
		sessionInfo := phoneAuth["sessionInfo"].(string)
		req := identitytoolkitService.Relyingparty.VerifyPhoneNumber(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPhoneNumberRequest{
			Code:        code,
			SessionInfo: sessionInfo,
		})
		req.Context(ctx)
		response, err := req.Do()
		phoneAuth["verificationProof"] = response.VerificationProof
		phoneAuth["status"] = "done"
		if err == nil {
			err = store.Update("phoneauth", id, phoneAuth)
		}
	}
	logger.Warn("error retrieving social login state", "err", err)
	return nil, err
}

// HandlerRegistrar an addon registrar
type HandlerRegistrar interface {
	Add(name string, route AuthenticationRouter)
}

// AuthenticationRouter defines an authentication router
type AuthenticationRouter struct {
	name   string
	title  string
	config map[string]interface{}
	route  func(*gin.RouterGroup)
}

// Name name of auth route
func (a *AuthenticationRouter) Name() string {
	return a.name
}

// Title title of auth route
func (a *AuthenticationRouter) Title() string {
	return a.title
}

// Config config of auth route
func (a *AuthenticationRouter) Config() map[string]interface{} {
	return a.config
}

// Router actual router
func (a *AuthenticationRouter) Router() func(*gin.RouterGroup) {
	return a.route
}

// Register injects an addon into a registry
func Register(ar HandlerRegistrar) {
	ar.Add("gpa", AuthenticationRouter{
		config: AddonConfig,
		route: func(gr *gin.RouterGroup) {
			gr.GET("gpa/send/:phone/:recaptcha", func(c *gin.Context) {
				r, err := sendCode(c)
				if err != nil {
					c.AbortWithError(400, err)
				} else {
					c.JSON(200, r)
				}
			})
			gr.GET("gpa/send/:phone/:recaptcha", func(c *gin.Context) {
				r, err := verifyCode(c)
				if err != nil {
					c.AbortWithError(400, err)
				} else {
					c.JSON(200, r)
				}
			})
		},
	})
}
