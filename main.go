package main
import (
"database/sql";"encoding/csv";"fmt";"log";"net/http";"net/url";"os";"strconv"
_ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
func main(){
 if u:=os.Getenv("DATABASE_URL"); u!="" {
   db,_=sql.Open("pgx", u)
   db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)`)
   db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)
 }
 http.HandleFunc("/",home)
 http.HandleFunc("/ongeza",ongeza)
 http.HandleFunc("/mwanachama",mwanachama)
 http.HandleFunc("/export",exportCSV)
 http.HandleFunc("/delete",del)
 http.HandleFunc("/edit",edit)
 http.HandleFunc("/futa-zote",futaZote)
 p:=os.Getenv("PORT"); if p==""{p="10000"}
 log.Printf("Uwashusema Mahida V5 - Edit Mode Port %s", p)
 http.ListenAndServe("0.0.0.0:"+p, nil)
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
 alert:=""
 if r.URL.Query().Get("futa")=="ok"{alert=`<div class=card style="background:#ffebee;border-left:4px solid red">✅ Data zote zimefutwa!</div>`}
 if r.URL.Query().Get("edit")=="ok"{alert=`<div class=card style="background:#e8f5e9;border-left:4px solid green">✅ Muamala umehaririwa!</div>`}
 if r.URL.Query().Get("add")=="ok"{alert=`<div class=card style="background:#e8f5e9;border-left:4px solid green">✅ Muamala umeongezwa!</div>`}
 fmt.Fprintf(w,`
<!DOCTYPE html><html><head><meta name=viewport content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title>
<style>body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:0 10px 90px}.header{padding:18px 6px 8px;font-weight:800;font-size:20px}.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px;box-shadow:0 4px 12px #0000000d}.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}.big{font-size:28px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}.small{font-size:11px;opacity:.8}.big2{font-size:18px;font-weight:800}.btn{width:100%%;padding:12px;border-radius:12px;border:0;font-weight:700;margin:6px 0;cursor:pointer}.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}.btn-red{background:#fff;border:1px solid #ef5350;color:#ef5350}input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0;background:#fafafa}table{width:100%%;border-collapse:collapse;font-size:12px}th{text-align:left;padding:8px;background:#f6f6f6}td{padding:8px;border-bottom:1px solid #f0f0f0}.badge{padding:3px 7px;border-radius:20px;background:#eef2ee;font-size:10px}.note{font-size:11px;background:#fff8e1;padding:8px;border-radius:8px;margin-top:6px;border-left:4px solid #ffb300}.act{padding:5px 8px;border-radius:8px;font-size:11px;text-decoration:none;margin-right:4px}.act-edit{background:#e3f2fd;color:#1565c0;border:1px solid #90caf9}.act-del{background:#ffebee;color:#c62828;border:1px solid #ef9a9a}</style>
<script>function checkAina(){var a=document.getElementById('aina').value;var j=document.getElementById('jina');var n=document.getElementById('note');if(a=='fedha_dharura'||a=='matumizi'||a=='mkopo'){j.value='Washusema Nzima';j.readOnly=true;n.style.display='block';n.innerHTML='ℹ️ Ya Washusema Nzima (kikundi).';}else{j.readOnly=false;if(j.value=='Washusema Nzima')j.value='';n.style.display='none';}}</script>
</head><body>
<div class=header>Uwashusema Mahida</div>
%s
<div class="card bakaa"><div class=small>BAKAA - Pesa Iliyobaki Baada ya Kukopesha</div><div class=big>TZS %d</div><div class=small>Mapato (%d) + Dharura (%d) - Matumizi (%d) - Mikopo (%d)</div></div>
<div class=grid2>
<div class=card><div class=small>MAPATO</div><div class=big2>%d</div></
