package app

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"os"
)

var DB *sql.DB

func InitDB() {
	connstring := fmt.Sprintf("mydb.sqlite")
	var err error

	_, oserr := os.Stat(connstring)
	isexist := os.IsExist(oserr)

	revel.INFO.Println("is exist? : ", isexist)

	DB, err = sql.Open("sqlite3", connstring)
	if err != nil {
		revel.INFO.Println("DB error: ", err)
	} else {
		revel.INFO.Println("DB connected")

		if false == isexist {
			_, errsql := DB.Exec(`
				CREATE TABLE Users(
					id INTEGER PRIMARY KEY AUTOINCREMENT, 
					login TEXT, 
					password TEXT);
				CREATE TABLE Sessions(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					secret TEXT,
					user_id INTEGER);
				CREATE TABLE Rooms(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT);
				CREATE TABLE UsersInRooms(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER,
					room_id INTEGER
					); 
				CREATE TABLE Messages (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					sender TEXT,
					message TEXT,
					room_id INTEGER,
					receiver_id INTEGER
					timestamp TEXT)`)
			if errsql != nil {
				revel.INFO.Println("error: ", errsql)
			}
		}
	}
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
