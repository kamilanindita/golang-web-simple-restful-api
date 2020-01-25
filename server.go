package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"
    "encoding/json"
 
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)


type Buku struct {
    Id    int  `json: "id"`
    Penulis  string  `json: "penulis"`
    Judul string  `json: "judul"`
    Kota string  `json: "kota"`
    Penerbit string  `json: "penerbit"`
    Tahun int  `json: "tahun"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []Buku `json:"data"`
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "website_crud"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}


func HandlerIndex(w http.ResponseWriter, r *http.Request) {
    var tmp = template.Must(template.ParseFiles(
        "views/Header.html",
        "views/Menu.html",
        "views/Index.html",
        "views/Footer.html",
    ))
    data:=""
    var error = tmp.ExecuteTemplate(w,"Index",data)
    if error != nil {
        http.Error(w, error.Error(), http.StatusInternalServerError)
    }
}


func HandlerBuku(w http.ResponseWriter, r *http.Request) {

    db := dbConn()
    selDB, err := db.Query("SELECT id,penulis,judul,kota,penerbit,tahun FROM buku")
    if err != nil {
        panic(err.Error())
    }
    
    var buku=Buku{}
    data := []Buku{}
    var response Response

    for selDB.Next() {
        var id, tahun int
        var penulis, judul, kota, penerbit string
        err = selDB.Scan(&id, &penulis, &judul, &kota, &penerbit, &tahun)
        if err != nil {
            panic(err.Error())
        }
        buku = Buku{id, penulis, judul, kota, penerbit, tahun}
        data = append(data, buku)
    }
    defer db.Close()
    if len(data) > 0{
        response.Status = true
        response.Message = "Data Found"
        response.Data = data
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }else{
        response.Status = false
        response.Message = "Data not Found"
        response.Data = data
        w.WriteHeader(http.StatusNotFound)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response) 
    }
    
}


func HandlerBukuById(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]
    
    db := dbConn()
    selDB, err := db.Query("SELECT id,penulis,judul,kota,penerbit,tahun FROM buku WHERE id=?",key)
    if err != nil {
        panic(err.Error())
    }
    
    var buku=Buku{}
    data := []Buku{}
    var response Response

    for selDB.Next() {
        var id, tahun int
        var penulis, judul, kota, penerbit string
        err = selDB.Scan(&id, &penulis, &judul, &kota, &penerbit, &tahun)
        if err != nil {
            panic(err.Error())
        }
        buku = Buku{id, penulis, judul, kota, penerbit, tahun}
        data = append(data, buku)
    }
    defer db.Close()
    
    if len(data) > 0{
        response.Status = true
        response.Message = "Data Found"
        response.Data = data
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }else{
        response.Status = false
        response.Message = "Data not Found"
        response.Data = data
        w.WriteHeader(http.StatusNotFound)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
    
}


func HandlerSave(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    penulis := r.FormValue("penulis")
    judul := r.FormValue("judul")
    kota := r.FormValue("kota")
    penerbit := r.FormValue("penerbit")
    tahun := r.FormValue("tahun")
    insForm, err := db.Prepare("INSERT INTO buku (penulis,judul,kota,penerbit,tahun) VALUES(?,?,?,?,?)")
    if err != nil {
        panic(err.Error())
    }
    insForm.Exec(penulis, judul,  kota, penerbit, tahun)
    defer db.Close()
 
    data := []Buku{}
    var response Response
    response.Status = true
    response.Message = "Data has been created"
    response.Data = data
    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

}


func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    
    penulis := r.FormValue("penulis")
    judul := r.FormValue("judul")
    kota := r.FormValue("kota")
    penerbit := r.FormValue("penerbit")
    tahun := r.FormValue("tahun")
    vars := mux.Vars(r)
    key := vars["id"]
    insForm, err := db.Prepare("UPDATE buku SET penulis=?, judul=?, kota=?, penerbit=?, tahun=? WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    insForm.Exec(penulis, judul,  kota, penerbit, tahun, key)
    defer db.Close()    
    data := []Buku{}
    var response Response
    response.Status = true
    response.Message = "Data has been updated"
    response.Data = data
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

}


func HandlerDelete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    vars := mux.Vars(r)
    key := vars["id"]
    delForm, err := db.Prepare("DELETE FROM buku WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(key)
    defer db.Close()
    data := []Buku{}
    var response Response
    response.Status = true
    response.Message = "Data has been deleted"
    response.Data = data
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    
}

func main() { 
    Route := mux.NewRouter().StrictSlash(true)

    log.Println("Server started on: http://localhost:8080")
    Route.HandleFunc("/",HandlerIndex)
    Route.HandleFunc("/buku",HandlerBuku).Methods("GET")
    Route.HandleFunc("/buku", HandlerSave).Methods("POST")
    Route.HandleFunc("/buku/{id}",HandlerBukuById).Methods("GET")
    Route.HandleFunc("/buku/{id}",HandlerUpdate).Methods("PUT")
    Route.HandleFunc("/buku/{id}",HandlerDelete).Methods("DELETE")
    
    log.Fatal(http.ListenAndServe(":8080", Route))
} 