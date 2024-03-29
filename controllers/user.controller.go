package controllers

import (
	"context"
	"fmt"
	"markie-backend/database"
	"markie-backend/models"
	"markie-backend/repository"
	"markie-backend/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userRepository = repository.Repository("users")
var spaceRepository = repository.Repository("spaces")
var client = database.DB

func hashPassword(password string) string {
	rounds, _ := strconv.Atoi(os.Getenv("SALT"))
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), rounds)
	return string(hash)
}

func comparePassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func CreateSpace(space *models.Space) {

}

func ResgisterUser(c *gin.Context) {
	session, sessionErr := client.StartSession()
	if sessionErr != nil {
		utils.Failure(c, bson.M{"message": "Session Error", "err": sessionErr}, http.StatusInternalServerError)
		return
	}
	defer session.EndSession(context.Background())
	transactionErr := session.StartTransaction()
	if transactionErr != nil {
		utils.Failure(c, bson.M{"message": "Transaction Error", "err": transactionErr}, http.StatusInternalServerError)
		return
	}
	var newUser models.User
	newUser.Id = primitive.NewObjectID()
	if err := c.BindJSON(&newUser); err != nil {
		utils.Failure(c, bson.M{"message": "Malformed Payload", "err": err}, http.StatusBadRequest)
		return
	}
	cur := userRepository.FindOne(&bson.M{"email": newUser.Email})
	var user models.User
	existingUser := cur.Decode(&user)
	if existingUser == nil {
		utils.Failure(c, bson.M{"message": "User with email already exisits"}, http.StatusConflict)
		return
	}
	passwordChannel := make(chan string)

	go func(password string) {
		hash := hashPassword(password)
		passwordChannel <- hash
	}(newUser.Password)
	hashedPassword := <-passwordChannel
	newUser.Password = hashedPassword
	var newSpace = &models.Space{
		Id:    primitive.NewObjectID(),
		Users: make([]models.SpaceUsers, 0),
		Name:  newUser.Firstname + "'s" + " Sapce",
		Admin: newUser.Id,
	}
	spaceIdChannel := make(chan string)
	go func(space *models.Space) {
		_, err := spaceRepository.InsertOne(space)
		if err != nil {
			utils.Failure(c, bson.M{"message": "Space creation failed"}, http.StatusInternalServerError)
			session.AbortTransaction(context.Background())
			return
		}
		spaceIdChannel <- newSpace.Id.String()
	}(newSpace)
	newUser.MySpaces = append(newUser.MySpaces, newSpace.Id)
	res, err := userRepository.InsertOne(newUser)
	if err != nil {
		utils.Failure(c, bson.M{"message": "Error occured while inserting doc", "err": err}, http.StatusInternalServerError)
		session.AbortTransaction(context.Background())
		return
	}
	if commitErr := session.CommitTransaction(context.Background()); commitErr != nil {
		utils.Failure(c, bson.M{"message": "Commit Error", "err": commitErr}, http.StatusInternalServerError)
		session.AbortTransaction(context.Background())
		return
	}
	utils.Success(c, bson.M{"id": res.InsertedID}, http.StatusCreated)
}

func GetUserById(c *gin.Context) {
	var id = c.Param("id")

	if len(id) != 24 {
		utils.Failure(c, bson.M{"message": "Invalid objectid"}, http.StatusBadRequest)
		return
	}

	validId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Failure(c, err, http.StatusBadRequest)
		return
	}

	var user models.User
	err = userRepository.FindOne(&bson.M{"_id": validId}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		utils.Failure(c, bson.M{"message": "User not found"}, http.StatusNotFound)
		return
	} else if err != nil {
		utils.Failure(c, err, http.StatusInternalServerError)
		return
	}

	utils.Success(c, user, http.StatusOK)

}

func GetUsers(c *gin.Context) {
	curr, err := userRepository.FindAll(&bson.M{})

	if err != nil {
		utils.Failure(c, err, http.StatusBadRequest)
		return
	}
	var results []models.User
	if err = curr.All(context.Background(), &results); err != nil {
		utils.Failure(c, err, http.StatusInternalServerError)
		return
	}
	if len(results) == 0 {
		utils.Success(c, bson.M{"users": results}, http.StatusNoContent)
		return
	}
	utils.Success(c, bson.M{"users": results}, http.StatusOK)
}

func Login(c *gin.Context) {
	var login models.Login

	if err := c.BindJSON(&login); err != nil {
		utils.Failure(c, err, http.StatusBadRequest)
		return
	}

	_, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	var user models.User
	err := userRepository.FindOne(&bson.M{"email": login.Email}).Decode(&user)
	if err == context.DeadlineExceeded {
		utils.Failure(c, bson.M{"message": "Database query timeout"}, http.StatusRequestTimeout)
		return
	} else if err == mongo.ErrNoDocuments {
		utils.Failure(c, bson.M{"message": "User not found", "email": login.Email}, http.StatusNotFound)
		return
	} else if err != nil {
		utils.Failure(c, err, http.StatusInternalServerError)
		return
	}
	passwordChannel := make(chan error)

	go func(password string, hashedPassword string) {
		err := comparePassword(password, hashedPassword)
		passwordChannel <- err
	}(login.Password, user.Password)

	err_ := <-passwordChannel

	if err_ != nil {
		utils.Failure(c, bson.M{"message": "Invalid password"}, http.StatusUnauthorized)
		return
	}

	id := string(user.Id.Hex())
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        id,
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		fmt.Println(err, jwtSecret)
		return
	}

	var result models.LoginSuccess

	result.Id = user.Id
	result.Email = user.Email
	result.Firstname = user.Firstname
	result.Lastname = user.Lastname
	result.UUID = uuid.New().String()

	if err != nil {
		utils.Failure(c, bson.M{"Message": "Redis Client failed to connect"}, http.StatusInternalServerError)
		return
	}

	if err := database.SetData(database.RedisClient, result.UUID, "Bearer "+tokenString, time.Hour*24); err != nil {
		utils.Failure(c, bson.M{"Message": "Could not set Token, Try Again"}, http.StatusInternalServerError)
	}

	utils.Success(c, bson.M{"message": "Login successful", "user": result}, http.StatusOK)
}

func VerifyUser(c *gin.Context) {
	userID, exists := c.Get("userId")

	if !exists {
		utils.Failure(c, bson.M{"error": "User ID not found in context"}, http.StatusInternalServerError)
		return
	}
	id, err_ := primitive.ObjectIDFromHex(userID.(string))
	if err_ != nil {
		utils.Failure(c, bson.M{"error": "Count not find user"}, http.StatusInternalServerError)
		return
	}

	var user models.User
	err := userRepository.FindOne(&bson.M{"_id": id}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		utils.Failure(c, bson.M{"message": "User not found"}, http.StatusNotFound)
		return
	} else if err != nil {
		utils.Failure(c, err, http.StatusInternalServerError)
		return
	}

	utils.Success(c, bson.M{"verified": true}, http.StatusOK)

}
