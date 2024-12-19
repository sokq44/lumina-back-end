package handlers

import (
	"backend/middleware"
	"net/http"
)

const (
	UserPath     = "/user"
	AssetsPath   = "/assets"
	ArticlesPath = "/articles"
)

func Empty(http.ResponseWriter, *http.Request) {}

func InitHandlers() {
	assignUserHandlers()
	assignAssetsHandlers()
	assignArticlesHandlers()
}

func assignUserHandlers() {
	http.HandleFunc(UserPath+"/login", middleware.CORS(middleware.Method("POST", LoginUser)))

	http.HandleFunc(UserPath+"/register", middleware.CORS(middleware.Method("POST", RegisterUser)))

	http.HandleFunc(UserPath+"/verify-email", middleware.CORS(middleware.Method("PATCH", VerifyEmail)))

	http.HandleFunc(UserPath+"/logout", middleware.CORS(middleware.Authenticate(middleware.Method("DELETE", LogoutUser))))

	http.HandleFunc(UserPath+"/logged-in", middleware.CORS(middleware.Authenticate(middleware.Method("GET", Empty))))

	http.HandleFunc(UserPath+"/get-user", middleware.CORS(middleware.Authenticate(middleware.Method("GET", GetUser))))

	http.HandleFunc(UserPath+"/modify-user", middleware.CORS(middleware.Authenticate(middleware.Method("PATCH", ModifyUser))))

	http.HandleFunc(UserPath+"/change-password", middleware.CORS(middleware.Method("PATCH", ChangePassword)))

	http.HandleFunc(UserPath+"/password-change-init", middleware.CORS(middleware.Method("POST", PasswordChangeInit)))
}

func assignAssetsHandlers() {
	http.HandleFunc(AssetsPath+"/add", middleware.CORS(middleware.Method("POST", middleware.Authenticate(AddAsset))))
}

func assignArticlesHandlers() {
	http.HandleFunc(ArticlesPath+"/add", middleware.CORS(middleware.Method("POST", middleware.Authenticate(AddArticle))))
}
