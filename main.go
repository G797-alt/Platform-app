package main
import ("database/sql"; "fmt"; "html"; "log"; "net/http"; "os"; "strconv"; _ "github.com/lib/pq")
var db *sql.DB
func main(){
 db,_=sql.Open("postgres",os.Getenv("DATABASE_URL"))
 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/del", del)
 p:=os.Getenv("PORT"); if p==""{p="10000"}
 log.Println("Listening",p)
 http.ListenAndServe(":"+p,nil)
}
func home(w http.ResponseWriter, r *http.Request){
 var mapato,matumizi int64
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida`).Scan(&mapato)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'`).Scan(&matumizi)
 baki:=mapato-matumizi
 hist:=""; rows,_:=db.Query(`SELECT m.id,w.jina,m.aina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id ORDER BY m.id DESC LIMIT 20`)
 if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j,a string; var k int64; rows.Scan(&id,&j,&a,&k)
 hist+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%d</td><td><button onclick="editD(%d,%d)">✏️</button></td></tr>`,html.EscapeString(j),a,k,id,k)}}
 fmt.Fprintf(w,`MAPATO %d | MATUMIZI %d | BAKAA %d<br><br><table border=1>%s</table><script>function editD(id,o){var k=prompt('Edit',o); if(k) fetch('/edit?type=d&id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}</script>`,mapato,matumizi,baki,hist)
}
func save(w http.ResponseWriter, r *http.Request){r.ParseForm(); mid,_:=strconv.Atoi(r.FormValue("mid")); k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64); db.Exec(`INSERT INTO mahida(mwanachama_id,kiasi,aina) VALUES($1,$2,$3)`,mid,k,r.FormValue("aina")); http.Redirect(w,r,"/",302)}
func edit(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); k:=r.URL.Query().Get("k"); kv,_:=strconv.Atoi(k); db.Exec(`UPDATE mahida SET kiasi=$1 WHERE id=$2`,kv,id); w.Write([]byte("ok"))}
func del(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); db.Exec(`DELETE FROM mahida WHERE id=$1`,id); w.Write([]byte("ok"))}