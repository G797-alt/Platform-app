package main
import(
"context"
"database/sql"
"encoding/csv"
"encoding/json"
"fmt"
"html"
"net/http"
"os"
"strconv"
"strings"
"time"
_ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
var err error
func init(){
u:=os.Getenv("DATABASE_URL")
if u==""{return}
db,err=sql.Open("pgx",u)
ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
defer cancel()
db.PingContext(ctx)
db.Exec("CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY, jina TEXT)")
db.Exec("CREATE TABLE IF NOT EXISTS mahida(id SERIAL PRIMARY KEY, wanachama_id INT, mwezi INT, mwaka INT, hisa TEXT, klasi TEXT, maelezo TEXT, kiasi BIGINT, aina TEXT, tarehe TIMESTAMP DEFAULT NOW())")
db.Exec("CREATE TABLE IF NOT EXISTS mikopo_chukua(id SERIAL PRIMARY KEY, wanachama_id INT, kiasi BIGINT)")
db.Exec("CREATE TABLE IF NOT EXISTS mikopo_rudisha(id SERIAL PRIMARY KEY, wanachama_id INT, kiasi BIGINT)")
var c int
db.QueryRow("SELECT COUNT(*) FROM wanachama").Scan(&c)
if c==0{
db.Exec("INSERT INTO wanachama(jina) VALUES ('Charles Joseph Mgaya')")
db.Exec("INSERT INTO mikopo_chukua(wanachama_id,kiasi) VALUES (1,3000000)")
db.Exec("INSERT INTO mikopo_rudisha(wanachama_id,kiasi) VALUES (1,300000)")
}
}
func main(){
http.HandleFunc("/",home)
http.HandleFunc("/add",addM)
http.HandleFunc("/del",delM)
http.HandleFunc("/save",save)
http.HandleFunc("/chukua",chukua)
http.HandleFunc("/rudisha",rudisha)
http.HandleFunc("/delr",delr)
http.HandleFunc("/backup",doBackup)
http.HandleFunc("/csv",doCSV)
p:=os.Getenv("PORT")
if p==""{p="10000"}
http.ListenAndServe(":"+p,nil)
}
func fmtN(n int64)string{
s:=strconv.FormatInt(n,10)
var b strings.Builder
for i,ch:=range s{
if i>0&&(len(s)-i)%3==0{b.WriteString(",")}
b.WriteRune(ch)
}
return b.String()
}
func home(w http.ResponseWriter,r *http.Request){
w.Header().Set("Content-Type","text/html; charset=utf-8")
mz:=r.URL.Query().Get("mwezi")
if mz==""{mz="Okt"}
mwS:=r.URL.Query().Get("mwaka")
if mwS==""{mwS="2026"}
mw,_:=strconv.Atoi(mwS)
mp:=map[string]int{"Jan":1,"Feb":2,"Mar":3,"Apr":4,"Mei":5,"Jun":6,"Jul":7,"Ago":8,"Sep":9,"Okt":10,"Nov":11,"Des":12}
mNum:=mp[mz]
if mNum==0{mNum=10}
var mapato,matumizi int64
type W struct{ID int; Jina string}
var list []W
if db!=nil{
db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina='mapato' AND mwezi=$1 AND mwaka=$2",mNum,mw).Scan(&mapato)
db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina='matumizi' AND mwezi=$1 AND mwaka=$2",mNum,mw).Scan(&matumizi)
rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY id")
if rows!=nil{defer rows.Close();for rows.Next(){var id int;var j string;rows.Scan(&id,&j);list=append(list,W{id,j})}}
}
baki:=mapato-matumizi
mo:="";for _,o:=range []string{"Jan","Feb","Mar","Apr","Mei","Jun","Jul","Ago","Sep","Okt","Nov","Des"}{s:="";if o==mz{s="selected"};mo+=fmt.Sprintf("<option %s>%s</option>",s,o)}
yo:="";for y:=2023;y<=2035;y++{s:="";if strconv.Itoa(y)==mwS{s="selected"};yo+=fmt.Sprintf("<option %s>%d</option>",s,y)}
mop:="";for _,wm:=range list{mop+=fmt.Sprintf("<option value=%d>%s</option>",wm.ID,html.EscapeString(wm.Jina))}
fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title><style>*{box-sizing:border-box;font-family:system-ui}body{margin:0;background:#f3f6f3}.h{background:#1a7a3b;color:#fff;padding:12px;text-align:center;font-weight:700;position:sticky;top:0}.c{padding:10px;max-width:650px;margin:auto}.r{display:flex;gap:8px;margin-bottom:10px}select,input{flex:1;padding:12px;border:1px solid #ddd;border-radius:10px;background:#fff}.b{padding:12px 16px;border:0;border-radius:10px;font-weight:700;color:#fff;background:#1a7a3b;cursor:pointer}.cards{display:grid;grid-template-columns:1fr 1fr 1fr;gap:8px;margin-bottom:14px}.card{padding:12px;border-radius:12px;color:#fff;text-align:center;font-weight:700}.card small{display:block;font-weight:400;font-size:11px}.card b{font-size:20px}.box{background:#fff;border-radius:12px;padding:12px;margin-bottom:12px;box-shadow:0 1px 3px rgba(0,0,0,.07)}.box h3{margin:0 0 10px;color:#1a7a3b;font-size:15px;display:flex;justify-content:space-between;align-items:center}.mem{border-bottom:1px solid #eee;padding:10px 0;display:flex;justify-content:space-between;align-items:center}.full{width:100%%;margin:7px 0;padding:13px;border:1px solid #ddd;border-radius:10px;font-size:14px}.big{width:100%%;background:#1a7a3b;padding:14px;color:#fff;border:0;border-radius:10px;font-weight:700;margin-top:10px;font-size:16px}.t{width:100%%;border-collapse:collapse;font-size:13px}.t th{background:#e8f5e9;padding:8px;text-align:left;color:#1a7a3b}.t td{padding:8px;border-bottom:1px solid #eee}.btn-o{width:100%%;padding:13px;border-radius:10px;border:0;font-weight:700;margin:7px 0;cursor:pointer;color:#fff;font-size:14px}</style></head><body><div class=h>Uwashusema Mahida - 2023 hadi 2035</div><div class=c><form method=GET action=/ class=r><select name=mwezi>%s</select><select name=mwaka>%s</select><button class=b>Sawa</button></form><div class=cards><div class=card style="background:#1a7a3b"><small>MAPATO TZS</small><b>%s</b></div><div class=card style="background:#d93025"><small>MATUMIZI TZS</small><b>%s</b></div><div class=card style="background:#1565c0"><small>Bakaa TZS</small><b>%s</b></div></div><div class=box><h3>Wanachama <button onclick="document.getElementById('ad').style.display='block'" class=b style="padding:6px 14px;font-size:13px">+ Mpya</button></h3><div id=ad style="display:none;margin-bottom:8px"><form method=POST action=/add class=r><input name=jina placeholder="Jina la Mwanachama" required class=full><button class=b>Hifadhi</button></form></div>`,mo,yo,fmtN(mapato),fmtN(matumizi),fmtN(baki))
for _,wm:=range list{fmt.Fprintf(w,`<div class=mem><span>%s</span><a href="/del?id=%d" onclick="return confirm('Futa?')" style="background:#fce4ec;padding:6px 10px;border-radius:8px;text-decoration:none;color:#c62828">X</a></div>`,html.EscapeString(wm.Jina),wm.ID)}
fmt.Fprintf(w,`</div><div class=box><h3>Ingiza Data</h3><form method=POST action=/save><input type=hidden name=mwezi value=%d><input type=hidden name=mwaka value=%d><select name=wanachama_id class=full required><option value="">Chagua Mwanachama</option>%s</select><select name=klasi class=full required><optgroup label="MAPATO"><option value="Hisa">Hisa</option><option value="Ada">Ada</option><option value="Marejesho ya Mkopo">Marejesho ya Mkopo</option><option value="Faini">Faini</option><option value="Fomu">Fomu</option><option value="Michango">Michango</option></optgroup><optgroup label="MATUMIZI"><option value="Mkopo kwa Mwanachama">Mkopo kwa Mwanachama</option><option value="Fedha ya Dharura">Fedha ya Dharura</option><option value="Mkono wa Pole">Mkono wa Pole</option><option value="Matumizi Mengine">Matumizi Mengine</option></optgroup></select><input name=maelezo placeholder="Maelezo (mf: mwezi wa...)" class=full><input name=kiasi type=number placeholder="Kiasi TZS" class=full required><select name=aina class=full id="aina"><option value=mapato>Mapato</option><option value=matumizi>Matumizi</option></select><button class=big>Hifadhi</button></form></div>`,mNum,mw,mop)
fmt.Fprint(w,`<div class=box><h3>Wanaorudisha Mikopo</h3><table class=t><tr><th>Jina</th><th>Kiasi</th><th>X</th></tr>`)
if db!=nil{
rows,_:=db.Query("SELECT w.jina,m.kiasi,m.id FROM mikopo_rudisha m JOIN wanachama w ON w.id=m.wanachama_id ORDER BY m.id DESC LIMIT 20")
if rows!=nil{
has:=false
for rows.Next(){has=true;var j string;var k int64;var id int;rows.Scan(&j,&k,&id);fmt.Fprintf(w,`<tr><td>%s</td><td>%s</td><td><a href="/delr?id=%d">X</a></td></tr>`,html.EscapeString(j),fmtN(k),id)}
if!has{fmt.Fprint(w,`<tr><td colspan=3 style="color:#999">Hakuna</td></tr>`)}
rows.Close()
}
}
fmt.Fprint(w,`</table><form method=POST action=/rudisha class=r style="margin-top:10px"><select name=wanachama_id required><option value="">Chagua</option>`)
for _,wm:=range list{fmt.Fprintf(w,`<option value=%d>%s</option>`,wm.ID,html.EscapeString(wm.Jina))}
fmt.Fprint(w,`</select><input name=kiasi type=number placeholder="Kiasi" required><button class=b>Rudisha</button></form></div><div class=box><h3>Wenye Deni la Mkopo</h3><table class=t><tr><th>Jina</th><th>Chukua</th><th>Rudisha</th><th>Baki</th></tr>`)
if db!=nil{
hasD:=false
for _,wm:=range list{
var ch,ru int64
db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mikopo_chukua WHERE wanachama_id=$1",wm.ID).Scan(&ch)
db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mikopo_rudisha WHERE wanachama_id=$1",wm.ID).Scan(&ru)
if ch>0{hasD=true;fmt.Fprintf(w,`<tr><td>%s</td><td>%s</td><td>%s</td><td style="color:#c62828;font-weight:700">%s</td></tr>`,html.EscapeString(wm.Jina),fmtN(ch),fmtN(ru),fmtN(ch-ru))}
}
if!hasD{fmt.Fprint(w,`<tr><td colspan=4 style="color:#999">Hakuna</td></tr>`)}
}
fmt.Fprint(w,`</table><form method=POST action=/chukua class=r style="margin-top:10px"><select name=wanachama_id required><option value="">Chagua</option>`)
for _,wm:=range list{fmt.Fprintf(w,`<option value=%d>%s</option>`,wm.ID,html.EscapeString(wm.Jina))}
fmt.Fprint(w,`</select><input name=kiasi type=number placeholder="Mkopo" required><button class=b>Weka</button></form></div><div class=box><h3>Historia ya Mwezi Huu</h3>`)
if db!=nil{
rows,_:=db.Query("SELECT w.jina,m.klasi,m.kiasi,TO_CHAR(m.tarehe,'DD/MM') FROM mahida m LEFT JOIN wanachama w ON w.id=m.wanachama_id WHERE m.mwezi=$1 AND m.mwaka=$2 ORDER BY m.id DESC LIMIT 30",mNum,mw)
if rows!=nil{
hasH:=false
fmt.Fprint(w,`<table class=t><tr><th>Jina</th><th>Klasi</th><th>Kiasi</th><th>Tarehe</th></tr>`)
for rows.Next(){hasH=true;var j,kl,ta string;var k int64;rows.Scan(&j,&kl,&k,&ta);fmt.Fprintf(w,`<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,html.EscapeString(j),html.EscapeString(kl),fmtN(k),ta)}
fmt.Fprint(w,`</table>`)
if!hasH{fmt.Fprint(w,`<div style="color:#999">Hakuna mwezi huu</div>`)}
rows.Close()
}
}
fmt.Fprintf(w,`</div><div class=box><h3>Hifadhi ya Kudumu & Tuma</h3><a href=/backup style="text-decoration:none"><button class=btn-o style="background:#ff6f00">💾 Pakua Backup ya Kudumu (JSON)</button></a><a href=/csv style="text-decoration:none"><button class=btn-o style="background:#546e7a">📊 Pakua Excel CSV</button></a><button class=btn-o style="background:#1a7a3b" onclick="tumaWA()">📄 Tuma WhatsApp Ripoti</button></div></div><script>function tumaWA(){var txt="*Uwashusema Mahida Ripoti*%%0A%%0AMAPATO: %s%%0AMATUMIZI: %s%%0ABAKAA: %s%%0AMwezi: %s %s";window.open("https://wa.me/?text="+txt,"_blank")} document.querySelector('select[name=klasi]').addEventListener('change',function(){var v=this.value;var mapato=["Hisa","Ada","Marejesho ya Mkopo","Faini","Fomu","Michango"];if(mapato.includes(v)){document.getElementById('aina').value='mapato'}else{document.getElementById('aina').value='matumizi'}})</script></body></html>`,fmtN(mapato),fmtN(matumizi),fmtN(baki),mz,mwS)
}
func addM(w http.ResponseWriter,r *http.Request){r.ParseForm();j:=r.FormValue("jina");if j!=""&&db!=nil{db.Exec("INSERT INTO wanachama(jina) VALUES ($1)",j)};http.Redirect(w,r,"/",302)}
func delM(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("id");if db!=nil{db.Exec("DELETE FROM wanachama WHERE id=$1",id)};http.Redirect(w,r,"/",302)}
func save(w http.ResponseWriter,r *http.Request){r.ParseForm();m,_:=strconv.Atoi(r.FormValue("mwezi"));y,_:=strconv.Atoi(r.FormValue("mwaka"));k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64);wid,_:=strconv.Atoi(r.FormValue("wanachama_id"));klasi:=r.FormValue("klasi");aina:=r.FormValue("aina");if strings.Contains(klasi,"Mkopo kwa")||strings.Contains(klasi,"Dharura")||strings.Contains(klasi,"Pole")||strings.Contains(klasi,"Mengine"){aina="matumizi"}else{if klasi=="Mkopo kwa Mwanachama"{aina="matumizi"}else{}};if db!=nil{db.Exec("INSERT INTO mahida(wanachama_id,mwezi,mwaka,hisa,klasi,maelezo,kiasi,aina) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",wid,m,y,klasi,klasi,r.FormValue("maelezo"),k,aina)};http.Redirect(w,r,"/",302)}
func chukua(w http.ResponseWriter,r *http.Request){r.ParseForm();wid,_:=strconv.Atoi(r.FormValue("wanachama_id"));k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64);if db!=nil{db.Exec("INSERT INTO mikopo_chukua(wanachama_id,kiasi) VALUES ($1,$2)",wid,k)};http.Redirect(w,r,"/",302)}
func rudisha(w http.ResponseWriter,r *http.Request){r.ParseForm();wid,_:=strconv.Atoi(r.FormValue("wanachama_id"));k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64);if db!=nil{db.Exec("INSERT INTO mikopo_rudisha(wanachama_id,kiasi) VALUES ($1,$2)",wid,k)};http.Redirect(w,r,"/",302)}
func delr(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("id");if db!=nil{db.Exec("DELETE FROM mikopo_rudisha WHERE id=$1",id)};http.Redirect(w,r,"/",302)}
func doBackup(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","application/json");w.Header().Set("Content-Disposition","attachment; filename=uwashusema-backup.json");type D struct{Wanachama []map[string]interface{} `json:"wanachama"`;Mahida []map[string]interface{} `json:"mahida"`};var d D;if db!=nil{rows,_:=db.Query("SELECT id,jina FROM wanachama");if rows!=nil{defer rows.Close();for rows.Next(){var id int;var j string;rows.Scan(&id,&j);d.Wanachama=append(d.Wanachama,map[string]interface{}{"id":id,"jina":j})}};rows,_=db.Query("SELECT wanachama_id,mwezi,mwaka,klasi,maelezo,kiasi,aina FROM mahida");if rows!=nil{defer rows.Close();for rows.Next(){var wid,m,yr int;var kl,ma,ai string;var k int64;rows.Scan(&wid,&m,&yr,&kl,&ma,&k,&ai);d.Mahida=append(d.Mahida,map[string]interface{}{"wid":wid,"m":m,"yr":yr,"kl":kl,"ma":ma,"k":k,"aina":ai})}}};json.NewEncoder(w).Encode(d)}
func doCSV(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/csv");w.Header().Set("Content-Disposition","attachment; filename=uwashusema.csv");cw:=csv.NewWriter(w);cw.Write([]string{"Jina","Mwezi","Mwaka","Klasi","Maelezo","Kiasi","Aina"});if db!=nil{rows,_:=db.Query("SELECT w.jina,m.mwezi,m.mwaka,m.klasi,m.maelezo,m.kiasi,m.aina FROM mahida m LEFT JOIN wanachama w ON w.id=m.wanachama_id ORDER BY m.id DESC");if rows!=nil{defer rows.Close();for rows.Next(){var j,kl,ma,ai string;var m,yr int;var k int64;rows.Scan(&j,&m,&yr,&kl,&ma,&k,&ai);cw.Write([]string{j,strconv.Itoa(m),strconv.Itoa(yr),kl,ma,strconv.FormatInt(k,10),ai})}}};cw.Flush()}