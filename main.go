package main
import ("database/sql"; "fmt"; "html"; "net/http"; "os"; "strconv"; _ "github.com/lib/pq")
var db *sql.DB
func main(){dsn:=os.Getenv("DATABASE_URL"); db,_=sql.Open("postgres",dsn)
http.HandleFunc("/", home); http.HandleFunc("/save", save); http.HandleFunc("/edit", edit); http.HandleFunc("/del", del)
p:=os.Getenv("PORT"); if p==""{p="10000"}; http.ListenAndServe(":"+p,nil)}
func home(w http.ResponseWriter, r *http.Request){
var mapato,matumizi int64
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%hisa%' OR aina ILIKE '%ada%' OR aina ILIKE '%faini%' OR aina ILIKE '%fomu%'`).Scan(&mapato)
var marejesho int64; db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%rejesho%'`).Scan(&marejesho)
mapato+=marejesho
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'`).Scan(&matumizi)
baki:=mapato-matumizi

wOpt:=""; wList:=""; rows,_:=db.Query(`SELECT id,jina FROM wanachama ORDER BY id`)
if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j string; rows.Scan(&id,&j)
wOpt+=fmt.Sprintf(`<option value=%d>%s</option>`,id,html.EscapeString(j))
wList+=fmt.Sprintf(`<div style='display:flex;justify-content:space-between;padding:4px;border-bottom:1px solid #eee'><span>%s</span><span><button onclick="editW(%d,'%s')">✏️</button> <button onclick="delW(%d)">❌</button></span></div>`,html.EscapeString(j),id,html.EscapeString(j),id)}}

rud:=""; rows,_=db.Query(`SELECT m.id,w.jina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.aina ILIKE '%rejesho%' ORDER BY m.id DESC`)
if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j string; var k int64; rows.Scan(&id,&j,&k)
rud+=fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">X</button></td></tr>`,html.EscapeString(j),k,id,k,id)}}
if rud==""{rud="<tr><td colspan=3>Hakuna</td></tr>"}

deni:=""; rows,_=db.Query(`SELECT w.id, w.jina, COALESCE(SUM(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END),0) as ch, COALESCE(SUM(CASE WHEN m.aina ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END),0) as ru, MAX(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.id END) as ch_id, MAX(CASE WHEN m.aina ILIKE '%%rejesho%%' THEN m.id END) as ru_id FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id GROUP BY w.id,w.jina HAVING SUM(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END)>0`)
if rows!=nil{defer rows.Close(); for rows.Next(){var wid int; var j string; var ch,ru int64; var ch_id, ru_id sql.NullInt64; rows.Scan(&wid,&j,&ch,&ru,&ch_id,&ru_id)
b:=ch-ru
deni+=fmt.Sprintf(`<tr><td>%s</td><td>%d <button onclick="editD(%d,%d)">✏️</button></td><td>%d <button onclick="editD(%d,%d)">✏️</button></td><td style='color:red'>%d</td></tr>`,html.EscapeString(j),ch,ch_id.Int64,ch,ru,ru_id.Int64,ru,b)}}
if deni==""{deni="<tr><td colspan=4>Hakuna deni</td></tr>"}

hist:=""; rows,_=db.Query(`SELECT m.id,w.jina,m.aina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id ORDER BY m.id DESC LIMIT 30`)
if rows!=nil{defer rows.Close(); for rows.Next(){var id int; var j,aina string; var k int64; rows.Scan(&id,&j,&aina,&k)
hist+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%d</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">❌</button></td></tr>`,html.EscapeString(j),html.EscapeString(aina),k,id,k,id)}}
if hist==""{hist="<tr><td colspan=4>Hakuna - Ingiza data</td></tr>"}

fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><style>body{font-family:sans-serif;margin:0;background:#f0f0f0}.top{padding:10px;display:flex;gap:6px}.c{flex:1;padding:12px;border-radius:10px;color:#fff;text-align:center;font-weight:bold}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}.box{background:#fff;margin:8px;padding:10px;border-radius:10px} table{width:100%%} th,td{padding:6px;border:1px solid #eee;font-size:13px}</style></head><body>
<div class=top><div class='c g'>MAPATO<br>TZS %d</div><div class='c r'>MATUMIZI<br>TZS %d</div><div class='c b'>Bakaa<br>TZS %d</div></div>
<div class=box><b>Wanachama (%d)</b> <button onclick="addW()">+ Mpya</button><div>%s</div></div>
<div class=box><b>Ingiza Data</b><form action=/save method=POST><select name=mid>%s</select><select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option><option>Mkopo</option><option>Marejesho</option></select><input name=kiasi type=number placeholder=Kiasi required><input name=maelezo placeholder=Maelezo><button style='width:100%%;background:#2e7d32;color:#fff;padding:12px;border:none;border-radius:8px'>Hifadhi</button></form></div>
<div class=box><b>Wanaorudisha (Edit = juu inabadilika LIVE)</b><table><tr><th>Jina</th><th>Kiasi</th><th>✏️/X</th></tr>%s</table></div>
<div class=box><b>Wenye Deni (Edit Chukua/Rudisha = Bakaa LIVE)</b><table><tr><th>Jina</th><th>Chukua ✏️</th><th>Rudisha ✏️</th><th>Baki</th></tr>%s</table></div>
<div class=box><b>Historia (Edit hapa = kila mahali LIVE)</b><table><tr><th>Jina</th><th>Aina</th><th>Kiasi</th><th>✏️/❌</th></tr>%s</table></div>
<script>
function addW(){var j=prompt('Jina'); if(j) fetch('/edit?add='+encodeURIComponent(j),{method:'POST'}).then(()=>location.reload())}
function editW(id,old){var j=prompt('Edit jina',old); if(j) fetch('/edit?id='+id+'&j='+encodeURIComponent(j),{method:'POST'}).then(()=>location.reload())}
function delW(id){if(confirm('Futa?')) fetch('/del?type=w&id='+id,{method:'POST'}).then(()=>location.reload())}
function editD(id,old){var k=prompt('Edit kiasi LIVE - itabadilika kila mahali!',old); if(k) fetch('/edit?type=d&id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}
function delD(id){if(confirm('Futa?')) fetch('/del?type=d&id='+id,{method:'POST'}).then(()=>location.reload())}
</script></body></html>`,mapato,matumizi,baki,11,wList,wOpt,rud,deni,hist)}
func save(w http.ResponseWriter, r *http.Request){r.ParseForm(); mid,_:=strconv.Atoi(r.FormValue("mid")); k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64); aina:=r.FormValue("aina"); ma:=r.FormValue("maelezo"); var j string; db.QueryRow(`SELECT jina FROM wanachama WHERE id=$1`,mid).Scan(&j); db.Exec(`INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,maelezo,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,10,2026)`,mid,j,k,aina,ma); http.Redirect(w,r,"/",302)}
func edit(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); j:=r.URL.Query().Get("j"); k:=r.URL.Query().Get("k"); add:=r.URL.Query().Get("add"); tp:=r.URL.Query().Get("type")
if add!=""{db.Exec(`INSERT INTO wanachama(jina) VALUES($1)`,add)} else if tp=="d"{kv,_:=strconv.ParseInt(k,10,64); db.Exec(`UPDATE mahida SET kiasi=$1 WHERE id=$2`,kv,id)} else {db.Exec(`UPDATE wanachama SET jina=$1 WHERE id=$2`,j,id); db.Exec(`UPDATE mahida SET jina=$1 WHERE mwanachama_id=$2`,j,id)}; w.Write([]byte("ok"))}
func del(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); tp:=r.URL.Query().Get("type"); if tp=="w"{db.Exec(`DELETE FROM wanachama WHERE id=$1`,id); db.Exec(`DELETE FROM mahida WHERE mwanachama_id=$1`,id)} else {db.Exec(`DELETE FROM mahida WHERE id=$1`,id)}; w.Write([]byte("ok"))}