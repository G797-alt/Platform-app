package main

import (
 "context"
 "database/sql"
 "fmt"
 "net/http"
 "os"
 "time"
 _ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
func init(){url:=os.Getenv("DATABASE_URL");db,_=sql.Open("pgx",url);ctx,_:=context.WithTimeout(context.Background(),10*time.Second);db.PingContext(ctx);db.Exec(`CREATE TABLE IF NOT EXISTS michango(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)}
func main(){
 http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","text/html")
  fmt.Fprint(w,`<!DOCTYPE html><html><head><meta name=viewport content="width=device-width,initial-scale=1"><style>body{font-family:sans-serif;padding:12px;background:#f5f5f5}.card{background:#fff;padding:15px;border-radius:10px;margin-bottom:10px}input,select,button{width:100%;padding:12px;margin:4px 0;border-radius:8px;border:1px solid #ccc}button{background:#0d7a3a;color:#fff;font-weight:bold}</style></head><body><div class=card><h2>Mahida</h2><form method=POST action=/ongeza><input name=jina placeholder="Jina" required><input name=kiasi type=number placeholder="Kiasi" required><select name=aina><option value=mapato>Mapato</option><option value=matoleo>Matoleo</option></select><button>Hifadhi</button></form></div><div class=card><h3>Ripoti</h3><iframe src=/ripoti style="width:100%;height:300px;border:0"></iframe></div></body></html>`)
 })
 http.HandleFunc("/ongeza",func(w http.ResponseWriter,r *http.Request){r.ParseForm();db.Exec("INSERT INTO michango(jina,kiasi,aina)VALUES($1,$2,$3)",r.FormValue("jina"),r.FormValue("kiasi"),r.FormValue("aina"));http.Redirect(w,r,"/",302)})
 http.HandleFunc("/ripoti",func(w http.ResponseWriter,r *http.Request){var a,b int64;db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM michango WHERE aina='mapato'").Scan(&a);db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM michango WHERE aina='matoleo'").Scan(&b);rows,_:=db.Query("SELECT jina,kiasi,aina FROM michango ORDER BY id DESC LIMIT 20");fmt.Fprintf(w,"<p>Mapato:%d<br>Matoleo:%d<br><b>Baki:%d</b></p><hr>",a,b,a-b);for rows.Next(){var j,aa string;var k int64;rows.Scan(&j,&k,&aa);fmt.Fprintf(w,"%s %d %s<br>",j,k,aa)}})
 port:=os.Getenv("PORT");if port==""{port="10000"};http.ListenAndServe(":"+port,nil)
}