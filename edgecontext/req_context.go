package edgecontext

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/fizx/baseplate.go/httpbp"
	"github.com/fizx/baseplate.go/log"
	"github.com/fizx/baseplate.go/set"
	"github.com/fizx/baseplate.go/thriftbp"

	"github.com/apache/thrift/lib/go/thrift"
)

// Signer is used to return a signature of the given string.
type Signer func(s string) (string, error)

// HeaderTrustHandlerSigner returns a Signer that can be used to sign an
// edge context HTTP header using TrustHeaderSignature.SignEdgeContextHeader.
func HeaderTrustHandlerSigner(handler httpbp.TrustHeaderSignature, duration time.Duration) Signer {
	return func(s string) (string, error) {
		h := httpbp.EdgeContextHeaders{EdgeRequest: s}
		return handler.SignEdgeContextHeader(h, duration)
	}
}

// An EdgeRequestContext contains context info about an edge request.
type EdgeRequestContext struct {
	impl *Impl

	// header and raw should always be set during initialization
	header string
	raw    NewArgs

	// token will be validated on first use
	tokenOnce sync.Once
	token     *AuthenticationToken
}

// AuthToken either validates the raw auth token and cache it,
// or return the cached token.
//
// If the validation failed, the error will be logged.
func (e *EdgeRequestContext) AuthToken() *AuthenticationToken {
	e.tokenOnce.Do(func() {
		if token, err := e.impl.ValidateToken(e.raw.AuthToken); err != nil {
			log.Errorw("token validation failed", "err", err)
			e.token = nil
		} else {
			e.token = token
		}
	})
	return e.token
}

// AttachToContext attaches the header to thrift context.
func (e *EdgeRequestContext) AttachToContext(ctx context.Context) context.Context {
	ctx = thrift.SetHeader(ctx, thriftbp.HeaderEdgeRequest, e.header)
	headers := set.StringSliceToSet(thrift.GetWriteHeaderList(ctx))
	headers.Add(thriftbp.HeaderEdgeRequest)
	ctx = thrift.SetWriteHeaderList(ctx, headers.ToSlice())
	return ctx
}

// AttachHTTPHeader attaches the header to the http Headers
func (e *EdgeRequestContext) AttachHTTPHeader(h http.Header) {
	h.Add(httpbp.EdgeContextHeader, e.header)
}

// AttachSignedHTTPHeader attaches the header to the http Headers along with
// a signature header.
func (e *EdgeRequestContext) AttachSignedHTTPHeader(h http.Header, signer Signer) error {
	sig, err := signer(e.header)
	if err != nil {
		return err
	}
	e.AttachHTTPHeader(h)
	h.Add(httpbp.EdgeContextSignatureHeader, sig)
	return nil
}

// SessionID returns the session id of this request.
func (e *EdgeRequestContext) SessionID() string {
	return e.raw.SessionID
}

// DeviceID returns the device id of this request.
func (e *EdgeRequestContext) DeviceID() string {
	return e.raw.DeviceID
}

// User returns the info about the user of this request.
func (e *EdgeRequestContext) User() User {
	return User{
		e: e,
	}
}

// OAuthClient returns the info about the oauth client of this request.
//
// ok will be false if this request does not have a valid auth token.
func (e *EdgeRequestContext) OAuthClient() (client OAuthClient, ok bool) {
	token := e.AuthToken()
	if token == nil {
		return
	}
	return OAuthClient(*token), true
}

// Service returns the info about the client service of this request.
//
// ok will be false if this request does not have a valid auth token.
func (e *EdgeRequestContext) Service() (service Service, ok bool) {
	token := e.AuthToken()
	if token == nil {
		return
	}
	return Service(*token), true
}
