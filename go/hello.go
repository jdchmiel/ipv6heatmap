package main

import (
    "fmt"
    "time"
    "net/http"
    "html/template"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "encoding/json"
)




//a single coordinate
type coordinate struct {
    Lat float64
    Lon float64
}

//a set of coordinates
type coordinates struct {
    point []coordinate
}

type Page struct {
    Title string
    Body  string
}

func loadPage(s string) *Page{
    body := "testing stuff in here right now for \"" + s + "\" and all other pages"
    return &Page{Title: s, Body: body}
}



func welcome() {
    fmt.Println("*******************")
    log("") //print date
    fmt.Println("* GO GO GO        *")
    fmt.Println("* This is GO      *")
    fmt.Println("*******************")
}

func log(s string) {

    fmt.Println( time.Now().String() + " serve: " + s);
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
    fileName := "/usr/src/myapp/" + r.URL.Path[1:]
    log(fileName)
    http.ServeFile(w, r, fileName)
}
func jsHandler(w http.ResponseWriter, r *http.Request) {
    fileName := "/usr/src/myapp/" + r.URL.Path[1:]
    log(fileName)
    http.ServeFile(w, r, fileName)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
    pageName := r.URL.Path[1:]
    log(pageName)
    p := loadPage(pageName)
    t, _ := template.ParseFiles("/usr/src/myapp/views/app.html")
    t.Execute(w, p)
}

func serveSingle(pattern string, filename string) {
    http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        log("Single File: " + pattern)
        http.ServeFile(w, r, filename)
    })
}

func getFloatFromQS(name string, r *http.Request) float64 {
    val, err := strconv.ParseFloat(r.URL.Query().Get(name),64)
    if (err != nil) {
        val = 0 // todo decide what a safe value could actually be vs just bailing out and returning a 40X error
        log("Bad querystring value, expected a float for: " + name)
        //todo return invalid request code
    }
    return val
}


func serveAPI(pattern string, stmt sql.Stmt) {
    http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        if (r.URL.Path != "/api/coordinates") {
            log("Bad endpoint" + r.URL.Path)
            //todo return 404
            return
        }
        if (r.Method != "GET") {
            //todo 10.4.6 405 Method Not Allowed
            log("Bad request type")
            //The method specified in the Request-Line is not allowed for the resource identified by the Request-URI. The response MUST include an Allow header containing a list of valid methods for the requested resource.
            return
        }
        minLat := getFloatFromQS("minLat", r)
        maxLat := getFloatFromQS("maxLat", r)
        minLon := getFloatFromQS("minLon", r)
        maxLon := getFloatFromQS("maxLon", r)

        rows, err := stmt.Query(minLat,maxLat, minLon, maxLon)
        if ( err != nil) {
            log("query failure")
            //todo log this
            return
        }
        var coord coordinate
        var coords []coordinate
        for rows.Next() {
            err = rows.Scan(&coord.Lat, &coord.Lon)
            if (err != nil) {
                log("ERROR fetching resultset row")
            }
            coords = append(coords,coord)
        }
        log("API: " + r.URL.Path)
        json.NewEncoder(w).Encode(coords)
    })
}


func main() {
    //spit a message at the docker logs to see the container is running
    welcome()

    //https://github.com/go-sql-driver/mysql/wiki/Examples
    db, err := sql.Open("mysql", "admin:jared@tcp(ipv6heatmap_db_1:3306)/ip2location_database") //docker container name is db so host = db
    if err != nil {
        //todo self, read the comment on next line from the wiki and dont spill dumps on the front end
        panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()
    err = db.Ping()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    stmt, err := db.Prepare("SELECT latitude, longitude FROM ip2location_database WHERE latitude >= ? " +
                            "and latitude <= ? and longitude >= ? and longitude <= ? LIMIT 10000")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmt.Close()


    //handle the static include dirs
    http.HandleFunc("/css/", cssHandler)
    http.HandleFunc("/js/", jsHandler)

    //handle the root hosted files
    serveSingle("/favicon.ico", "/usr/src/myapp/favicon.ico")
    //todo robots.txt
    //todo sitemap.xml

    serveAPI("/api/", *stmt)
    http.HandleFunc("/", pageHandler)

    http.ListenAndServe(":8080", nil)
}
