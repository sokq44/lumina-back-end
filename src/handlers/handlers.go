package handlers

import (
	"backend/docs"
	"backend/middleware"
	"backend/models"
	"backend/utils/database"
	"backend/utils/emails"
	"backend/utils/jwt"
	"backend/utils/problems"
	"fmt"
	"net/http"
)

const (
	AuthPath        = "/auth"
	UserPath        = "/user"
	AssetsPath      = "/assets"
	ArticlesPath    = "/articles"
	CommentsPath    = "/comments"
	DiscussionsPath = "/discussions"
)

var (
	db = database.GetDb()
	em = emails.GetEmails()

	CORS   = middleware.CORS
	Method = middleware.Method
	Auth   = middleware.Authenticate
)

func EmptyHandler(http.ResponseWriter, *http.Request) {}

func GetUserFromRequest(r *http.Request) (*models.User, *problems.Problem) {
	access, err := r.Cookie("access_token")
	if err != nil {
		return nil, &problems.Problem{
			Type:          problems.HandlerProblem,
			ServerMessage: fmt.Sprintf("error while retrieving the access_token cookie: %v", err),
			ClientMessage: "An unexpected error has occurred while processing your request.",
			Status:        http.StatusInternalServerError,
		}
	}

	claims, p := jwt.DecodePayload(access.Value)
	if p != nil {
		return nil, p
	}

	return db.GetUserById(claims["user"].(string))
}

func InitHandlers(dev bool, port string) {
	/* Docs */
	if dev {
		http.HandleFunc(docs.Path, CORS(Method("GET", docs.GetSwagger(port))))
		http.HandleFunc(docs.OpenAPIPath, CORS(Method("GET", docs.GetOpenAPISpec(port))))
	}

	/* Session */
	http.HandleFunc(AuthPath+"/login", CORS(Method("POST", Login)))
	http.HandleFunc(AuthPath+"/register", CORS(Method("POST", Register)))
	http.HandleFunc(AuthPath+"/check", CORS(Auth(Method("GET", EmptyHandler))))
	http.HandleFunc(AuthPath+"/logout", CORS(Auth(Method("DELETE", Logout))))

	/* User */
	http.HandleFunc(UserPath+"/get", CORS(Auth(Method("GET", GetUser))))
	http.HandleFunc(UserPath+"/update", CORS(Auth(Method("PATCH", ModifyUser))))
	http.HandleFunc(UserPath+"/email/init", CORS(Auth(Method("POST", EmailChangeInit))))
	http.HandleFunc(UserPath+"/email/verify", CORS(Method("PATCH", VerifyEmail)))
	http.HandleFunc(UserPath+"/email/change", CORS(Auth(Method("PATCH", ChangeEmail))))
	http.HandleFunc(UserPath+"/password/init", CORS(Method("POST", PasswordChangeInit)))
	http.HandleFunc(UserPath+"/password/valid", CORS(Method("GET", PasswordChangeValid)))
	http.HandleFunc(UserPath+"/password/change", CORS(Method("PATCH", ChangePassword)))

	/* Articles */
	http.HandleFunc(ArticlesPath+"/get", CORS(Method("GET", GetArticle)))
	http.HandleFunc(ArticlesPath+"/save", CORS(Method("PUT", Auth(SaveArticle))))
	http.HandleFunc(ArticlesPath+"/all", CORS(Method("GET", Auth(GetArticles))))
	http.HandleFunc(ArticlesPath+"/delete", CORS(Method("DELETE", Auth(DeleteArticle))))
	http.HandleFunc(ArticlesPath+"/suggested", CORS(Method("GET", Auth(GetSuggestedArticles))))

	/* Assets */
	http.HandleFunc(AssetsPath+"/add", CORS(Method("POST", Auth(AddAsset))))

	/* Comments */
	http.HandleFunc(CommentsPath+"/article/all", CORS(Method("GET", GetAllArticleComments)))
	http.HandleFunc(CommentsPath+"/article/create", CORS(Auth(Method("POST", CreateArticleComment))))
	http.HandleFunc(CommentsPath+"/article/update", CORS(Auth(Method("PATCH", UpdateArticleComment))))
	http.HandleFunc(CommentsPath+"/article/delete", CORS(Auth(Method("DELETE", DeleteArticleComment))))

	/* Discussions */
	http.HandleFunc(DiscussionsPath+"/article/create", CORS(Auth(Method("POST", CreateArticleDiscussion))))
	http.HandleFunc(DiscussionsPath+"/article/update", CORS(Auth(Method("PATCH", UpdateArticleDiscussion))))
}
