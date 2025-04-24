package handlers

import (
	"backend/middleware"
	"net/http"
)

const (
	UserPath        = "/user"
	AssetsPath      = "/assets"
	ArticlesPath    = "/articles"
	CommentsPath    = "/comments"
	DiscussionsPath = "/discussions"
)

var (
	CORS   = middleware.CORS
	Method = middleware.Method
	Auth   = middleware.Authenticate
)

func emptyHandler(http.ResponseWriter, *http.Request) {}

func InitHandlers() {
	/* User */
	http.HandleFunc(UserPath+"/login", CORS(Method("POST", LoginUser)))
	http.HandleFunc(UserPath+"/get-user", CORS(Auth(Method("GET", GetUser))))
	http.HandleFunc(UserPath+"/register", CORS(Method("POST", RegisterUser)))
	http.HandleFunc(UserPath+"/verify-email", CORS(Method("PATCH", VerifyEmail)))
	http.HandleFunc(UserPath+"/logout", CORS(Auth(Method("DELETE", LogoutUser))))
	http.HandleFunc(UserPath+"/logged-in", CORS(Auth(Method("GET", emptyHandler))))
	http.HandleFunc(UserPath+"/modify-user", CORS(Auth(Method("PATCH", ModifyUser))))
	http.HandleFunc(UserPath+"/change-password", CORS(Method("PATCH", ChangePassword)))
	http.HandleFunc(UserPath+"/password-change-init", CORS(Method("POST", PasswordChangeInit)))
	http.HandleFunc(UserPath+"/password-change-valid", CORS(Method("GET", PasswordChangeValid)))

	/* Articles */
	http.HandleFunc(ArticlesPath+"/get", CORS(Method("GET", GetArticle)))
	http.HandleFunc(ArticlesPath+"/save", CORS(Method("PUT", Auth(SaveArticle))))
	http.HandleFunc(ArticlesPath+"/get-all", CORS(Method("GET", Auth(GetArticles))))
	http.HandleFunc(ArticlesPath+"/delete", CORS(Method("DELETE", Auth(DeleteArticle))))
	http.HandleFunc(ArticlesPath+"/get-suggested", CORS(Method("GET", Auth(GetSuggestedArticles))))

	/* Assets */
	http.HandleFunc(AssetsPath+"/add", CORS(Method("POST", Auth(AddAsset))))

	/* Comments */
	http.HandleFunc(CommentsPath+"/article/comment", CORS(Auth(Method("POST", CreateArticleComment))))
	http.HandleFunc(CommentsPath+"/article/update", CORS(Auth(Method("PATCH", UpdateArticleComment))))  // TODO: Validate whether the user is the author of the comment
	http.HandleFunc(CommentsPath+"/article/delete", CORS(Auth(Method("DELETE", DeleteArticleComment)))) // TODO: Validate whether the user is the author of the comment

	/* Discussions */
	http.HandleFunc(DiscussionsPath+"/article/start", CORS(Auth(Method("POST", emptyHandler))))       // TODO: Implement
	http.HandleFunc(DiscussionsPath+"/article/contribute", CORS(Auth(Method("PATCH", emptyHandler)))) // TODO: Implement
}
