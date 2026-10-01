package main
import ("database/sql"; "fmt"; "html"; "log"; "net/http"; "os"; "strconv"; _ "github.com/lib/pq")
var db *sql.DB
func main(){db,_=sql.Open("postgres",os.Getenv("DATABASE_URL"))
http.HandleFunc("/", home); http.HandleFunc("/edit", edit); http.HandleFunc("/del", del); http.HandleFunc("/save", save)
p:=os.Getenv("PORT"); if p==""{p="10000"}; log.Println("Listening",p); http.ListenAndServe(":"+p,nil)}
func home(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","text/html; charset=utf-8")
var mapato,matumizi int64
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%hisa%' OR aina ILIKE '%ada%' OR aina ILIKE '%faini%' OR aina ILIKE '%fomu%' OR aina ILIKE '%rejesho%'`).Scan(&mapato)
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'`).Scan(&matumizi)
baki:=mapato-matumizi
opt:=""; rows,_:=db.Query(`SELECT id,jina FROM wanachama ORDER BY jina`); if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j string; rows.Scan(&id,&j); opt+=fmt.Sprintf(`<option value=%d>%s</option>`,id,html.EscapeString(j))}}
hist:=""; rows,_=db.Query(`SELECT m.id,w.jina,m.aina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id ORDER BY m.id DESC LIMIT 50`)
if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j,a string; var k int64; rows.Scan(&id,&j,&a,&k)
hist+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%d</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">❌</button></td></tr>`,html.EscapeString(j),html.EscapeString(a),k,id,k,id)}}
fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><style>body{font-family:sans-serif;margin:0;padding:8px;background:#f5f5f5}.top{display:flex;gap:6px;margin-bottom:8px}.c{flex:1;padding:14px;border-radius:12px;color:#fff;text-align:center;font-weight:bold}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}.box{background:#fff;padding:10px;border-radius:10px;margin-bottom:8px}table{width:100%%}td,th{padding:6px;border:1px solid #eee;font-size:13px}</style></head><body>
<div class=top><div class='c g'>MAPATO<br>TZS %d</div><div class='c r'>MATUMIZI<br>TZS %d</div><div class='c b'>Bakaa<br>TZS %d</div></div>
<div class=box><b>Ingiza</b><form action=/save method=POST style="display:flex;gap:4px"><select name=mid>%s</select><select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option><option>Mkopo</option><option>Marejesho</option></select><input name=kiasi type=number placeholder=kiasi required><button>Save</button></form></div>
<div class=box><b>Historia - Bonyeza ✏️ kiasi kinabadilika juu LIVE!</b><table><tr><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Edit</th></tr>%s</table></div>
<script>
function editD(id,o){var k=prompt('Badilisha kiasi - juu itabadilika LIVE!',o); if(k) fetch('/edit?type=d&id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}
function delD(id){if(confirm('Futa?')) fetch('/del?type=d&id='+id,{method:'POST'}).then(()=>location.reload())}
</script></body></html>`,mapato,matumizi,baki,opt,hist)}
func save(w http.ResponseWriter, r *http.Request){r.ParseForm(); mid,_:=strconv.Atoi(r.FormValue("mid")); k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64); var j string; db.QueryRow(`SELECT jina FROM wanachama WHERE id=$1`,mid).Scan(&j); db.Exec(`INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,10,2026)`,mid,j,k,r.FormValue("aina")); http.Redirect(w,r,"/",302)}
func edit(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); k:=r.URL.Query().Get("k"); kv,_:=strconv.ParseInt(k,10,64); db.Exec(`UPDATE mahida SET kiasi=$1 WHERE id=$2`,kv,id); w.Write([]byte("ok"))}
func del(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); db.Exec(`DELETE FROM mahida WHERE id=$1`,id); w.Write([]byte("ok"))}