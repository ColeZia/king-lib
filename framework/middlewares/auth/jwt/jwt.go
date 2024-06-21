package jwt

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt"
	"gl.king.im/king-lib/framework/transport/metadata"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type authKey struct{}

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized(reason, "JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized(reason, "keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized(reason, "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized(reason, "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized(reason, "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized(reason, "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized(reason, "Token provider is missing")
	ErrSignToken              = errors.Unauthorized(reason, "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized(reason, "Can not get key while signing token")
)

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
	tokenHeader   map[string]interface{}
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithClaims with customer claim
// If you use it in Server, f needs to return a new jwt.Claims object each time to avoid concurrent write problems
// If you use it in Client, f only needs to return a single object to provide performance
func WithClaims(f func() jwt.Claims) Option {
	return func(o *options) {
		o.claims = f
	}
}

// WithTokenHeader withe customer tokenHeader for client side
func WithTokenHeader(header map[string]interface{}) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(keyFunc jwt.Keyfunc, opts ...Option) middleware.Middleware {
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				if keyFunc == nil {
					return nil, ErrMissingKeyFunc
				}
				auths := strings.SplitN(header.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}
				jwtToken := auths[1]
				var (
					tokenInfo *jwt.Token
					err       error
				)
				if o.claims != nil {
					tokenInfo, err = jwt.ParseWithClaims(jwtToken, o.claims(), keyFunc)
				} else {
					tokenInfo, err = jwt.Parse(jwtToken, keyFunc)
				}
				if err != nil {
					ve, ok := err.(*jwt.ValidationError)
					if !ok {
						return nil, errors.Unauthorized(reason, err.Error())
					}
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						return nil, ErrTokenInvalid
					}
					if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						return nil, ErrTokenExpired
					}
					return nil, ErrTokenParseFail
				}
				if !tokenInfo.Valid {
					return nil, ErrTokenInvalid
				}
				if tokenInfo.Method != o.signingMethod {
					return nil, ErrUnSupportSigningMethod
				}
				ctx = NewContext(ctx, tokenInfo.Claims)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// Client is a client jwt middleware.
func Client(opts ...Option) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			//			str, strok := transport.FromServerContext(ctx)
			//			ctr, ctrok := transport.FromClientContext(ctx)
			//			strkind := ""
			//			if strok {
			//				strkind = str.Kind().String()
			//			}
			//
			//			ctrkind := ""
			//			if ctrok {
			//				ctrkind = ctr.Kind().String()
			//			}
			//			fmt.Println("111 str.Kind():", strkind, "; strok:", strok, "; ctr.Kind():", ctrkind, "; ctrok", ctrok)
			//
			//			smd, sok := kmd.FromServerContext(ctx)
			//			cmd, cok := kmd.FromClientContext(ctx)
			//
			//			fmt.Println("111 smd:", smd, sok, "; cmd:", cmd, cok)

			ctx = metadata.BuildInnerMDCtxFromServerCtx(ctx)
			//ctx = metadata.BuildInnerMDCtx()

			//			smd, sok = kmd.FromServerContext(ctx)
			//			cmd, cok = kmd.FromClientContext(ctx)
			//
			//			fmt.Println("222 smd:", smd, sok, "; cmd:", cmd, cok)
			//
			//			str, strok = transport.FromServerContext(ctx)
			//			ctr, ctrok = transport.FromClientContext(ctx)
			//			strkind = ""
			//			if strok {
			//				strkind = str.Kind().String()
			//			}
			//
			//			ctrkind = ""
			//			if ctrok {
			//				ctrkind = ctr.Kind().String()
			//			}
			//			fmt.Println("222 str.Kind():", strkind, "; strok:", strok, "; ctr.Kind():", ctrkind, "; ctrok", ctrok)
			return handler(ctx, req)
		}
	}
}

// NewContext put auth info into context
func NewContext(ctx context.Context, info jwt.Claims) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (token jwt.Claims, ok bool) {
	token, ok = ctx.Value(authKey{}).(jwt.Claims)
	return
}
