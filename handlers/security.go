package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/rest-api/utils"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		// starts with Bearer?
		if !strings.HasPrefix(token, "Bearer") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			//http.Error(w, "", http.StatusUnauthorized)
			w.Write([]byte(`{"message": "token is wrong/missing"}`))
			return
		}
		// verify the token
		claims, err := verifyToken(token)
		if err != nil {
			utils.ReturnError(&w, fmt.Sprintf(`{"message": "%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}
		mapClaims := claims.(jwt.MapClaims)
		output.DebugLog("", fmt.Sprintf("received these claims: %v", claims))
		// inject token info into the call context
		// context creation to store product
		ctx := context.WithValue(r.Context(), "claims", mapClaims)
		// create a new request with the new context
		req := r.WithContext(ctx)

		// server the next handler
		next.ServeHTTP(w, req)
	})
}

// Login verifies username and password and create a JWT token in return.
func Login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Login handler func execution")
	// unmarshall json body
	var jBody map[string]string
	e := json.NewDecoder(r.Body)
	e.Decode(&jBody)

	username := jBody["username"]
	password := jBody["password"]

	if len(username) == 0 || len(password) == 0 {
		utils.ReturnError(&w, "please provide username and password to get the token", http.StatusBadRequest)
		return
	}
	// in a real scenario the user and password are checked on a database
	if username == "andrea" && password == "test" {
		// create the token
		token, err := createToken(username)
		if err != nil {
			utils.ReturnError(&w, "please provide username and password to get the token", http.StatusBadRequest)
			return
		}
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
	} else {
		utils.ReturnError(&w, `{"message": "credentials are not correct"}`, http.StatusUnauthorized)
		return
	}
}

// createToken creates the JWT token
func createToken(name string) (string, error) {
	signingKey := []byte(utils.Server.TokenPwd)
	// set expiration to 30 seconds
	expTime := time.Now().Add(5 * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"role": "admin",
		"exp":  expTime,
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

// verifyToken checks the token to confirm that is correctly signed and not expired
func verifyToken(token string) (jwt.Claims, error) {
	// check signature
	token = token[7:]
	if len(token) == 0 {
		return nil, fmt.Errorf("token has an incorrect format")
	}
	signingKey := []byte(utils.Server.TokenPwd)
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	exp, ok := t.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("incorrect exp payload property")
	}
	if ok := t.Claims.(jwt.MapClaims).VerifyExpiresAt(int64(exp), false); !ok {
		return nil, fmt.Errorf("token expired")
	}
	return t.Claims, nil
}
