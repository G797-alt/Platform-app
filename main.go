package main
import (
 "database/sql"; "fmt"; "html"; "log"; "net/http"; "os"; "strconv"
 _ "github.com/lib/pq"
)
var db *sql.DB
func main(){
 dsn:=os.Getenv("DATABASE_URL"); var err error; db,err=sql.Open("postgres",dsn); if err!=nil{log.Fatal(err)}
 db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY, jina TEXT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS mahida(id SERIAL PRIMARY KEY, mwanachama_id INT, jina TEXT, kiasi BIGINT, aina TEXT, maelezo TEXT, mwezi INT, mwaka INT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS mikopo_chukua(id SERIAL PRIMARY KEY, mwanachama_id INT, kiasi BIGINT, mwezi INT, mwaka INT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS mikopo_rudisha(id SERIAL PRIMARY KEY, mwanachama_id INT, kiasi BIGINT, mwezi INT, mwaka INT)`)
 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 http.HandleFunc("/editM", editM)
 http.HandleFunc("/delM", delM)
 http.HandleFunc("/editAll", editAll)
 p:=os.Getenv("PORT"); if p==""{p="10000"}; http.ListenAndServe(":"+p,nil)
}
func home(w http.ResponseWriter, r *http.Request){
 mz:=r.URL.Query().Get("mwezi"); if mz==""{mz="Okt"}
 mwS:=r.URL.Query().Get("mwaka"); if mwS==""{mwS="2026"}
 mw,_:=strconv.Atoi(mwS); mp:=map[string]int{"Jan":1,"Feb":2,"Mar":3,"Apr":4,"Mei":5,"Jun":6,"Jul":7,"Ago":8,"Sep":9,"Okt":10,"Nov":11,"Des":12}; mNum:=mp[mz]; if mNum==0{mNum=10}

 var mapato, matumizi int64
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE mwezi=$1 AND mwaka=$2 AND LOWER(aina) NOT IN ('mkopo','dharura','pole')`,mNum,mw).Scan(&mapato)
 var marejesho int64
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mikopo_rudisha WHERE mwezi=$1 AND mwaka=$2`,mNum,mw).Scan(&marejesho)
 mapato+=marejesho
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mikopo_chukua WHERE mwezi=$1 AND mwaka=$2`,mNum,mw).Scan(&matumizi)
 var matMahida int64
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE mwezi=$1 AND mwaka=$2 AND LOWER(aina) IN ('mkopo','dharura','pole')`,mNum,mw).Scan(&matMahida)
 matumizi+=matMahida
 baki:=mapato-matumizi

 // Wanachama
 wOpt:=""; wList:=""; rows,_:=db.Query(`SELECT id,jina FROM wanachama ORDER BY id`)
 if rows!=nil{ defer rows.Close(); for rows.Next(){ var id int; var j string; rows.Scan(&id,&j)
   wOpt+=fmt.Sprintf(`<option value=%d>%s</option>`,id,html.EscapeString(j))
   wList+=fmt.Sprintf(`<div style='display:flex;justify-content:space-between'><span>%s</span><span><button onclick="editW(%d,'%s')">✏️</button> <button onclick="delW(%d)">❌</button></span></div>`,html.EscapeString(j),id,html.EscapeString(j),id)
 }}

 // WANAORUDISHA - FIX: chukua kutoka mikopo_rudisha + mahida marejesho
 rudishaHTML:=""; rows,_=db.Query(`SELECT w.jina, SUM(r.kiasi) as tot, MAX(r.id) as lastid FROM mikopo_rudisha r JOIN wanachama w ON w.id=r.mwanachama_id WHERE r.mwezi=$1 AND r.mwaka=$2 GROUP BY w.jina`,mNum,mw)
 if rows!=nil{ defer rows.Close(); for rows.Next(){ var j string; var tot int64; var lid int; rows.Scan(&j,&tot,&lid)
   rudishaHTML+=fmt.Sprintf(`<tr><td>%s</td><td><input id='rud%d' value='%d' style='width:80px'><button onclick='saveEdit("r",%d)'>✏️</button></td><td><button onclick='delEdit("r",%d)'>X</button></td></tr>`,html.EscapeString(j),lid,tot,lid,lid)
 }}
 if rudishaHTML==""{rudishaHTML="<tr><td colspan=3>Hakuna</td></tr>"}

 // WENYE DENI - edit moja inabadilisha kila mahali
 deniHTML:=""; rows,_=db.Query(`
 SELECT w.id, w.jina, COALESCE(SUM(c.kiasi),0) as ch, COALESCE(SUM(r.kiasi),0) as ru
 FROM wanachama w
 LEFT JOIN mikopo_chukua c ON c.mwanachama_id=w.id AND c.mwezi=$1 AND c.mwaka=$2
 LEFT JOIN mikopo_rudisha r ON r.mwanachama_id=w.id AND r.mwezi=$1 AND r.mwaka=$2
 GROUP BY w.id,w.jina HAVING COALESCE(SUM(c.kiasi),0)>0
 `,mNum,mw)
 if rows!=nil{ defer rows.Close(); for rows.Next(){ var id int; var j string; var ch,ru int64; rows.Scan(&id,&j,&ch,&ru)
   bakiDeni:=ch-ru
   deniHTML+=fmt.Sprintf(`<tr><td>%s</td><td><input id='ch%d' value='%d' style='width:70px'></td><td><input id='ru%d' value='%d' style='width:70px'></td><td style='color:red'>%d</td><td><button onclick='saveDeni(%d)'>✏️</button></td></tr>`,html.EscapeString(j),id,ch,id,ru,bakiDeni,id)
 }}
 if deniHTML==""{deniHTML="<tr><td colspan=5>Hakuna deni</td></tr>"}

 // HISTORIA - FIX: onyesha yote ya mwezi huu
 histHTML:=""; rows,_=db.Query(`
 SELECT w.jina,m.aina,m.kiasi,m.id,'mahida' as tp FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.mwezi=$1 AND m.mwaka=$2
 UNION ALL SELECT w.jina,'Chukua',c.kiasi,c.id,'chukua' FROM mikopo_chukua c JOIN wanachama w ON w.id=c.mwanachama_id WHERE c.mwezi=$1 AND c.mwaka=$2
 UNION ALL SELECT w.jina,'Rudisha',r.kiasi,r.id,'rudisha' FROM mikopo_rudisha r JOIN wanachama w ON w.id=r.mwanachama_id WHERE r.mwezi=$1 AND r.mwaka=$2
 ORDER BY 1`,mNum,mw)
 if rows!=nil{ defer rows.Close(); for rows.Next(){ var j,aina,tp string; var k int64; var id int; rows.Scan(&j,&aina,&k,&id,&tp)
   histHTML+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%d</td><td><input id='h_%s_%d' value='%d' style='width:60px'><button onclick='saveHist("%s",%d)'>✏️</button></td></tr>`,html.EscapeString(j),html.EscapeString(aina),k,tp,id,k,tp,id)
 }}
 if histHTML==""{histHTML="<tr><td colspan=4>Hakuna mwezi huu - chagua mwezi mwingine au hifadhi data</td></tr>"}

 fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><style>body{font-family:sans-serif;background:#f5f5f5;margin:0}.top{background:#2e7d32;color:#fff;padding:12px;text-align:center;font-weight:bold}.f{display:flex;gap:5px;padding:8px}.cards{display:flex;gap:6px;padding:8px}.c{flex:1;padding:12px;border-radius:8px;color:#fff;text-align:center;font-weight:bold}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}.box{background:#fff;margin:8px;padding:10px;border-radius:8px} table{width:100%%;border-collapse:collapse} th,td{border:1px solid #eee;padding:6px;font-size:13px} select,input{width:100%%;padding:10px;margin:4px 0;border:1px solid #ccc;border-radius:6px}.btn{width:100%%;background:#2e7d32;color:#fff;padding:12px;border:none;border-radius:8px}</style></head><body>
<div class=top>Uwashusema Mahida - 2023 hadi 2035</div>
<div class=f><select id=mz><option %s>Okt</option><option>Sep</option><option>Nov</option></select><select id=mw><option>2026</option><option>2025</option></select><button onclick='location.href="/?mwezi="+document.getElementById("mz").value+"&mwaka="+document.getElementById("mw").value' style='background:#2e7d32;color:#fff'>Sawa</button></div>
<div class=cards><div class='c g'>MAPATO<br>TZS %d</div><div class='c r'>MATUMIZI<br>TZS %d</div><div class='c b'>Bakaa<br>TZS %d</div></div>
<div class=box><b>Wanachama</b> <button onclick='addW()' style='float:right'>+ Mpya</button><div>%s</div></div>
<div class=box><b>Ingiza Data</b><form action=/save method=POST><select name=mwanachama_id>%s</select><select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option><option>Mkopo</option><option>Marejesho</option><option>Dharura</option><option>Pole</option></select><input name=kiasi type=number placeholder=Kiasi><input name=maelezo placeholder=Maelezo><input type=hidden name=mwezi value=%s><input type=hidden name=mwaka value=%s><input type=hidden name=mNum value=%d><input type=hidden name=mwNum value=%d><button class=btn>Hifadhi</button></form></div>
<div class=box><b>Wanaorudisha Mikopo (LIVE - Edit inaonekana juu)</b><table><tr><th>Jina</th><th>Kiasi</th><th>✏️/X</th></tr>%s</table></div>
<div class=box><b>Wenye Deni la Mkopo (LIVE - Edit moja inabadilisha Mapato/Matumizi/Bakaa)</b><table><tr><th>Jina</th><th>Chukua</th><th>Rudisha</th><th>Baki</th><th>✏️</th></tr>%s</table></div>
<div class=box><b>Historia ya Mwezi Huu (LIVE)</b><table><tr><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Edit</th></tr>%s</table></div>
<div class=box><b>Hifadhi ya Kudumu & Tuma</b><button class=btn style='background:orange'>📦 Pakua Backup (JSON)</button></div>
<script>
function addW(){var j=prompt('Jina'); if(j) fetch('/editM?j='+j,{method:'POST'}).then(()=>location.reload())}
function editW(id,old){var j=prompt('Edit',old); if(j) fetch('/editM?id='+id+'&j='+j,{method:'POST'}).then(()=>location.reload())}
function delW(id){if(confirm('Futa?')) fetch('/delM?id='+id,{method:'POST'}).then(()=>location.reload())}
function saveEdit(tp,id){var k=document.getElementById((tp=='r'?'rud':'')+id).value; fetch('/editAll?tp='+tp+'&id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}
function delEdit(tp,id){if(confirm('Futa?')) fetch('/editAll?tp='+tp+'&id='+id+'&del=1',{method:'POST'}).then(()=>location.reload())}
function saveDeni(wid){var ch=document.getElementById('ch'+wid).value; var ru=document.getElementById('ru'+wid).value; fetch('/editAll?tp=deni&wid='+wid+'&ch='+ch+'&ru='+ru+'&m=%d&mw=%d',{method:'POST'}).then(()=>location.reload())}
function saveHist(tp,id){var k=document.getElementById('h_'+tp+'_'+id).value; fetch('/editAll?tp='+tp+'&id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload())}
</script></body></html>`,func()string{if mz=="Okt"{return "selected"};return ""}(), mapato, matumizi, baki, wList, wOpt, mz, mwS, mNum, mw, rudishaHTML, deniHTML, histHTML, mNum, mw)
}
func save(w http.ResponseWriter, r *http.Request){
 r.ParseForm(); mID,_:=strconv.Atoi(r.FormValue("mwanachama_id")); k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64); aina:=r.FormValue("aina"); ma:=r.FormValue("maelezo"); mNum,_:=strconv.Atoi(r.FormValue("mNum")); mwNum,_:=strconv.Atoi(r.FormValue("mwNum")); mz:=r.FormValue("mwezi"); mwS:=r.FormValue("mwaka")
 var j string; db.QueryRow(`SELECT jina FROM wanachama WHERE id=$1`,mID).Scan(&j)
 if aina=="Marejesho" || aina=="Rudisha"{ db.Exec(`INSERT INTO mikopo_rudisha(mwanachama_id,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4)`,mID,k,mNum,mwNum) } else if aina=="Mkopo"{ db.Exec(`INSERT INTO mikopo_chukua(mwanachama_id,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4)`,mID,k,mNum,mwNum) } else { db.Exec(`INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,maelezo,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6,$7)`,mID,j,k,aina,ma,mNum,mwNum) }
 http.Redirect(w,r,"/?mwezi="+mz+"&mwaka="+mwS,302)
}
func editM(w http.ResponseWriter, r *http.Request){ j:=r.URL.Query().Get("j"); id:=r.URL.Query().Get("id"); if id==""{db.Exec(`INSERT INTO wanachama(jina) VALUES($1)`,j)} else {db.Exec(`UPDATE wanachama SET jina=$1 WHERE id=$2`,j,id)}; w.Write([]byte("ok")) }
func delM(w http.ResponseWriter, r *http.Request){ id:=r.URL.Query().Get("id"); db.Exec(`DELETE FROM wanachama WHERE id=$1`,id); w.Write([]byte("ok")) }
func editAll(w http.ResponseWriter, r *http.Request){
 tp:=r.URL.Query().Get("tp"); id:=r.URL.Query().Get("id"); k:=r.URL.Query().Get("k"); del:=r.URL.Query().Get("del")
 kv,_:=strconv.ParseInt(k,10,64)
 if del=="1"{ if tp=="r"{db.Exec(`DELETE FROM mikopo_rudisha WHERE id=$1`,id)}; w.Write([]byte("ok")); return }
 if tp=="r" || tp=="rudisha"{ db.Exec(`UPDATE mikopo_rudisha SET kiasi=$1 WHERE id=$2`,kv,id) } else if tp=="c" || tp=="chukua"{ db.Exec(`UPDATE mikopo_chukua SET kiasi=$1 WHERE id=$2`,kv,id) } else if tp=="deni"{
   wid:=r.URL.Query().Get("wid"); ch:=r.URL.Query().Get("ch"); ru:=r.URL.Query().Get("ru"); m:=r.URL.Query().Get("m"); mw:=r.URL.Query().Get("mw"); chV,_:=strconv.ParseInt(ch,10,64); ruV,_:=strconv.ParseInt(ru,10,64); mI,_:=strconv.Atoi(m); mwI,_:=strconv.Atoi(mw)
   db.Exec(`DELETE FROM mikopo_chukua WHERE mwanachama_id=$1 AND mwezi=$2 AND mwaka=$3`,wid,mI,mwI); db.Exec(`DELETE FROM mikopo_rudisha WHERE mwanachama_id=$1 AND mwezi=$2 AND mwaka=$3`,wid,mI,mwI)
   if chV>0{db.Exec(`INSERT INTO mikopo_chukua(mwanachama_id,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4)`,wid,chV,mI,mwI)}
   if ruV>0{db.Exec(`INSERT INTO mikopo_rudisha(mwanachama_id,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4)`,wid,ruV,mI,mwI)}
 } else { db.Exec(`UPDATE mahida SET kiasi=$1 WHERE id=$2`,kv,id) }
 w.Write([]byte("ok"))
}