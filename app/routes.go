package main

import (
	"app/data"
	"cloud0.christosploutarchou.com/cploutarchou/MicroGO/mailer"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

//routes The application routes.
func (a *application) routes() *chi.Mux {
	// Note : All the middlewares must come before the routes
	a.use(a.Middleware.CheckRemember)
	// Add routes here
	a.get("/", a.Handlers.Home)
	a.get("/go-page", a.Handlers.GoHome)
	a.get("/jet-page", a.Handlers.JetHome)
	a.get("/sessions", a.Handlers.SessionTest)
	a.get("/users/login", a.Handlers.UserLogin)
	a.post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/logout", a.Handlers.Logout)
	a.get("/users/forgot-password", a.Handlers.Forgot)
	a.get("/form", a.Handlers.Form)

	a.get("/json", a.Handlers.Json)
	a.get("/xml", a.Handlers.XML)
	a.get("/download", a.Handlers.Download)
	a.get("/encryption", a.Handlers.TestEncryption)
	a.get("/cache-test", a.Handlers.ShowCachePage)
	a.get("/cache-test", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)
	a.get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "your@mail.com",
			FromName:    "Christos Ploutarchou",
			To:          "cploutarchou@gmail.com",
			BCC:         nil,
			CC:          nil,
			Subject:     "Test SMTP - Send using channel",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		//a.App.Mailer.Jobs <- msg
		//res := <-a.App.Mailer.Results
		//if res.Error != nil {
		//	a.App.ErrorLog.Println(res.Error)
		//}
		err := a.App.Mailer.SentSMTPMessage(msg)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		fmt.Fprint(w, "Email Successfully send!")
	})
	a.App.Routes.Post("/form", a.Handlers.PostForm)
	a.get("/create_user", func(writer http.ResponseWriter, request *http.Request) {
		u := data.User{
			FirstName: "Christos",
			LastName:  "Ploutarchou",
			Email:     "cploutarchou@gmail.com",
			Active:    1,
			Password:  "mypassword",
		}

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		_, err = fmt.Fprintf(writer, "%d: %s", id, u.FirstName)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
	})
	a.get("/get_all_users", func(writer http.ResponseWriter, request *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		for _, x := range users {
			_, err := fmt.Fprintf(writer, x.FirstName)
			if err != nil {
				a.App.ErrorLog.Println(err)
				return
			}
		}
	})

	a.get("/get_user/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(request, "id"))
		user, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		_, err = fmt.Fprintf(writer, "%s %s %s %s", user.FirstName, user.LastName, user.Email, user.CreatedAt)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

	})
	a.get("/update_user/{id}", func(writer http.ResponseWriter, request *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(request, "id"))
		user, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		user.LastName = a.App.CreateRandomString(10)
		validator := a.App.Validator(nil)
		user.LastName = ""
		user.Validate(validator)
		if !validator.Valid() {
			fmt.Fprint(writer, "Validation failed")
			return
		}
		err = user.Update(*user)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		_, err = fmt.Fprintf(writer, "Update user id %d last Name to %s", id, user.LastName)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
	})
	// Routes for static files.
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
