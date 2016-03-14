package api

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type GrantType string

const (
	GrantPassword GrantType = "password"
	GrantRefreshToken = "refresh_token"
)

type ResponseType string

const (
	ResponseToken ResponseType = "token"
)

type TokenType string

const (
	TokenBearer TokenType = "Bearer"
)

type AuthRequest struct {
	ClientID string
	ClientSecret string
	GrantType GrantType
	Password string
	RedirectURI string
	ResponseType ResponseType
	State string
	Username string
}

type AuthResponse struct {
	Response
	AccessToken string
	ExpiresIn int
	TokenType TokenType
	Scope string
	Uid string
	RefreshToken string
	UserStatus int
	PrivateKey string
	EncPrivateKey string
	EventID string
}

type AuthCookiesRequest struct {
	ClientID string
	ResponseType ResponseType
	GrantType GrantType
	RefreshToken string
	RedirectURI string
	State string
}

func encrypt(user *backend.User, token string) (encrypted string, err error) {
	entitiesList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(user.PrivateKey))
	if err != nil {
		return
	}
	if len(entitiesList) == 0 {
		err = errors.New("Key ring does not contain any key")
		return
	}

	entity := entitiesList[0]

	var tokenBuffer bytes.Buffer
	armorWriter, err := armor.Encode(&tokenBuffer, "PGP MESSAGE", map[string]string{})
	if err != nil {
		return
	}

	
	w, err := openpgp.Encrypt(armorWriter, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return
	}

	w.Write([]byte(token))
	w.Close()

	armorWriter.Close()

	encrypted = tokenBuffer.String()
	return
}

func Auth(ctx *macaron.Context, req AuthRequest) {
	if req.GrantType != GrantPassword {
		ctx.JSON(200, &ErrorResponse{
			Response: Response{400},
			Error: "invalid_grant",
			ErrorDescription: "GrantType must be set to password",
		})
		return
	}

	user, err := backend.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(200, &ErrorResponse{
			Response: Response{401},
			Error: "invalid_grant",
			ErrorDescription: err.Error(),
		})
		return
	}

	encryptedToken, err := encrypt(user, "access_token")
	if err != nil {
		ctx.JSON(200, &ErrorResponse{
			Response: Response{500},
			Error: "invalid_key",
			ErrorDescription: err.Error(),
		})
		return
	}

	ctx.JSON(200, &AuthResponse{
		Response: Response{1000},
		AccessToken: encryptedToken,
		ExpiresIn: 360000,
		TokenType: TokenBearer,
		Scope: "full mail payments reset keys",
		Uid: user.Uid,
		RefreshToken: "refresh_token",
		PrivateKey: user.PrivateKey,
		EncPrivateKey: user.PrivateKey,
		EventID: "gnFPgsx4P9uXvB7IW8sIAUEcxEGGGH7mmRFiCmWwcn1jY3hxPxnCh39qvQInv5LkQFPn5rYh8qzfP_bJPrvHrg==",
	})
}

func AuthCookies(ctx *macaron.Context, req AuthCookiesRequest) {
	if req.GrantType != GrantRefreshToken {
		ctx.JSON(200, &ErrorResponse{
			Response: Response{400},
			Error: "invalid_grant",
			ErrorDescription: "GrantType must be set to refresh_token",
		})
		return
	}

	// TODO
}