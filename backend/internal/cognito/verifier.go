package cognito

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"EMS/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWK represents a single JSON Web Key
type JWK struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
	Use string `json:"use"`
}

// JWKS represents the JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWTVerifier struct {
	Issuer   string
	ClientID string
	JWKSURL  string

	// Cache the keys so we don't do an HTTP request on every API call!
	keysMu sync.RWMutex
	keys   map[string]*rsa.PublicKey
}

func NewJWTVerifier(cfg *config.Config) (*JWTVerifier, error) {
	jwksURL := strings.TrimRight(
		cfg.CognitoOpenIDConfigURL,
		"/",
	) + "/.well-known/jwks.json"

	verifier := &JWTVerifier{
		Issuer:   cfg.CognitoOpenIDConfigURL,
		ClientID: cfg.CognitoClientID,
		JWKSURL:  jwksURL,
		keys:     make(map[string]*rsa.PublicKey),
	}

	// Fetch keys initially
	if err := verifier.refreshKeys(); err != nil {
		return nil, fmt.Errorf("failed to load initial JWKS: %w", err)
	}

	return verifier, nil
}

// refreshKeys downloads the JWKS from Cognito and converts them into RSA Public Keys
func (v *JWTVerifier) refreshKeys() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(v.JWKSURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch jwks: %s", resp.Status)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty != "RSA" {
			continue
		}

		// Decode the Modulus (N)
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			continue
		}
		n := new(big.Int).SetBytes(nBytes)

		// Decode the Exponent (E)
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			continue
		}

		// Convert E bytes to integer
		var eInt int
		if len(eBytes) < 4 {
			ndata := make([]byte, 4)
			copy(ndata[4-len(eBytes):], eBytes)
			eInt = int(binary.BigEndian.Uint32(ndata))
		} else {
			eInt = int(binary.BigEndian.Uint32(eBytes))
		}

		newKeys[key.Kid] = &rsa.PublicKey{
			N: n,
			E: eInt,
		}
	}

	v.keysMu.Lock()
	v.keys = newKeys
	v.keysMu.Unlock()

	return nil
}

// keyFunc is called by the jwt parser to find the right key based on the 'kid' header
func (v *JWTVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	// Verify that the signing method is RSA
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing kid in token header")
	}

	v.keysMu.RLock()
	pubKey, exists := v.keys[kid]
	v.keysMu.RUnlock()

	if !exists {
		// If key not found, maybe Cognito rotated keys. Fetch them again once.
		if err := v.refreshKeys(); err != nil {
			return nil, fmt.Errorf("failed to refresh keys: %v", err)
		}

		v.keysMu.RLock()
		pubKey, exists = v.keys[kid]
		v.keysMu.RUnlock()

		if !exists {
			return nil, fmt.Errorf("kid %s not found in JWKS", kid)
		}
	}

	return pubKey, nil
}

func (v *JWTVerifier) VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	// Validate the token is RSA256 signed, who is issuer, expiry.
	token, err := jwt.Parse(
		tokenString,
		v.keyFunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.Issuer),
	)

	if err != nil {
		return nil, err
	}

	// extracting claims from token payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify the token has the correct type
	tokenUse, ok := claims["token_use"].(string)
	if !ok || tokenUse != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Verify the token has the correct client
	clientID, ok := claims["client_id"].(string)
	if !ok || clientID != v.ClientID {
		return nil, fmt.Errorf("invalid client")
	}

	return token, nil
}

func (v *JWTVerifier) VerifyIDToken(tokenString string) (*jwt.Token, error) {
	// Validate the token is RSA256 signed, who is issuer, expiry.
	token, err := jwt.Parse(
		tokenString,
		v.keyFunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.Issuer),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify the token has the correct type
	tokenUse, ok := claims["token_use"].(string)
	if !ok || tokenUse != "id" {
		return nil, fmt.Errorf("invalid token type, expected id")
	}

	// For ID tokens, verify 'aud' is the client id
	aud, ok := claims["aud"].(string)
	if !ok || aud != v.ClientID {
		return nil, fmt.Errorf("invalid audience")
	}

	return token, nil
}
