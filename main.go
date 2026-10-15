package main
import (
"context";"database/sql";"encoding/csv";"fmt";"net/http";"net/url";"os";"time"
_ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
func init(){
 u:=os.Getenv("DATABASE_URL"); if u==""{return}
 db,_=sql.Open("pgx",u)
 ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
 db.PingContext(ctx)
 db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)
}
func main(){
 http.HandleFunc("/",home)
 http.HandleFunc("/ongeza",ongeza)
 http.HandleFunc("/mwanachama",mwanachama)
 http.HandleFunc("/export",exportCSV)
 http.HandleFunc("/delete",del)
 p:=os.Getenv("PORT"); if p==""{p="10000"}
 http.ListenAndServe(":"+p,nil)
}
func getTotals() (mapato, dharura, matumizi, mikopo int64, br map[string]int64){
 br=map[string]int64{"hisa_anzia":0,"ada_mwezi":0,"marejesho":0,"faini":0,"fedha_fomu":0,"michango_mingine":0}
 if db==nil{return}
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')`).Scan(&mapato)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'`).Scan(&dharura)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'`).Scan(&matumizi)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'`).Scan(&mikopo)
 for k := range br { var v int64; db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina=$1`,k).Scan(&v); br[k]=v }
 return
}
func home(w http.ResponseWriter,r *http.Request){
 mapato,dharura,matumizi,mikopo,br:=getTotals()
 bakaa:=(mapato+dharura)-(matumizi+mikopo)
 waText:=fmt.Sprintf("UWASHUSEMA MAHIDA RIPOTI\nBAKAA: TZS %d\nMAPATO: %d\n- Hisa Anzia:%d\n- Ada Mwezi:%d\n- Marejesho:%d\n- Faini:%d\n- Fomu:%d\n- Mengine:%d\nDHARURA: %d\nMATUMIZI: %d\nMIKOPO: %d\nhttps://platform-app-cj1v.onrender.com",bakaa,mapato,br["hisa_anzia"],br["ada_mwezi"],br["marejesho"],br["faini"],br["fedha_fomu"],br["michango_mingine"],dharura,matumizi,mikopo)
 waLink:="https://wa.me/?text="+url.QueryEscape(waText)
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 fmt.Fprintf(w,`
<!DOCTYPE html><html><head><meta name=viewport content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title>
<style>
body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:0 10px 80px}
.header{padding:18px 6px 8px;font-weight:800;font-size:20px;color:#0f172a}
.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px;box-shadow:0 4px 12px #0000000d}
.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}
.big{font-size:28px;font-weight:900;margin:4px 0}
.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.small{font-size:11px;opacity:.8}
.big2{font-size:18px;font-weight:800}
.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700;margin:6px 0}
.btn-green{background:#0f7a3a;color:#fff}
.btn-white{background:#fff;border:1px solid #ddd;color:#222}
.btn-wa{background:#25D366;color:#fff}
input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0;background:#fafafa}
table{width:100%%;border-collapse:collapse;font-size:12px}
th{text-align:left;padding:8px;background:#f6f6f6}
td{padding:8px;border-bottom:1px solid #f0f0f0}
.badge{padding:3px 7px;border-radius:20px;background:#eef2ee;font-size:10px}
</style></head><body>
<div class=header>Uwashusema Mahida</div>
<div class="card bakaa"><div class=small>BAKAA - Pesa Iliyobaki Baada ya Kukopesha</div><div class=big>TZS %d</div><div class=small>Mapato (%d) + Dharura (%d) - Matumizi (%d) - Mikopo (%d)</div></div>
<div class=grid2>
<div class=card><div class=small>MAPATO (Jumla)</div><div class=big2>%d</div></div>
<div class=card><div class=small>FEDHA YA DHARURA</div><div class=big2>%d</div></div>
<div class=card><div class=small>MATUMIZI</div><div class=big2 style="color:#d32f2f">%d</div></div>
<div class=card><div class=small>MIKOPO ILIYOTOLEWA</div><div class=big2 style="color:#ef6c00">%d</div></div>
</div>
<div class=card>
<a href=/export><button class="btn btn-white">📥 Pakua Excel (CSV)</button></a>
<a href="%s" target=_blank><button class="btn btn-wa">💬 Tuma Ripoti WhatsApp</button></a>
</div>
<div class=card><h3 style="margin:0 0 6px">+ Ongeza Muamala</h3>
<form method=POST action=/ongeza>
<input name=jina placeholder="Jina la Mwanachama" required>
<input name=kiasi type=number placeholder="Kiasi TZS" required>
<select name=aina required>
<option value="" disabled selected>Chagua Aina ya Muamala</option>
<option value=hisa_anzia>Hisa Anzia</option>
<option value=ada_mwezi>Ada ya Mwezi</option>
<option value=marejesho>Marejesho</option>
<option value=faini>Faini</option>
<option value=fedha_fomu>Fedha ya Fomu</option>
<option value=michango_mingine>Michango Mingine</option>
<option value=fedha_dharura>Fedha ya Dharura</option>
<option value=matumizi>Matumizi</option>
<option value=mkopo>Mkopo (Kukopesha)</option>
</select>
<button class="btn btn-green">Ongeza</button>
</form></div>
<div class=card><h3 style="margin:0 0 6px">👤 Mwanachama Mpya</h3>
<form method=POST action=/mwanachama><input name=jina placeholder="Jina Kamili" required><input name=simu placeholder="Simu"><input name=namba placeholder="Namba"><button class="btn btn-green">Sajili</button></form></div>
<div class="card"><h3>📜 Historia</h3><table><tr><th>Tarehe</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th></th></tr>
`,bakaa,mapato,dharura,matumizi,mikopo,mapato,dharura,matumizi,mikopo,waLink)
 if db!=nil{
  rows,_:=db.Query(`SELECT id,jina,kiasi,aina,TO_CHAR(tarehe,'DD/MM HH:MI') FROM miamala ORDER BY id DESC LIMIT 100`)
  if rows!=nil{for rows.Next(){var id,k int64;var j,a,t string;rows.Scan(&id,&j,&k,&a,&t);fmt.Fprintf(w,`<tr><td>%s</td><td>%s</td><td><span class=badge>%s</span></td><td>%d</td><td><a href=/delete?id=%d style="color:red">X</a></td></tr>`,t,j,a,k,id)}}
 }
 fmt.Fprint(w,`</table></div></body></html>`)
}
func mwanachama(w http.ResponseWriter,r *http.Request){r.ParseForm();if db!=nil{db.Exec("INSERT INTO wanachama(jina,simu,namba)VALUES($1,$2,$3)",r.FormValue("jina"),r.FormValue("simu"),r.FormValue("namba"))};http.Redirect(w,r,"/",302)}
func ongeza(w http.ResponseWriter,r *http.Request){r.ParseForm();if db!=nil{db.Exec("INSERT INTO miamala(jina,kiasi,aina)VALUES($1,$2,$3)",r.FormValue("jina"),r.FormValue("kiasi"),r.FormValue("aina"))};http.Redirect(w,r,"/",302)}
func del(w http.ResponseWriter,r *http.Request){if db!=nil{db.Exec("DELETE FROM miamala WHERE id=$1",r.URL.Query().Get("id"))};http.Redirect(w,r,"/",302)}
func exportCSV(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","text/csv")
 w.Header().Set("Content-Disposition","attachment; filename=Uwashusema_Mahida_Ripoti.csv")
 writer:=csv.NewWriter(w)
 writer.Write([]string{"Tarehe","Jina","Aina","Kiasi"})
 if db!=nil{
  rows,_:=db.Query(`SELECT TO_CHAR(tarehe,'YYYY-MM-DD HH:MI'),jina,aina,kiasi FROM miamala ORDER BY id DESC`)
  if rows!=nil{for rows.Next(){var t,j,a string;var k int64;rows.Scan(&t,&j,&a,&k);writer.Write([]string{t,j,a,fmt.Sprintf("%d",k)})}}
  mapato,dharura,matumizi,mikopo,_:=getTotals();bakaa:=(mapato+dharura)-(matumizi+mikopo)
  writer.Write([]string{""});writer.Write([]string{"BAKAA",fmt.Sprintf("%d",bakaa)});writer.Write([]string{"MAPATO",fmt.Sprintf("%d",mapato)});writer.Write([]string{"DHARURA",fmt.Sprintf("%d",dharura)});writer.Write([]string{"MATUMIZI",fmt.Sprintf("%d",matumizi)});writer.Write([]string{"MIKOPO",fmt.Sprintf("%d",mikopo)})
 }
 writer.Flush()
}
