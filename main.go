package main
import ("database/sql"; "fmt"; "html"; "log"; "net/http"; "os"; "strconv"; _ "github.com/lib/pq")
var db *sql.DB
func main(){db,_=sql.Open("postgres",os.Getenv("DATABASE_URL"))
http.HandleFunc("/", home); http.HandleFunc("/debug", debug); http.HandleFunc("/save", save); http.HandleFunc("/edit", edit); http.HandleFunc("/del", del)
p:=os.Getenv("PORT"); if p==""{p="10000"}; log.Println("Listening",p); http.ListenAndServe(":"+p,nil)}
func home(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","text/html; charset=utf-8")
var total,cnt int64; db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(kiasi),0) FROM mahida`).Scan(&cnt,&total)
var mapato,matumizi int64
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%hisa%' OR aina ILIKE '%ada%' OR aina ILIKE '%faini%' OR aina ILIKE '%fomu%' OR aina ILIKE '%rejesho%' OR aina ILIKE '%michango%'`).Scan(&mapato)
if mapato==0 {mapato=total}
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'`).Scan(&matumizi)
baki:=mapato-matumizi
opt:=""; rows,_:=db.Query(`SELECT id,jina FROM wanachama ORDER BY jina`); if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j string; rows.Scan(&id,&j); opt+=fmt.Sprintf(`<option value=%d>%s</option>`,id,html.EscapeString(j))}}
if opt==""{opt="<option>Hamna wanachama - Nenda /debug</option>"}
hist:=""; rows,_=db.Query(`SELECT id, COALESCE(jina,'-'), COALESCE(aina,'-'), COALESCE(kiasi,0) FROM mahida ORDER BY id DESC LIMIT 100`)
if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j,a string; var k int64; rows.Scan(&id,&j,&a,&k)
hist+=fmt.Sprintf(`<tr><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">❌</button></td></tr>`,id,html.EscapeString(j),html.EscapeString(a),k,id,k,id)}}
if hist==""{hist=fmt.Sprintf(`<tr><td colspan=5 style='color:red'>MAHIDA IKO 0 ROWS! Data imefutika. Total rows=%d. <a href=/debug>Bonyeza hapa /debug uone</a></td></tr>`,cnt)}
fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><style>body{font-family:sans-serif;margin:0;padding:8px;background:#f5f5f5}.top{display:flex;gap:6px}.c{flex:1;padding:14px;border-radius:12px;color:#fff;text-align:center;font-weight:bold}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}.box{background:#fff;padding:10px;border-radius:10px;margin-bottom:8px}table{width:100%%}td,th{padding:6px;border:1px solid #ddd;font-size:12px}</style></head><body>
<div class=top><div class='c g'>MAPATO<br>TZS %d</div><div class='c r'>MATUMIZI<br>TZS %d</div><div class='c b'>Bakaa<br>TZS %d</div></div>
<div style='background:yellow;padding:6px;font-size:11px'>DEBUG: Mahida rows=%d Total=%d - <a href=/debug>/debug</a> kuona tables</div>
<div class=box><b>Ingiza</b><form action=/save method=POST style="display:flex;gap:4px"><select name=mid>%s</select><select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option><option>Mkopo</option><option>Marejesho</option></select><input name=kiasi type=number placeholder=kiasi required><button>Save</button></form></div>
<div class=box><b>Historia (Raw - bila JOIN)</b><table><tr><th>ID</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Edit</th></tr>%s</table></div>
<script>
function editD(id,o){var k=prompt('Badilisha kiasi LIVE!',o); if(k) fetch('/edit?type=d&id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}
function delD(id){if(confirm('Futa?')) fetch('/del?type=d&id='+id,{method:'POST'}).then(()=>location.reload())}
</script></body></html>`,mapato,matumizi,baki,cnt,total,opt,hist)}
func debug(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","text/html")
fmt.Fprintf(w,"<h2>DEBUG TABLES</h2>")
var c1,c2 int64; db.QueryRow(`SELECT COUNT(*) FROM wanachama`).Scan(&c1); db.QueryRow(`SELECT COUNT(*) FROM mahida`).Scan(&c2)
fmt.Fprintf(w,"Wanachama=%d Mahida=%d<br><hr>",c1,c2)
fmt.Fprintf(w,"<b>wanachama:</b><br>"); rows,_:=db.Query(`SELECT id,jina FROM wanachama LIMIT 50`); if rows!=nil{for rows.Next(){var id int; var j string; rows.Scan(&id,&j); fmt.Fprintf(w,"%d - %s<br>",id,j)}; rows.Close()}
fmt.Fprintf(w,"<hr><b>mahida raw 20:</b><br>"); rows,_=db.Query(`SELECT id,mwanachama_id,jina,aina,kiasi FROM mahida ORDER BY id DESC LIMIT 20`); if rows!=nil{for rows.Next(){var id,mid int; var j,a string; var k int64; rows.Scan(&id,&mid,&j,&a,&k); fmt.Fprintf(w,"%d mid=%d j=%s aina=%s k=%d<br>",id,mid,j,a,k)}; rows.Close()}
}
func save(w http.ResponseWriter, r *http.Request){r.ParseForm(); mid,_:=strconv.Atoi(r.FormValue("mid")); k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64); var j string; db.QueryRow(`SELECT jina FROM wanachama WHERE id=$1`,mid).Scan(&j); db.Exec(`INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,10,2026)`,mid,j,k,r.FormValue("aina")); http.Redirect(w,r,"/",302)}
func edit(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); k:=r.URL.Query().Get("k"); kv,_:=strconv.ParseInt(k,10,64); db.Exec(`UPDATE mahida SET kiasi=$1 WHERE id=$2`,kv,id); w.Write([]byte("ok"))}
func del(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); db.Exec(`DELETE FROM mahida WHERE id=$1`,id); w.Write([]byte("ok"))}