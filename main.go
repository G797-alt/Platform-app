package main
import ("database/sql"; "fmt"; "html"; "log"; "net/http"; "os"; _ "github.com/lib/pq")
var db *sql.DB
func main(){
 dsn:=os.Getenv("DATABASE_URL")
 log.Println("DSN:", dsn[:20])
 var err error
 db,err=sql.Open("postgres",dsn)
 if err!=nil{log.Println("DB open error:",err)}
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
   w.Header().Set("Content-Type","text/html")
   fmt.Fprintf(w,"<h1>CJ1V LIVE - %s</h1><p>DB: %v</p><a href='/test'>Test DB</a>",os.Getenv("PORT"),err)
   if db!=nil{
     var c int; db.QueryRow("SELECT COUNT(*) FROM mahida").Scan(&c)
     fmt.Fprintf(w,"<p>Mahida count: %d</p>",c)
   }
 })
 http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request){
   rows,_:=db.Query("SELECT jina,kiasi,aina FROM mahida LIMIT 5")
   fmt.Fprint(w,"<table border=1>")
   for rows.Next(){var j,a string; var k int64; rows.Scan(&j,&k,&a); fmt.Fprintf(w,"<tr><td>%s</td><td>%d</td><td>%s</td></tr>",html.EscapeString(j),k,a)}
   fmt.Fprint(w,"</table>")
 })
 p:=os.Getenv("PORT"); if p==""{p="10000"}; log.Println("Listening on",p); http.ListenAndServe(":"+p,nil)
}
