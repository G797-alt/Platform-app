package main

import (
 "database/sql"
 "fmt"
 "html"
 "log"
 "net/http"
 "os"
 "strconv"
 _ "github.com/lib/pq"
)

var db *sql.DB

func main(){
 db,_ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 p:=os.Getenv("PORT")
 if p=="" { p="10000" }
 http.ListenAndServe(":"+p, nil)
}

func home(w http.ResponseWriter, r *http.Request){
 var mapato int64
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida").Scan(&mapato)
 
 opt:=""
 rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY jina")
 if rows!=nil{
  for rows.Next(){
   var id int; var j string
   rows.Scan(&id,&j)
   opt+="<option value='"+strconv.Itoa(id)+"'>"+html.EscapeString(j)+"</option>"
  }
  rows.Close()
 }
 hist:=""
 r2,_:=db.Query("SELECT m.id,w.jina,m.aina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id ORDER BY m.id DESC LIMIT 100")
 if r2!=nil{
  for r2.Next(){
   var id int; var j,a string; var k int64
   r2.Scan(&id,&j,&a,&k)
   hist+="<tr><td>"+strconv.Itoa(id)+"</td><td>"+j+"</td><td>"+a+"</td><td>"+strconv.FormatInt(k,10)+"</td></tr>"
  }
  r2.Close()
 }
 if hist==""{hist="<tr><td colspan=4>Hakuna data</td></tr>"}

 fmt.Fprint(w, "<html><body><h1>MAPATO TZS "+strconv.FormatInt(mapato,10)+"</h1><form action=/save method=POST><select name=mid>"+opt+"</select><input name=kiasi><button>Save</button></form><table>"+hist+"</table></body></html>")
}
func save(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 mid,_:=strconv.Atoi(r.FormValue("mid"))
 k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64)
 var j string
 db.QueryRow("SELECT jina FROM wanachama WHERE id=$1",mid).Scan(&j)
 db.Exec("INSERT INTO mahida(mwanachama_id,jina,kiasi,aina) VALUES($1,$2,$3,'Hisa')",mid,j,k)
 log.Println("saved")
 http.Redirect(w,r,"/",302)
}
