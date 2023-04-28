package main

import (
	"backend/internal/models"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go get what you need",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.DB.AllProducts()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, products)

}

func (app *application) HistoriesByUser(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	reqToken := r.Header.Get("Authorization")

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userID := app.auth.RetrieveUser(reqToken)

	// get history by user
	histories, err := app.DB.GetHistoryByUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, histories)
}

func (app *application) GenerateCoupon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	product, _ := app.DB.OneProduct(productID)
	if product == nil {
		app.errorJSON(w, err)
		return
	}

	// get user id from context
	reqToken := r.Header.Get("Authorization")

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userID := app.auth.RetrieveUser(reqToken)

	// generate coupon
	coupon := models.Coupon{
		ProductID: productID,
		UserID:    userID,
		Code:      app.genRandomString(),
		Rate:      app.GetDiscountRateOnPrice(product.Price),
		ExpireAt:  time.Now().Add(15 * time.Minute),
		Active:    false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	// save to table coupon
	err = app.DB.SaveCoupon(&coupon)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, coupon.Code)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	//read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	//	check pw
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	//	create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	//	generate token
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println(tokens.Token)
	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)
		}
	}
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreate

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate user against database
	u, err := app.DB.GetUserByEmail(user.Email)
	if u != nil {
		app.errorJSON(w, errors.New("email already taken"), http.StatusBadRequest)
		return
	}

	//	hash pw
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		app.errorJSON(w, errors.New("err when hashing password"), http.StatusBadRequest)
		return
	}
	user.Password = string(hashedPassword)

	entity := models.MapCreateUser(user)

	err = app.DB.InsertUser(entity)
	if err != nil {
		app.errorJSON(w, errors.New("err when insert user"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, nil)
}

func (app *application) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	product, err := app.DB.OneProduct(productID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, product)
}

func (app *application) Purchase(w http.ResponseWriter, r *http.Request) {
	var purchase models.Purchase

	err := app.readJSON(w, r, &purchase)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	token := w.Header().Get("Authorization")

	userID := app.auth.RetrieveUser(token)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//check if the coupon match with coupon saved in table coupon base on the code and product_id and expired_at, active
	coupon, err := app.DB.GetCouponByCode(purchase.Coupon)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if coupon.ProductID != purchase.ProductID {
		app.errorJSON(w, errors.New("coupon code is not applicable for this product"))
		return
	}

	if !coupon.Active || coupon.ExpireAt.Before(time.Now()) {
		app.errorJSON(w, errors.New("coupon code expired"))
		return
	}

	order := models.History{
		ProductID: purchase.ProductID,
		UserID:    userID,
		Quantity:  purchase.Quantity,
		Discount:  coupon.Rate,
	}

	app.DB.AddHistory(&order)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
