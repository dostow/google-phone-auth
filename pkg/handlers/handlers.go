package handlers

import (
	"context"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/mgutz/logxi/v1"
	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

var logger = log.New("google-phone-auth")

// TODO: These interfaces need to be moved to an auth plugin module
// AuthenticationContext context
type AuthenticationContext interface {
	Get(key string) (value interface{}, exists bool)
	MustGet(string) interface{}
	Param(string) string
}

// AuthenticationRouterFunction auth route function
type AuthenticationRouterFunction func(AuthenticationContext) (interface{}, error)

// PluginDataStore describes a store
type PluginDataStore interface {
	Save(string, map[string]interface{}) (string, error)
	Update(string, string, map[string]interface{}) error
	Search(string, map[string]interface{}) ([]map[string]interface{}, error)
}

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

// AddonConfig defines configuration of this handler
var AuthRequestConfig = map[string]interface{}{
	"name":        "authrequest",
	"title":       "Authentication Request",
	"description": "Stores authentication request data",
	"properties": map[string]interface{}{
		"ip":          map[string]interface{}{"type": "string", "description": "ip of request"},
		"sessionInfo": map[string]interface{}{"type": "string", "description": "ip of request"},
		"phone":       map[string]interface{}{"type": "string", "description": "ip of request"},
		"recaptcha":   map[string]interface{}{"type": "string", "description": "ip of request"},
	},
	"required":             []string{"ip", "sessionInfo", "phone", "recaptcha"},
	"type":                 "object",
	"additionalProperties": false,
}

func onlyNumeric(v string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return v
	}
	return reg.ReplaceAllString(v, "")
}
func sendCode(c AuthenticationContext) (interface{}, error) {
	var err error
	phone := onlyNumeric(c.Param("phone"))
	recaptcha := c.Param("recaptcha")

	authConfig := c.MustGet("authConfig").(map[string]interface{})
	headers := c.MustGet("headers").(http.Header)
	store := c.MustGet("pluginDataStore").(PluginDataStore)
	ctx := context.Background()

	secret := authConfig["secret"].(string)
	identitytoolkitService, err := identitytoolkit.NewService(ctx, option.WithAPIKey(secret))
	if err != nil {
		return nil, err
	}

	req := identitytoolkitService.Relyingparty.SendVerificationCode(&identitytoolkit.IdentitytoolkitRelyingpartySendVerificationCodeRequest{
		PhoneNumber:    "+" + phone,
		RecaptchaToken: recaptcha,
	})

	req.Context(ctx)
	response, err := req.Do()
	if err != nil {
		return nil, err
	}
	var id string
	var authRequest = map[string]interface{}{
		"ip":          connectingIP(headers),
		"sessionInfo": response.SessionInfo,
		"phone":       phone,
		"recaptcha":   recaptcha,
		"status":      "pending",
	}
	logger.Info("created auth request", "authRequest", authRequest)
	id, err = store.Save("authRequest", authRequest)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func verifyCode(c AuthenticationContext) (interface{}, error) {
	ctx := context.Background()
	phone := onlyNumeric(c.Param("phone"))
	authConfig := c.MustGet("authConfig").(map[string]interface{})
	headers := c.MustGet("headers").(http.Header)
	store := c.MustGet("pluginDataStore").(PluginDataStore)

	secret := authConfig["secret"].(string)
	code := c.Param("code")
	identitytoolkitService, err := identitytoolkit.NewService(ctx, option.WithAPIKey(secret))
	if err != nil {
		return nil, err
	}

	res, err := store.Search("authRequest", map[string]interface{}{
		"ip":     connectingIP(headers),
		"phone":  phone,
		"status": "pending",
	})
	if err == nil {
		logger.Info("found auth requests", "res", res)
		phoneAuth := res[0]
		id := phoneAuth["id"].(string)
		// existingCode := phoneAuth["code"].(string)
		sessionInfo := phoneAuth["sessionInfo"].(string)
		req := identitytoolkitService.Relyingparty.VerifyPhoneNumber(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPhoneNumberRequest{
			Code:        code,
			SessionInfo: sessionInfo,
		})
		req.Context(ctx)
		var response *identitytoolkit.IdentitytoolkitRelyingpartyVerifyPhoneNumberResponse
		response, err = req.Do()
		if response != nil {
			phoneAuth["verificationProof"] = response.VerificationProof
			phoneAuth["status"] = "done"
			phoneAuth["phone"] = phone
			if err == nil {
				err = store.Update(id, "authRequest", phoneAuth)
				// get or create a user with the phone number
				return phoneAuth, err
			}
		}
	}
	return nil, err
}

// HandlerRegistrar an addon registrar
type HandlerRegistrar interface {
	Add(name string,
		config map[string]interface{},
		route func(*gin.RouterGroup),
		schemas []map[string]interface{},
	) error
}

// Register injects an addon into a registry
func Register(ar HandlerRegistrar) {
	ar.Add("gpa", AddonConfig, func(gr *gin.RouterGroup) {
		gr.GET("send/:phone/:recaptcha", func(c *gin.Context) {
			r, err := sendCode(c)
			if err != nil {
				c.AbortWithStatusJSON(400, map[string]interface{}{"err": err.Error()})
			} else {
				c.JSON(200, r)
			}
		})
		gr.GET("verify/:phone/:code", func(c *gin.Context) {
			r, err := verifyCode(c)
			if err != nil {
				c.AbortWithStatusJSON(400, map[string]interface{}{"err": err.Error()})
			} else {
				c.Set("verified_phone", onlyNumeric(r.(map[string]interface{})["phone"].(string)))
				provider_id := r.(map[string]interface{})["phone"].(string)
				c.Set("provider", "verified_phone")
				c.Set("provider_id", provider_id)
				signInWithProvider := c.MustGet("SignInWithProvider").(func(*gin.Context))
				signInWithProvider(c)
			}
		})
	},
		[]map[string]interface{}{AuthRequestConfig},
	)
}
