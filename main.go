package main
import (
"context";"database/sql";"fmt";"net/http";"os";"time"
_ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
func init(){
 u:=os.Getenv("DATABASE_URL")
 if u==""{return}
 db,_=sql.Open("pgx",u)
 ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
 db.PingContext(ctx)
 db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,wanachama_id INT,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)
}
func main(){
 http.HandleFunc("/",home)
 http.HandleFunc("/ongeza",ongeza)
 http.HandleFunc("/mwanachama",mwanachama)
 http.HandleFunc("/wanachama",wanachamaList)
 http.HandleFunc("/history",history)
 http.HandleFunc("/delete",del)
 p:=os.Getenv("PORT"); if p==""{p="10000"}
 http.ListenAndServe(":"+p,nil)
}
func home(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 var mapato, dharura, matumizi, mikopo int64
 if db!=nil{
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')`).Scan(&mapato)
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'`).Scan(&dharura)
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'`).Scan(&matumizi)
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'`).Scan(&mikopo)
 }
 bakaa := (mapato + dharura) - (matumizi + mikopo)
 fmt.Fprintf(w,`
<!DOCTYPE html><html><head><meta name=viewport content="width=device-width,initial-scale=1"><title>Uwashusema V3</title>
<style>*{box-sizing:border-box}body{font-family:sans-serif;background:#eef2ee;margin:0;padding:10px}
.card{background:#fff;padding:14px;border-radius:14px;margin-bottom:12px;box-shadow:0 2px 8px #0001}
.grid{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.big{font-size:22px;font-weight:800}
input,select,button{width:100%%;padding:12px;margin:6px 0;border-radius:10px;border:1px solid #ccc}
button{background:#0d7a3a;color:#fff;font-weight:700;border:0}
.badge{display:inline-block;padding:4px 8px;border-radius:20px;background:#eee;font-size:11px;margin:2px}
.bakaa{background:#0d7a3a;color:#fff}
</style></head><body>
<div class="card bakaa"><small>BAKAA - Pesa Iliyobaki Baada ya Kukopesha</small><div class="big">TZS %d</div><small>Mapato (%d) + Dharura (%d) - Matumizi (%d) - Mikopo (%d)</small></div>
<div class="grid">
<div class="card"><small>MAPATO (Jumla)</small><div class="big">%d</div><small>Hisa+Ada+Marejesho+Faini+Fomu+Michango</small></div>
<div class="card"><small>FEDHA YA DHARURA</small><div class="big">%d</div></div>
<div class="card"><small>MATUMIZI</small><div class="big" style="color:red">%d</div></div>
<div class="card"><small>MIKOPO ILIYOTOLEWA</small><div class="big" style="color:#b45309">%d</div></div>
</div>

<div class="card"><h3>➕ Ongeza Muamala</h3>
<form method=POST action=/ongeza>
<input name=jina placeholder="Jina la Mwanachama" required>
<input name=kiasi type=number placeholder="Kiasi TZS" required>
<select name=aina required>
<optgroup label="MAPATO">
<option value=hisa_anzia>a) Hisa Anzia</option>
<option value=ada_mwezi>b) Ada ya Mwezi</option>
<option value=marejesho>c) Marejesho (Rudisho la Mkopo)</option>
<option value=faini>d) Faini</option>
<option value=fedha_fomu>e) Fedha ya Fomu</option>
<option value=michango_mingine>f) Michango Mingine</option>
</optgroup>
<optgroup label="VINGINE">
<option value=fedha_dharura>FEDHA YA DHARURA (Peke Yake)</option>
<option value=matumizi>MATUMIZI (Peke Yake)</option>
<option value=mkopo>MKOPO - Kukopesha (itaondoa kwenye Bakaa)</option>
</optgroup>
</select>
<button>Ongeza</button>
</form></div>

<div class="card"><h3>👤 Mwanachama Mpya</h3>
<form method=POST action=/mwanachama>
<input name=jina placeholder="Jina Kamili" required><input name=simu placeholder="Simu"><input name=namba placeholder="Namba"><button>Sajili</button>
</form>
<a href=/wanachama><button style="background:#eee;color:#333;margin-top:8px">👥 Orodha ya Wanachama</button></a>
<a href=/history><button style="background:#eee;color:#333">📜 Historia Yote</button></a>
</div>
</body></html>
`,bakaa,mapato,dharura,matumizi,mikopo,mapato,dharura,matumizi,mikopo)
}
func mwanachama(w http.ResponseWriter,r *http.Request){r.ParseForm();db.Exec("INSERT INTO wanachama(jina,simu,namba)VALUES($1,$2,$3)",r.FormValue("jina"),r.FormValue("simu"),r.FormValue("namba"));http.Redirect(w,r,"/",302)}
func ongeza(w http.ResponseWriter,r *http.Request){r.ParseForm();db.Exec("INSERT INTO miamala(jina,kiasi,aina)VALUES($1,$2,$3)",r.FormValue("jina"),r.FormValue("kiasi"),r.FormValue("aina"));http.Redirect(w,r,"/",302)}
func wanachamaList(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html; charset=utf-8");fmt.Fprint(w,`<a href=/>⬅ Rudi</a><h2>Wanachama</h2><table border=1 width=100%><tr><th>Jina</th><th>Simu</th><th>Namba</th></tr>`);rows,_:=db.Query("SELECT jina,simu,namba FROM wanachama ORDER BY id DESC");for rows.Next(){var j,s,n sql.NullString;rows.Scan(&j,&s,&n);fmt.Fprintf(w,"<tr><td>%s</td><td>%s</td><td>%s</td></tr>",j.String,s.String,n.String)};fmt.Fprint(w,"</table>")}
func history(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html; charset=utf-8");fmt.Fprint(w,`<a href=/>⬅ Rudi</a><h2>Historia</h2><table border=1 width=100%><tr><th>Tarehe</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Futa</th></tr>`);rows,_:=db.Query("SELECT id,jina,kiasi,aina,TO_CHAR(tarehe,'DD/MM HH:MI') FROM miamala ORDER BY id DESC LIMIT 200");for rows.Next(){var id,k int64;var j,a,t sql.NullString;rows.Scan(&id,&j,&k,&a,&t);fmt.Fprintf(w,"<tr><td>%s</td><td>%s</td><td><span class=badge>%s</span></td><td>%d</td><td><a href=/delete?id=%d>X</a></td></tr>",t.String,j.String,a.String,k,id)};fmt.Fprint(w,"</table>")}
func del(w http.ResponseWriter,r *http.Request){db.Exec("DELETE FROM miamala WHERE id=$1",r.URL.Query().Get("id"));http.Redirect(w,r,"/history",302)}
