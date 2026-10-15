package main

import (
 "database/sql"
 "encoding/csv"
 "fmt"
 "log"
 "net/http"
 "net/url"
 "os"
 "time"
 _ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
 if u := os.Getenv("DATABASE_URL"); u != "" {
  db, _ = sql.Open("pgx", u)
  db.Exec("CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)")
  db.Exec("CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe DATE DEFAULT CURRENT_DATE)")
 }
 http.HandleFunc("/", home)
 http.HandleFunc("/ongeza", ongeza)
 http.HandleFunc("/mwanachama", mwanachama)
 http.HandleFunc("/export", exportCSV)
 http.HandleFunc("/delete", del)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/futa-zote", futaZote)
 p := os.Getenv("PORT")
 if p == "" { p = "10000" }
 log.Println("V7 start", p)
 http.ListenAndServe("0.0.0.0:"+p, nil)
}

func getTotals() (mapato, dharura, matumizi, mikopo int64) {
 if db == nil { return }
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')").Scan(&mapato)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'").Scan(&dharura)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'").Scan(&matumizi)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'").Scan(&mikopo)
 return
}

func home(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-cache")
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 mapato, dharura, matumizi, mikopo := getTotals()
 bakaa := (mapato + dharura) - (matumizi + mikopo)
 mwaka := time.Now().Year()
 waLink := "https://wa.me/?text=" + url.QueryEscape(fmt.Sprintf("UWASHUSEMA MAHIDA %d BAKAA TZS %d", mwaka, bakaa))
 msg := ""
 if r.URL.Query().Get("futa") == "ok" { msg = "<div class=card style=background:#ffebee>✅ Data zote zimefutwa</div>" }
 if r.URL.Query().Get("edit") == "ok" { msg = "<div class=card style=background:#e8f5e9>✅ Imehaririwa vizuri</div>" }

 today := time.Now().Format("2006-01-02")

 fmt.Fprint(w, "<!DOCTYPE html><html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><title>Uwashusema</title>")
 fmt.Fprint(w, "<style>body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:0 10px 90px}.header{padding:18px 6px 8px;font-weight:800;font-size:20px}.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px}.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}.big{font-size:26px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}.small{font-size:11px;opacity:.8}.big2{font-size:18px;font-weight:800}.btn{width:100%;padding:12px;border-radius:12px;border:0;font-weight:700;margin:6px 0}.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}.btn-red{border:1px solid #ef5350;color:#ef5350;background:#fff}input,select{width:100%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0}table{width:100%;border-collapse:collapse;font-size:12px}th{padding:8px;background:#f6f6f6;text-align:left}td{padding:8px;border-bottom:1px solid #eee}.badge{background:#eef2ee;padding:3px 7px;border-radius:20px;font-size:10px}.act{padding:6px 10px;border-radius:8px;font-size:12px;text-decoration:none;margin-right:5px;display:inline-block}.act-edit{background:#1976d2;color:#fff}.act-del{background:#d32f2f;color:#fff}.year{font-size:12px;background:#00000022;padding:2px 8px;border-radius:10px}</style>")
 fmt.Fprint(w, "<script>function checkAina(){var a=document.getElementById('aina').value;var j=document.getElementById('jina');if(a=='fedha_dharura'||a=='matumizi'||a=='mkopo'){j.value='Washusema Nzima';j.readOnly=true}else{j.readOnly=false;if(j.value=='Washusema Nzima')j.value='';}}</script></head><body>")
 fmt.Fprintf(w, "<div class=header>Uwashusema Mahida <span class=year>%d</span></div>%s", mwaka, msg)
 fmt.Fprintf(w, "<div class=\"card bakaa\"><div class=small>BAKAA %d - Pesa Iliyobaki</div><div class=big>TZS %d</div><div class=small>Mapato %d + Dharura %d - Matumizi %d - Mikopo %d</div></div>", mwaka, bakaa, mapato, dharura, matumizi, mikopo)
 fmt.Fprintf(w, "<div class=grid2><div class=card><div class=small>MAPATO</div><div class=big2>%d</div></div><div class=card><div class=small>DHARURA</div><div class=big2>%d</div></div><div class=card><div class=small>MATUMIZI</div><div class=big2 style=color:#d32f2f>%d</div></div><div class=card><div class=small>MIKOPO</div><div class=big2 style=color:#ef6c00>%d</div></div></div>", mapato, dharura, matumizi, mikopo)
 fmt.Fprintf(w, "<div class=card><a href=/export><button class=\"btn btn-white\">📥 Pakua Excel</button></a><a href=\"%s\" target=_blank><button class=\"btn btn-wa\">💬 WhatsApp</button></a><a href=/futa-zote onclick=\"return confirm('Futa DATA YOTE?')\"><button class=\"btn btn-red\">🗑️ Futa Data Zote</button></a></div>", waLink)
 fmt.Fprintf(w, "<div class=card><h3>+ Ongeza Muamala</h3><form method=POST action=/ongeza><input id=jina name=jina placeholder=\"Jina la Mwanachama\" required><input name=kiasi type=number placeholder=\"Kiasi TZS\" required><input type=date name=tarehe value=\"%s\" required><select id=aina name=aina required onchange=checkAina()><option value=\"\" disabled selected>Chagua Aina</option><option value=hisa_anzia>Hisa Anzia</option><option value=ada_mwezi>Ada ya Mwezi</option><option value=marejesho>Marejesho</option><option value=faini>Faini</option><option value=fedha_fomu>Fedha ya Fomu</option><option value=michango_mingine>Michango Mingine</option><option value=fedha_dharura>FEDHA YA DHARURA (Kikundi)</option><option value=matumizi>MATUMIZI (Kikundi)</option><option value=mkopo>MKOPO (Kikundi)</option></select><button class=\"btn btn-green\">Ongeza Sasa</button></form></div>", today)
 fmt.Fprint(w, "<div class=card><h3>👤 Mwanachama Mpya</h3><form method=POST action=/mwanachama><input name=jina placeholder=\"Jina Kamili\" required><input name=simu placeholder=\"Simu\"><input name=namba placeholder=\"Namba\"><button class=\"btn btn-green\">Sajili Mwanachama</button></form></div>")
 fmt.Fprint(w, "<div class=card><h3>📜 Historia - Tarehe Kamili</h3><table><tr><th>Tarehe</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Action</th></tr>")
 if db != nil {
  rows, _ := db.Query("SELECT id,jina,kiasi,aina,TO_CHAR(tarehe,'DD/MM/YYYY') FROM miamala ORDER BY tarehe DESC, id DESC LIMIT 200")
  if rows != nil {
   for rows.Next() {
    var id, k int64
    var j, a, t string
    rows.Scan(&id, &j, &k, &a, &t)
    fmt.Fprintf(w, "<tr><td><b>%s</b></td><td>%s</td><td><span class=badge>%s</span></td><td>%d</td><td><a class=act act-edit href=\"/edit?id=%d\">EDIT</a><a class=act act-del href=\"/delete?id=%d\" onclick=\"return confirm('Futa?')\">X</a></td></tr>", t, j, a, k, id, id)
   }
  }
 }
 fmt.Fprint(w, "</table></div></body></html>")
}

func mwanachama(w http.ResponseWriter, r *http.Request) {
 r.ParseForm()
 if db != nil { db.Exec("INSERT INTO wanachama(jina,simu,namba)VALUES($1,$2,$3)", r.FormValue("jina"), r.FormValue("simu"), r.FormValue("namba")) }
 http.Redirect(w, r, "/?add=ok", 302)
}
func ongeza(w http.ResponseWriter, r *http.Request) {
 r.ParseForm()
 jina := r.FormValue("jina")
 aina := r.FormValue("aina")
 tarehe := r.FormValue("tarehe")
 if aina == "fedha_dharura" || aina == "matumizi" || aina == "mkopo" { jina = "Washusema Nzima" }
 if db != nil {
  if tarehe != "" {
   db.Exec("INSERT INTO miamala(jina,kiasi,aina,tarehe)VALUES($1,$2,$3,$4::date)", jina, r.FormValue("kiasi"), aina, tarehe)
  } else {
   db.Exec("INSERT INTO miamala(jina,kiasi,aina)VALUES($1,$2,$3)", jina, r.FormValue("kiasi"), aina)
  }
 }
 http.Redirect(w, r, "/?add=ok", 302)
}
func del(w http.ResponseWriter, r *http.Request) {
 if db != nil { db.Exec("DELETE FROM miamala WHERE id=$1", r.URL.Query().Get("id")) }
 http.Redirect(w, r, "/", 302)
}
func futaZote(w http.ResponseWriter, r *http.Request) {
 if db != nil { db.Exec("DELETE FROM miamala"); db.Exec("DELETE FROM wanachama") }
 http.Redirect(w, r, "/?futa=ok", 302)
}
func edit(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Query().Get("id")
 if id == "" { id = r.FormValue("id") }
 if r.Method == "POST" {
  if db != nil { db.Exec("UPDATE miamala SET jina=$1,kiasi=$2,aina=$3,tarehe=$4::date WHERE id=$5", r.FormValue("jina"), r.FormValue("kiasi"), r.FormValue("aina"), r.FormValue("tarehe"), r.FormValue("id")) }
  http.Redirect(w, r, "/?edit=ok", 302)
  return
 }
 var jina, aina, tarehe string
 var kiasi int64
 if db != nil { db.QueryRow("SELECT jina,kiasi,aina,TO_CHAR(tarehe,'YYYY-MM-DD') FROM miamala WHERE id=$1", id).Scan(&jina, &kiasi, &aina, &tarehe) }
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 fmt.Fprintf(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;padding:12px}.card{background:#fff;border-radius:16px;padding:18px}input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:8px 0}.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700}.btn-green{background:#0f7a3a;color:#fff}</style></head><body><div class=card><h3>✏️ Hariri Muamala #%s</h3><form method=POST action=/edit><input type=hidden name=id value=\"%s\"><label>Jina</label><input name=jina value=\"%s\" required><label>Kiasi</label><input name=kiasi type=number value=\"%d\" required><label>Tarehe (DD/MM/YYYY)</label><input type=date name=tarehe value=\"%s\" required><label>Aina</label><select name=aina required><option value=hisa_anzia>Hisa</option><option value=ada_mwezi>Ada</option><option value=marejesho>Marejesho</option><option value=faini>Faini</option><option value=fedha_fomu>Fomu</option><option value=michango_mingine>Mingine</option><option value=fedha_dharura>Dharura</option><option value=matumizi>Matumizi</option><option value=mkopo>Mkopo</option></select><button class=\"btn btn-green\">Hifadhi Mabadiliko</button></form><br><a href=\"/\">Rudi Nyuma</a></div></body></html>", id, id, jina, kiasi, tarehe)
}
func exportCSV(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "text/csv")
 w.Header().Set("Content-Disposition", "attachment; filename=Uwashusema.csv")
 writer := csv.NewWriter(w)
 writer.Write([]string{"Tarehe", "Jina", "Aina", "Kiasi"})
 if db != nil {
  rows, _ := db.Query("SELECT TO_CHAR(tarehe,'YYYY-MM-DD'),jina,aina,kiasi FROM miamala ORDER BY tarehe DESC")
  if rows != nil {
   for rows.Next() {
    var t, j, a string
    var k int64
    rows.Scan(&t, &j, &a, &k)
    writer.Write([]string{t, j, a, fmt.Sprintf("%d", k)})
   }
  }
 }
 writer.Flush()
}
