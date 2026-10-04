
package main

import (
 "database/sql"
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
 if err != nil { log.Fatal(err) }
 db.Ping()

 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 http.HandleFunc("/del", del)
 
 p:=os.Getenv("PORT")
 if p=="" { p="10000" }
 log.Println("CJ1V APP LIVE")
 http.ListenAndServe(":"+p, nil)
}

func home(w http.ResponseWriter, r *http.Request){
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 var mapato int64
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida").Scan(&mapato)

 // 1. WANACHAMA dropdown - kutoka wanachama table
 opt:=""
 rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY jina")
 if rows!=nil{
  for rows.Next(){
   var id int; var j string
   rows.Scan(&id,&j)
   opt+="<option value='"+strconv.Itoa(id)+"|"+html.EscapeString(j)+"'>"+html.EscapeString(j)+"</option>"
  }
  rows.Close()
 }

 // 2. HISTORIA - moja kwa moja kutoka mahida, HAKUNA JOIN
 hist:=""
 r2,_:=db.Query("SELECT id,jina,aina,kiasi FROM mahida ORDER BY id DESC LIMIT 200")
 if r2!=nil{
  for r2.Next(){
   var id int; var j,a string; var k int64
   r2.Scan(&id,&j,&a,&k)
   hist+="<tr><td>"+strconv.Itoa(id)+"</td><td>"+html.EscapeString(j)+"</td><td>"+html.EscapeString(a)+"</td><td>"+strconv.FormatInt(k,10)+"</td><td><a href='/del?id="+strconv.Itoa(id)+"' onclick=\"return confirm('Futa?')\">Futa</a></td></tr>"
  }
  r2.Close()
 }
 if hist==""{ hist="<tr><td colspan=5 style='text-align:center;padding:20px'>Hakuna data kwenye mahida - lakini MAPATO bado ni "+strconv.FormatInt(mapato,10)+"</td></tr>" }

 htmlPage := `
<html><head><meta name=viewport content="width=device-width,initial-scale=1"><style>
body{font-family:sans-serif;background:#f5f5f5;margin:0;padding:10px}
.card{background:#fff;padding:15px;border-radius:10px;margin-bottom:10px;box-shadow:0 2px 5px #0001}
.top{background:#2e7d32;color:#fff;padding:20px;border-radius:10px;text-align:center;font-size:22px;font-weight:bold;margin-bottom:10px}
table{width:100%;border-collapse:collapse} th{background:#e8f5e9;padding:10px} td{padding:10px;border-bottom:1px solid #eee}
input,select{padding:10px;border-radius:5px;border:1px solid #ccc}
button{padding:10px 20px;background:#2e7d32;color:#fff;border:none;border-radius:5px;font-weight:bold}
</style></head><body>
<div class=top>MAPATO TZS `+strconv.FormatInt(mapato,10)+`</div>
<div class=card>
<h3>Ongeza Muamala - APP</h3>
<form action=/save method=POST>
<select name=mid>`+opt+`</select>
<select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Mkopo</option><option>Marejesho</option><option>Fomu</option></select>
<input name=kiasi type=number placeholder="Kiasi" required>
<button type=submit>Save kwenye APP</button>
</form>
</div>
<div class=card>
<h3>Historia ya APP (data zote)</h3>
<table><tr><th>ID</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Action</th></tr>`+hist+`</table>
</div>
</body></html>`
 fmt.Fprint(w, htmlPage)
}

func save(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 midFull := r.FormValue("mid") // id|jina
 k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64)
 aina:=r.FormValue("aina")
 
 // toa jina kutoka value
 jina := midFull
 id := 0
 for i,c := range midFull{
  if c=='|'{
   id,_=strconv.Atoi(midFull[:i])
   jina=midFull[i+1:]
   break
  }
 }
 // save - tumia jina moja kwa moja, ndio APP yako inavyofanya kazi
 _, err := db.Exec("INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6)", id, jina, k, aina, int(time.Now().Month()), time.Now().Year())
 if err != nil { log.Println("save err", err) }
 http.Redirect(w,r,"/",302)
}

func del(w http.ResponseWriter, r *http.Request){
 id:=r.URL.Query().Get("id")
 db.Exec("DELETE FROM mahida WHERE id=$1", id)
 http.Redirect(w,r,"/",302)
}
