package main
import (
"database/sql";"encoding/csv";"fmt";"log";"net/http";"net/url";"os"
_ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
func main(){
 // Fungua DB bila ku-block port
 if u:=os.Getenv("DATABASE_URL"); u!="" {
   var err error
   db, err = sql.Open("pgx", u)
   if err==nil {
     db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)`)
     db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)
   }
 }
 http.HandleFunc("/",home)
 http.HandleFunc("/ongeza",ongeza)
 http.HandleFunc("/mwanachama",mwanachama)
 http.HandleFunc("/export",exportCSV)
 http.HandleFunc("/delete",del)
 p:=os.Getenv("PORT"); if p==""{p="10000"}
 log.Printf("Uwashusema Mahida inaanza kwenye port %s", p)
 // Hii ndio muhimu kwa Render - bind 0.0.0.0
 err:=http.ListenAndServe("0.0.0.0:"+p, nil)
 if err!=nil{log.Fatal(err)}
}
func getTotals() (mapato, dharura, matumizi, mikopo int64){
 if db==nil{return}
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')`).Scan(&mapato)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'`).Scan(&dharura)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'`).Scan(&matumizi)
 db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'`).Scan(&mikopo)
 return
}
func home(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Cache-Control","no-cache, no-store, must-revalidate")
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 mapato,dharura,matumizi,mikopo:=getTotals()
 bakaa:=(mapato+dharura)-(matumizi+mikopo)
 waText:=fmt.Sprintf("UWASHUSEMA MAHIDA RIPOTI\nBAKAA: TZS %d\nMAPATO: %d\nDHARURA: %d\nMATUMIZI: %d\nMIKOPO: %d",bakaa,mapato,dharura,matumizi,mikopo)
 waLink:="https://wa.me/?text="+url.QueryEscape(waText)
 fmt.Fprintf(w,`
<!DOCTYPE html><html><head><meta name=viewport content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title>
<style>body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:0 10px 80px}.header{padding:18px 6px 8px;font-weight:800;font-size:20px}.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px;box-shadow:0 4px 12px #0000000d}.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}.big{font-size:28px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}.small{font-size:11px;opacity:.8}.big2{font-size:18px;font-weight:800}.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700;margin:6px 0}.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0;background:#fafafa}table{width:100%%;border-collapse:collapse;font-size:12px}th{text-align:left;padding:8px;background:#f6f6f6}td{padding:8px;border-bottom:1px solid #f0f0f0}.badge{padding:3px 7px;border-radius:20px;background:#eef2ee;font-size:10px}.note{font-size:11px;background:#fff8e1;padding:8px;border-radius:8px;margin-top:6px;border-left:4px solid #ffb300}</style>
<script>function checkAina(){var a=document.getElementById('aina').value;var j=document.getElementById('jina');var n=document.getElementById('note');if(a=='fedha_dharura'||a=='matumizi'||a=='mkopo'){j.value='Washusema Nzima';j.readOnly=true;n.style.display='block';n.innerHTML='ℹ️ Hii ni ya Washusema Nzima (kikundi).';}else{j.readOnly=false;if(j.value=='Washusema Nzima')j.value='';n.style.display='none';}}</script>
</head><body>
<div class=header>Uwashusema Mahida</div>
<div class="card bakaa"><div class=small>BAKAA - Pesa Iliyobaki Baada ya Kukopesha</div><div class=big>TZS %d</div><div class=small>Mapato (%d) + Dharura (%d) - Matumizi (%d) - Mikopo (%d)</div></div>
<div class=grid2>
<div class=card><div class=small>MAPATO</div><div class=big2>%d</div></div>
<div class=card><div class=small>FEDHA YA DHARURA</div><div class=big2>%d</div></div>
<div class=card><div class=small>MATUMIZI</div><div class=big2 style="color:#d32f2f">%d</div></div>
<div class=card><div class=small>MIKOPO ILIYOTOLEWA</div><div class=big2 style="color:#ef6c00">%d</div></div>
</div>
<div class=card><a href=/export><button class="btn btn-white">📥 Pakua Excel (CSV)</button></a><a href="%s" target=_blank><button class="btn btn-wa">💬 Tuma Ripoti WhatsApp</button></a></div>
<div class=card><h3>+ Ongeza Muamala</h3><form method=POST action=/ongeza><input id=jina name=jina placeholder="Jina la Mwanachama" required><input name=kiasi type=number placeholder="Kiasi TZS" required><select id=aina name=aina required onchange=checkAina()><option value="" disabled selected>Chagua Aina</option><option value=hisa_anzia>Hisa Anzia</option><option value=ada_mwezi>Ada ya Mwezi</option><option value=marejesho>Marejesho</option><option value=faini>Faini</option><option value=fedha_fomu>Fedha ya Fomu</option><option value=michango_mingine>Michango Mingine</option><option value=fedha_dharura>--- FEDHA YA DHARURA (Ya Kikundi) ---</option><option value=matumizi>--- MATUMIZI (Ya Kikundi) ---</option><option value=mkopo>--- MKOPO (Ya Kikundi) ---</option></select><div id=note class=note style="display:none"></div><button class="btn btn-green">Ongeza Sasa</button></form></div>
<div class=card><h3>👤 Mwanachama Mpya</h3><form method=POST action=/mwanachama><input name=jina placeholder="Jina Kamili" required><input name=simu placeholder="Simu"><input name=namba placeholder="Namba"><button class="btn btn-green">Sajili</button></form></div>
<div class="card"><h3>📜 Historia Live</h3><table><tr><th>Tarehe</th><th>Jina / Kikundi</th><th>Aina</th><th>Kiasi</th><th></th></tr>
`,bakaa,mapato,dharura,matumizi,mikopo,mapato,dharura,matumizi,mikopo,waLink)
 if db!=nil{
  rows,_:=db.Query(`SELECT id,jina,kiasi,aina,TO_CHAR(tarehe,'DD/MM HH24:MI:SS') FROM miamala ORDER BY id DESC LIMIT 150`)
  if rows!=nil{for rows.Next(){var id,k int64;var j,a,t string;rows.Scan(&id,&j,&k,&a,&t);fmt.Fprintf(w,`<tr><td>%s</td><td>%s</td><td><span class=badge>%s</span></td><td>%d</td><td><a href=/delete?id=%d style="color:red">X</a></td></tr>`,t,j,a,k,id)}}
 }
 fmt.Fprint(w,`</table></div></body></html>`)
}
func mwanachama(w http.ResponseWriter,r *http.Request){r.ParseForm();if db!=nil{db.Exec("INSERT INTO wanachama(jina,simu,namba)VALUES($1,$2,$3)",r.FormValue("jina"),r.FormValue("simu"),r.FormValue("namba"))};http.Redirect(w,r,"/",302)}
func ongeza(w http.ResponseWriter,r *http.Request){r.ParseForm();jina:=r.FormValue("jina");aina:=r.FormValue("aina");if aina=="fedha_dharura"||aina=="matumizi"||aina=="mkopo"{jina="Washusema Nzima"};if db!=nil{db.Exec("INSERT INTO miamala(jina,kiasi,aina)VALUES($1,$2,$3)",jina,r.FormValue("kiasi"),aina)};http.Redirect(w,r,"/",302)}
func del(w http.ResponseWriter,r *http.Request){if db!=nil{db.Exec("DELETE FROM miamala WHERE id=$1",r.URL.Query().Get("id"))};http.Redirect(w,r,"/",302)}
func exportCSV(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","text/csv")
 w.Header().Set("Content-Disposition","attachment; filename=Uwashusema_Mahida.csv")
 writer:=csv.NewWriter(w)
 writer.Write([]string{"Tarehe","Jina/Kikundi","Aina","Kiasi"})
 if db!=nil{
  rows,_:=db.Query(`SELECT TO_CHAR(tarehe,'YYYY-MM-DD HH24:MI:SS'),jina,aina,kiasi FROM miamala ORDER BY id DESC`)
  if rows!=nil{for rows.Next(){var t,j,a string;var k int64;rows.Scan(&t,&j,&a,&k);writer.Write([]string{t,j,a,fmt.Sprintf("%d",k)})}}
  mapato,dharura,matumizi,mikopo:=getTotals();bakaa:=(mapato+dharura)-(matumizi+mikopo)
  writer.Write([]string{""});writer.Write([]string{"BAKAA",fmt.Sprintf("%d",bakaa)});writer.Write([]string{"MAPATO",fmt.Sprintf("%d",mapato)});writer.Write([]string{"DHARURA",fmt.Sprintf("%d",dharura)});writer.Write([]string{"MATUMIZI",fmt.Sprintf("%d",matumizi)});writer.Write([]string{"MIKOPO",fmt.Sprintf("%d",mikopo)})
 }
 writer.Flush()
}
