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

var (
	CORS   = middleware.CORS
	Auth   = middleware.Authenticate
	Method = middleware.Method
)

func Empty(http.ResponseWriter, *http.Request) {}

func InitHandlers() {
	assignUserHandlers()
	assignAssetsHandlers()
	assignArticlesHandlers()
}

func assignUserHandlers() {
	http.HandleFunc(UserPath+"/login", CORS(Method("POST", LoginUser)))

	http.HandleFunc(UserPath+"/register", CORS(Method("POST", RegisterUser)))

	http.HandleFunc(UserPath+"/verify-email", CORS(Method("PATCH", VerifyEmail)))

	http.HandleFunc(UserPath+"/logout", CORS(Auth(Method("DELETE", LogoutUser))))

	http.HandleFunc(UserPath+"/logged-in", CORS(Auth(Method("GET", Empty))))

	http.HandleFunc(UserPath+"/get-user", CORS(Auth(Method("GET", GetUser))))

	http.HandleFunc(UserPath+"/modify-user", CORS(Auth(Method("PATCH", ModifyUser))))

	http.HandleFunc(UserPath+"/change-password", CORS(Method("PATCH", ChangePassword)))

	http.HandleFunc(UserPath+"/password-change-init", CORS(Method("POST", PasswordChangeInit)))
}

func assignAssetsHandlers() {
	http.HandleFunc(AssetsPath+"/add", CORS(Method("POST", Auth(AddAsset))))
}

func assignArticlesHandlers() {
	http.HandleFunc(ArticlesPath+"/add", CORS(Method("POST", Auth(AddArticle))))

	http.HandleFunc(ArticlesPath+"/get", CORS(Method("GET", Auth(GetArticles))))

	http.HandleFunc(ArticlesPath+"/delete", CORS(Method("DELETE", Auth(DeleteArticle))))
}
