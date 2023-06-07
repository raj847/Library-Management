package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing/handler/api"
	"testing/middleware"
	"testing/repository"
	"testing/service"
	"testing/utils"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	AuthorAPIHandler  *api.AuthorAPI
	BookAPIHandler    *api.BookAPI
	MappingAPIHandler *api.MappingAPI
	UserAPIHandler    *api.UserAPI
}

func main() {
	err := os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/testing")
	if err != nil {
		log.Fatalf("cannot set env: %v", err)
	}

	mux := http.NewServeMux()

	err = utils.ConnectDB()
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	db := utils.GetDBConnection()
	mux = RunServer(db, mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {

	authorRepo := repository.NewAuthorRepository(db)
	bookRepo := repository.NewBookRepository(db)
	mappingRepo := repository.NewMappingRepository(db)
	userRepo := repository.NewUserRepository(db)

	authorService := service.NewAuthorService(authorRepo)
	bookService := service.NewBookService(bookRepo)
	mappingService := service.NewMappingService(mappingRepo)
	userService := service.NewUserService(userRepo)

	authorAPIHandler := api.NewAuthorAPI(authorService)
	bookAPIHandler := api.NewBookAPI(bookService)
	mappingAPIHandler := api.NewMappingAPI(mappingService)
	userAPIHandler := api.NewUserAPI(userService)

	apiHandler := APIHandler{
		AuthorAPIHandler:  authorAPIHandler,
		BookAPIHandler:    bookAPIHandler,
		MappingAPIHandler: mappingAPIHandler,
		UserAPIHandler:    userAPIHandler,
	}

	//AUTHOR
	MuxRoute(mux, "POST", "/testing/v1.0/author/create",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.AuthorAPIHandler.CreateNewAuthor))))

	MuxRoute(mux, "GET", "/testing/v1.0/author",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.AuthorAPIHandler.GetAllAuthor)),
		),
		"?author_id=",
	)

	MuxRoute(mux, "PUT", "/testing/v1.0/author/update",
		middleware.Put(
			middleware.Auth(
				http.HandlerFunc(apiHandler.AuthorAPIHandler.UpdateAuthor))),
		"?author_id=",
	)

	MuxRoute(mux, "DELETE", "/testing/v1.0/author/delete",
		middleware.Delete(
			middleware.Auth(
				http.HandlerFunc(apiHandler.AuthorAPIHandler.DeleteAuthor))),
		"?author_id=",
	)

	//BOOK
	MuxRoute(mux, "POST", "/testing/v1.0/book/create",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BookAPIHandler.CreateNewBook))))

	MuxRoute(mux, "GET", "/testing/v1.0/book",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BookAPIHandler.GetAllBook)),
		),
		"?book_id=",
	)

	MuxRoute(mux, "PUT", "/testing/v1.0/book/update",
		middleware.Put(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BookAPIHandler.UpdateBook))),
		"?book_id=",
	)

	MuxRoute(mux, "DELETE", "/testing/v1.0/book/delete",
		middleware.Delete(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BookAPIHandler.DeleteBook))),
		"?book_id=",
	)

	//Mapping
	MuxRoute(mux, "POST", "/mapping/create",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.MappingAPIHandler.CreateNewMapping))))

	MuxRoute(mux, "GET", "/mapping/read",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.MappingAPIHandler.GetAllMapping))))

	//USER
	MuxRoute(mux, "POST", "/testing/v1.0/user/register", middleware.Post(
		http.HandlerFunc(
			apiHandler.UserAPIHandler.Register)))

	MuxRoute(mux, "POST", "/testing/v1.0/user/login", middleware.Post(
		http.HandlerFunc(
			apiHandler.UserAPIHandler.UserLogin)))

	return mux

}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}
