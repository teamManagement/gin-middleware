package ginmiddleware

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/goextension"
	"github.com/tjfoc/gmsm/sm2"
	"net/http"
)

func UserVerifyHeaderWithVerifyFn(fn func(header http.Header) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !fn(ctx.Request.Header) {
			ctx.Abort()
			ctx.Status(http.StatusMethodNotAllowed)
			return
		}
		ctx.Next()
	}
}

func UseVerifyUserAgent(targetUserAgent string) gin.HandlerFunc {
	return UseVerifyUserAgentWithVerifyFn(func(userAgent string) bool {
		return userAgent == targetUserAgent
	})
}

func UseVerifyUserAgentWithVerifyFn(fn func(userAgent string) bool) gin.HandlerFunc {
	return UserVerifyHeaderWithVerifyFn(func(header http.Header) bool {
		return fn(header.Get("User-Agent"))
	})
}

func DefaultOriginDataGet(ctx *gin.Context) []byte {
	return []byte(ctx.Request.RequestURI)
}

func DefaultRsaOriginDataGet(ctx *gin.Context) []byte {
	h, _ := coderutils.Hash(sha256.New(), []byte(ctx.Request.RequestURI))
	return h
}

func DefaultDigestGet(ctx *gin.Context) []byte {
	digest, _ := goextension.Bytes(ctx.GetHeader("_d")).DecodeBase64()
	return digest
}

type VerifySignRequestOption struct {
	// RSAPublicKey rsa公钥
	RSAPublicKey *rsa.PublicKey
	// SM2PublicKey Sm2公钥
	SM2PublicKey *sm2.PublicKey
	// OriginGet 原文获取， 默认使用 DefaultOriginDataGet, RSA默认使用 DefaultRsaOriginDataGet
	OriginGet func(ctx *gin.Context) []byte
	// DigestGet 签名结果获取, 默认使用 DefaultDigestGet
	DigestGet func(ctx *gin.Context) []byte
	// Hash hash算法获取, 默认使用 crypto.SHA256
	Hash crypto.Hash
}

func UseVerifySignRequestPath(option *VerifySignRequestOption) gin.HandlerFunc {
	if option != nil {
		if option.OriginGet == nil {
			if option.SM2PublicKey == nil && option.RSAPublicKey != nil {
				option.OriginGet = DefaultRsaOriginDataGet
			} else {
				option.OriginGet = DefaultOriginDataGet
			}
		}

		if option.DigestGet == nil {
			option.DigestGet = DefaultDigestGet
		}

		if option.Hash == 0 {
			option.Hash = crypto.SHA256
		}
	}
	return func(context *gin.Context) {
		if option == nil || (option.RSAPublicKey == nil && option.SM2PublicKey == nil) {
			context.Status(http.StatusBadRequest)
			context.Abort()
			return
		}

		origin := option.OriginGet(context)
		digest := option.DigestGet(context)

		verifyResult := false
		if option.SM2PublicKey != nil {
			verifyResult = option.SM2PublicKey.Verify(origin, digest)
		} else {
			verifyResult = rsa.VerifyPKCS1v15(option.RSAPublicKey, option.Hash, origin, digest) == nil
		}

		if !verifyResult {
			context.Status(http.StatusMethodNotAllowed)
			context.Abort()
			return
		}
		return
	}
}
