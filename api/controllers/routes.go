package controllers

import "github.com/brianhumphreys/library_app/api/middlewares"

func (s *Server) initializeRoutes() {

	s.Router.HandleFunc("/api/v1/login", middlewares.CORS(middlewares.SetMiddlewareJSON(s.Login))).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/api/v1/signup", middlewares.CORS(middlewares.SetMiddlewareJSON(s.CreateUser))).Methods("POST", "OPTIONS")

	s.Router.HandleFunc("/api/v1/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/api/v1/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/api/v1/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/api/v1/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	s.Router.HandleFunc("/api/v1/books", middlewares.CORS(middlewares.SetMiddlewareJSON(s.CreateBook))).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/api/v1/books", middlewares.CORS(middlewares.SetMiddlewareJSON(s.GetBooks))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/api/v1/books/{id}", middlewares.CORS(middlewares.SetMiddlewareJSON(s.GetBook))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/api/v1/books/{id}", middlewares.CORS(middlewares.SetMiddlewareJSON(s.UpdateBook))).Methods("PUT", "OPTIONS")
	s.Router.HandleFunc("/api/v1/books/{id}", middlewares.CORS(s.DeleteBook)).Methods("DELETE", "OPTIONS")

	s.Router.HandleFunc("/api/v1/checkouts/current-books/{id}", middlewares.CORS(middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCurrentlyCheckedOutBooksOfUserWithID)))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/api/v1/checkouts/all-books/{id}", middlewares.CORS(middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetBookCheckoutHistoryOfUserWithID)))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/api/v1/checkouts/all-users/{id}", middlewares.CORS(middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUserCheckoutHistoryOfBookWithID)))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/api/v1/checkouts/checkout", middlewares.CORS(middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CheckoutABook)))).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/api/v1/checkouts/checkin", middlewares.CORS(middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CheckinABook)))).Methods("POST", "OPTIONS")
}
