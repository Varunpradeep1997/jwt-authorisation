package controllers 

import(
"context"
"fmt"
"log"
"strconv"
"net/http"
"time"
"github.com/gin-gonic/gin"
"github.com/go-playground/validator/v10"
helper "github.com/Varunpradeep1997/golang-jwt-projects/helpers"
"github.com/Varunpradeep1997/golang-jwt-projects/models"
"github.com/Varunpradeep1997/golang-jwt-projects/database"
"golang.org/x/crypto/bcrypt"

"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/bson/primitive"
"go.mongodb.org/mongo-driver/mongo"
"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()


var (
	googleOauthConfig=&oauth2.Config{
		RedirectURL: "http://localhost:8080/callback",
		ClientID: "18530797278-0kjpkbdt043hupkq6gei41vbmik8b1ib.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile","openid"},
		Endpoint: google.Endpoint,

	}
	randomState= "random"

)


func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err!=nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!= nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check=false
	}
	return check, msg
}

func Signup()gin.HandlerFunc{

	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the email"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})
		defer cancel()
		if err!= nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for the phone number"})
		}

		if count >0{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exists"})
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr !=nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}

}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}

		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 !=nil || page<1{
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
	matchStage, groupStage, projectStage})
defer cancel()
if err!=nil{
	c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
}
var allusers []bson.M
if err = result.All(ctx, &allusers); err!=nil{
	log.Fatal(err)
}
c.JSON(http.StatusOK, allusers[0])}}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}




/*google integration starts here*/




func HandleHome(w http.ResponseWriter, r *http.Request){
	var html = `<html><body><a href="/login">Google-Login</a></body></html>`
	fmt.Fprint(w,html)

}

func HandleLogin(w http.ResponseWriter, r *http.Request){
	url := googleOauthConfig.AuthCodeURL(randomState)
fmt.Print(url)
	/* AuthCodeURL returns a URL to OAuth 2.0 provider's consent page
that asks for permissions for the required scopes explicitly.
State is a token to protect the user from CSRF attacks. You must always provide a non-empty 
string and validate that it matches the the state query parameter on your redirect callback
*/

http.Redirect(w,r, url, http.StatusTemporaryRedirect)/*Redirect replies to the request
 with a redirect to url, which may be a path relative to the request path*/


}


func HandleCallback(w http.ResponseWriter, r *http.Request){

if	r.FormValue("state") != randomState{         //FormValue returns the first value for the named component of the query.
fmt.Printf("state is not valid,%s required,got %s\n",randomState,r.FormValue("state"))
http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
/*Redirect replies to the request
 with a redirect to url, which may be a path relative to the request path*/
return


}    
code :=r.FormValue("code")
token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
if err != nil {
	fmt.Printf("could not get token: %s\n",err.Error())
	http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
}

resp, err :=http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token="+ token.AccessToken)  //provides access token

if err != nil{
	fmt.Printf("could not create get request: %s\n",err.Error())
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}
defer resp.Body.Close()

content, err := ioutil.ReadAll(resp.Body)
if err != nil {

	fmt.Printf("could not create get request: %s\n",err.Error())
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return



}
fmt.Fprintf(w,"Response: %s" , content)
fmt.Println(string(content))

}

