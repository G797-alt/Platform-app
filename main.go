package main

import (
 "database/sql"
 "encoding/csv"
 "encoding/json"
 "fmt"
 "html"
 "log"
 "net/http"
 "os"
 "strconv"
 "time"
 _ "github.com/lib/pq"
)

var db *sql.DB

func main(){
 var err error
 db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
 if err!= nil { log.Fatal(err) }
 if err = db.Ping(); err!= nil { log.Fatal(err) }
 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 http.HandleFunc("/save-mwanachama", saveMwanachama)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/del", del)
 http.HandleFunc("/backup", backup)
 http.HandleFunc("/csv", csvExport)
 http.HandleFunc("/restore", restore)
 p := os.Getenv("PORT")
 if p=="" { p="10000" }
 log.Println("CJ1V v2 MPYA LIVE on", p)
 http.ListenAndServe(":"+p, nil)
}

func fmtInt(n int64) string {
 if n==0 { return "0" }
 s := fmt.Sprintf("%d", n)
 out:=""
 c:=0
 for i:=len(s)-1;i>=0;i--{
  out=string(s[i])+out
  c++
  if c%3==0 && i!=0 { out=","+out }
 }
 return out
}

func home(w http.ResponseWriter, r *http.Request){
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 var mapato, matumizi int64
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida").Scan(&mapato)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%%mkopo%%' AND aina NOT ILIKE '%%rejesho%%'").Scan(&matumizi)
 baki := mapato - matumizi
 q := r.URL.Query().Get("q")

 opt:=""
 rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY jina ASC")
 if rows!=nil{
  for rows.Next(){
   var id int
   var j sql.NullString
   rows.Scan(&id,&j)
   jv:=""
   if j.Valid { jv=j.String }
   opt+=fmt.Sprintf(`<option value="%d|%s">%s</option>`, id, html.EscapeString(jv), html.EscapeString(jv))
  }
  rows.Close()
 }

 hist:=""
 sqlQ := "SELECT id,jina,aina,kiasi,COALESCE(mwezi,0),COALESCE(mwaka,0) FROM mahida "
 args:=[]interface{}{}
 if q!=""{
  sqlQ+=" WHERE jina ILIKE $1 OR aina ILIKE $1 "
  args
