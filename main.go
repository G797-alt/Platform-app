package main
import ("database/sql"; "encoding/csv"; "encoding/json"; "fmt"; "html"; "log"; "net/http"; "os"; "strconv"; "time"; _ "github.com/lib/pq")
var db *sql.DB
func main(){
db,_=sql.Open("postgres",os.Getenv("DATABASE_URL"))
http.HandleFunc("/", home); http.HandleFunc("/save", save); http.HandleFunc("/edit", edit); http.HandleFunc("/del", del)
http.HandleFunc("/backup", backup); http.HandleFunc("/restore", restore); http.HandleFunc("/csv", csvExport)
p:=os.Getenv("PORT"); if p==""{p="10000"}; log.Println("Live",p); http.ListenAndServe(":"+p,nil)
}
func home(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","text/html; charset=utf-8")
var mapato,matumizi int64
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina NOT ILIKE '%mkopo%' OR aina ILIKE '%rejesho%'`).Scan(&mapato)
if mapato==0{db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida`).Scan(&mapato)}
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'`).Scan(&matumizi)
baki:=mapato-matumizi

// wanachama options
opt:=""; r1,_:=db.Query(`SELECT id,jina FROM wanachama ORDER BY jina`); if r1!=nil{for r1.Next(){var id int; var j string; r1.Scan(&id,&j); opt+=fmt.Sprintf(`<option value=%d>%s</option>`,id,html.EscapeString(j))}; r1.Close()}

// Wanaorudisha - wote waliofanya marejesho mwezi huu
rud:=""; r2,_:=db.Query(`SELECT m.id,w.jina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.aina ILIKE '%rejesho%' ORDER BY m.id DESC`); if r2!=nil{for r2.Next(){var id int; var j string; var k int64; r2.Scan(&id,&j,&k); rud+=fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">X</button></td></tr>`,html.EscapeString(j),k,id,k,id)}; r2.Close()}
if rud==""{rud=`<tr><td colspan=3>Hakuna</td></tr>`}

// Wenye Deni - Charles Joseph Mgaya 3,000,000 / 300,000 / 2,700,000 kama pichani
deni:=""; r3,_:=db.Query(`SELECT w.id, w.jina, COALESCE(SUM(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END),0) as ch, COALESCE(SUM(CASE WHEN m.aina ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END),0) as ru FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id GROUP BY w.id,w.jina HAVING SUM(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END)>0`); if r3!=nil{for r3.Next(){var wid int; var j string; var ch,ru int64; r3.Scan(&wid,&j,&ch,&ru); b:=ch-ru
deni+=fmt.Sprintf(`<tr><td style='max-width:90px'>%s</td><td>%s</td><td>%s</td><td style='color:#c62828;font-weight:bold'>%s</td><td><button onclick="alert('Edit Chukua/Rudisha - nenda Historia chini')" style='background:#2e7d32;color:#fff;border:none;padding:6px 10px;border-radius:6px'>✏️</button></td></tr>`,html.EscapeString(j),fmtInt(ch),fmtInt(ru),fmtInt(b))}; r3.Close()}
if deni==""{deni=`<tr><td colspan=5>Hakuna deni</td></tr>`}

// Historia ya Mwezi Huu
hist:=""; curM:=int(time.Now().Month()); r4,_:=db.Query(`SELECT m.id,w.jina,m.aina,m.kiasi,m.mwezi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.mwezi=$1 ORDER BY m.id DESC`,curM); if r4!=nil{for r4.Next(){var id int; var j,a string; var k int64; var mw int; r4.Scan(&id,&j,&a,&k,&mw); hist+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">❌</button></td></tr>`,html.EscapeString(j),html.EscapeString(a),fmtInt(k),id,k,id)}; r4.Close()}
if hist==""{hist=`<tr><td colspan=4>Hakuna mwezi huu</td></tr>`}

fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><style>
body{font-family:sans-serif;margin:0;background:#f1f5f1;padding:8px}
.card{background:#fff;border-radius:12px;padding:12px;margin-bottom:10px;box-shadow:0 1px 3px #0001}
.top{display:flex;gap:6px;margin-bottom:10px}.c{flex:1;padding:12px;border-radius:10px;color:#fff;text-align:center;font-weight:bold;font-size:13px}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}
h3{color:#2e7d32;margin:8px 0;font-size:15px} table{width:100%%;font-size:13px} th{background:#e8f5e9;padding:8px;text-align:left} td{padding:8px;border-bottom:1px solid #eee}
.btn{width:100%%;padding:14px;border:none;border-radius:10px;color:#fff;font-weight:bold;margin-bottom:8px;font-size:14px}
</style></head><body>
<div class=top><div class='c g'>MAPATO<br>TZS %s</div><div class='c r'>MATUMIZI<br>TZS %s</div><div class='c b'>Bakaa<br>TZS %s</div></div>

<div class=card><form action=/save method=POST style="display:flex;gap:4px;flex-wrap:wrap"><select name=mid style="flex:1">%s</select><select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option><option>Mkopo</option><option>Marejesho</option></select><input name=kiasi type=number placeholder="Kiasi" required style="width:90px"><button style="background:#2e7d32;color:#fff;padding:10px 16px;border:none;border-radius:8px">Hifadhi</button></form></div>

<div class=card><h3>Wanaorudisha Mikopo</h3><table><tr><th>Jina</th><th>Kiasi</th><th>✏️ / X</th></tr>%s</table></div>

<div class=card><h3>Wenye Deni la Mkopo</h3><table><tr><th>Jina</th><th>Chukua</th><th>Rudisha</th><th>Baki</th><th>✏️</th></tr>%s</table></div>

<div class=card><h3>Historia ya Mwezi Huu</h3><table>%s</table></div>

<div class=card><h3>Hifadhi ya Kudumu & Tuma</h3>
<a href=/backup style="text-decoration:none"><button class=btn style="background:#ef6c00">💾 Pakua Backup ya Kudumu (JSON)</button></a>
<button class=btn style="background:#1565c0" onclick="document.getElementById('up').click()">📁 Rudisha Backup</button>
<form id=fup action=/restore method=POST enctype="multipart/form-data" style="display:none"><input id=up type=file name=file onchange="this.form.submit()"></form>
<a href=/csv style="text-decoration:none"><button class=btn style="background:#616161">📊 Pakua Excel CSV</button></a>
<button class=btn style="background:#2e7d32" onclick="shareWA()">📄 Tuma WhatsApp Ripoti</button>
</div>

<script>
function editD(id,o){var k=prompt('Edit kiasi - Ita-update kila mahali LIVE!',o); if(k) fetch('/edit?id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}
function delD(id){if(confirm('Futa?')) fetch('/del?id='+id,{method:'POST'}).then(()=>location.reload())}
function shareWA(){var t='*CJ1V Ripoti*%0AMAPATO: %s%0AMATUMIZI: %s%0ABAKAA: %s%0A'+new Date().toLocaleDateString(); window.open('https://wa.me/?text='+t,'_blank')}
</script></body></html>`,fmtInt(mapato),fmtInt(matumizi),fmtInt(baki),opt,rud,deni,hist,fmtInt(mapato),fmtInt(matumizi),fmtInt(baki))}
func fmtInt(n int64) string {s:=fmt.Sprintf("%d",n); if len(s)<=3{return s}; r:=""; c:=0; for i:=len(s)-1;i>=0;i--{r=string(s[i])+r; c++; if c%%3==0 && i!=0{r=","+r}}; return r}
func save(w http.ResponseWriter, r *http.Request){r.ParseForm(); mid,_:=strconv.Atoi(r.FormValue("mid")); k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64); var j string; db.QueryRow(`SELECT jina FROM wanachama WHERE id=$1`,mid).Scan(&j); db.Exec(`INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6)`,mid,j,k,r.FormValue("aina"),int(time.Now().Month()),time.Now().Year()); http.Redirect(w,r,"/",302)}
func edit(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); k:=r.URL.Query().Get("k"); kv,_:=strconv.ParseInt(k,10,64); db.Exec(`UPDATE mahida SET kiasi=$1 WHERE id=$2`,kv,id); w.Write([]byte("ok"))}
func del(w http.ResponseWriter, r *http.Request){id:=r.URL.Query().Get("id"); db.Exec(`DELETE FROM mahida WHERE id=$1`,id); w.Write([]byte("ok"))}
func backup(w http.ResponseWriter, r *http.Request){rows,_:=db.Query(`SELECT id,mwanachama_id,jina,aina,kiasi,mwezi,mwaka FROM mahida`); var arr []map[string]interface{}; if rows!=nil{for rows.Next(){var id,mid,mw,yr int; var j,a string; var k int64; rows.Scan(&id,&mid,&j,&a,&k,&mw,&yr); arr=append(arr,map[string]interface{}{"id":id,"mid":mid,"jina":j,"aina":a,"kiasi":k,"mwezi":mw,"mwaka":yr})}; rows.Close()}; w.Header().Set("Content-Disposition","attachment; filename=cj1v-backup.json"); json.NewEncoder(w).Encode(arr)}
func restore(w http.ResponseWriter, r *http.Request){f,_,_:=r.FormFile("file"); var arr []map[string]interface{}; json.NewDecoder(f).Decode(&arr); for _,x:=range arr{db.Exec(`INSERT INTO mahida(id,mwanachama_id,jina,aina,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(id) DO NOTHING`,int(x["id"].(float64)),int(x["mid"].(float64)),x["jina"],x["aina"],int64(x["kiasi"].(float64)),int(x["mwezi"].(float64)),int(x["mwaka"].(float64)))}; http.Redirect(w,r,"/",302)}
func csvExport(w http.ResponseWriter, r *http.Request){w.Header().Set("Content-Disposition","attachment; filename=cj1v.csv"); w.Header().Set("Content-Type","text/csv"); cw:=csv.NewWriter(w); cw.Write([]string{"ID","Jina","Aina","Kiasi","Mwezi"}); rows,_:=db.Query(`SELECT id,jina,aina,kiasi,mwezi FROM mahida`); if rows!=nil{for rows.Next(){var id,mw int; var j,a string; var k int64; rows.Scan(&id,&j,&a,&k,&mw); cw.Write([]string{fmt.Sprintf("%d",id),j,a,fmt.Sprintf("%d",k),fmt.Sprintf("%d",mw)})}; rows.Close()}; cw.Flush()}