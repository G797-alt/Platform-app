package main

import (
 "database/sql"
 "encoding/csv"
 "fmt"
 "log"
 "net/http"
 "net/url"
 "os"
 "strconv"
 "time"
 _ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
 if u := os.Getenv("DATABASE_URL"); u != "" {
  var err error
  db, err = sql.Open("pgx", u)
  if err != nil { log.Println(err) } else {
   db.Exec("CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)")
   db.Exec("CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())")
  }
 }
 http.HandleFunc("/", home)
 http.HandleFunc("/ongeza", ongeza)
 http.HandleFunc("/mwanachama", mwanachama)
 http.HandleFunc("/export", exportCSV)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/delete", del)
 p := os.Getenv("PORT")
 if p == "" { p = "10000" }
 log.Println("V10 FIX start", p)
 http.ListenAndServe("0.0.0.0:"+p, nil)
}

func getTotals() (mapato, dharura, matumizi, mikopo int64, count int) {
 if db == nil { return }
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')").Scan(&mapato)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'").Scan(&dharura)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'").Scan(&matumizi)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'").Scan(&mikopo)
 db.QueryRow("SELECT COUNT(*) FROM miamala").Scan(&count)
 return
}

func home(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-cache")
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 mapato, dharura, matumizi, mikopo, count := getTotals()
 bakaa := (mapato + dharura) - (matumizi + mikopo)
 mwaka := time.Now().Year()
 waLink := "https://wa.me/?text=" + url.QueryEscape(fmt.Sprintf("UWASHUSEMA %d BAKAA %d", mwaka, bakaa))
 msg := ""
 if r.URL.Query().Get("edit") == "ok" { msg = "<div class=card style=background:#e8f5e9;padding:12px>✅ Imehaririwa!</div>" }
 if r.URL.Query().Get("add") == "ok" { msg = "<div class=card style=background:#e8f5e9;padding:12px>✅ Imehifadhiwa! Jumla Miamala: "+fmt.Sprintf("%d", count)+"</div>" }
 today := time.Now().Format("2006-01-02")

 fmt.Fprint(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:10px 10px 90px}.header{padding:16px 6px;font-weight:800;font-size:20px}.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px}.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}.big{font-size:28px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}.small{font-size:11px;opacity:.9}.big2{font-size:18px;font-weight:800}.btn{width:100%;padding:13px;border-radius:12px;border:0;font-weight:700;margin:6px 0}.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}input,select{width:100%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0}table{width:100%;border-collapse:collapse;font-size:12px}th{padding:8px;background:#222;color:#fff;text-align:left}td{padding:10px;border-bottom:1px solid #eee}.badge{background:#eef2ee;padding:3px 7px;border-radius:20px;font-size:10px}.editbtn{background:#1976d2;color:#fff;padding:8px 14px;border-radius:8px;text-decoration:none;font-weight:700;display:inline-block}</style><script>function checkAina(){var a=document.getElementById('aina').value;var j=document.getElementById('jina');if(a=='fedha_dharura'||a=='matumizi'||a=='mkopo'){j.value='Washusema Nzima';j.readOnly=true}else{j.readOnly=false;if(j.value=='Washusema Nzima')j.value='';}}</script></head><body>")
 fmt.Fprintf(w, "<div class=header>Uwashusema Mahida <span style=font-size:12px;background:#0002;padding:3px 8px;border-radius:10px>%d | Miamala %d</span></div>%s", mwaka, count, msg)
 fmt.Fprintf(w, "<div class=\"card bakaa\"><div class=small>BAKAA %d</div><div class=big>TZS %d</div><div class=small>Mapato %d + Dharura %d - Matumizi %d - Mikopo %d</div></div>", mwaka, bakaa, mapato, dharura, matumizi, mikopo)
 fmt.Fprintf(w, "<div class=grid2><div class=card><div class=small>MAPATO</div><div class=big2>%d</div></div><div class=card><div class=small>DHARURA</div><div class=big2>%d</div></div><div class=card><div class=small>MATUMIZI</div><div class=big2 style=color:red>%d</div></div><div class=card><div class=small>MIKOPO</div><div class=big2 style=color:orange>%d</div></div></div>", mapato, dharura, matumizi, mikopo)
 fmt.Fprintf(w, "<div class=card><a href=/export><button class=\"btn btn-white\">📥 Pakua Excel %d</button></a><a href=\"%s\" target=_blank><button class=\"btn btn-wa\">WhatsApp</button></a></div>", mwaka, waLink)
 fmt.Fprintf(w, "<div class=card><h3>+ Ongeza Muamala</h3><form method=POST action=/ongeza><input id=jina name=jina placeholder=\"Jina\" required><input name=kiasi type=number placeholder=\"Kiasi TZS\" required><input type=date name=tarehe value=\"%s\" required><select id=aina name=aina required onchange=checkAina()><option value=\"\" disabled selected>Chagua Aina</option><option value=hisa_anzia>Hisa Anzia</option><option value=ada_mwezi>Ada ya Mwezi</option><option value=marejesho>Marejesho</option><option value=faini>Faini</option><option value=fedha_fomu>Fedha ya Fomu</option><option value=michango_mingine>Michango Mingine</option><option value=fedha_dharura>FEDHA YA DHARURA</option><option value=matumizi>MATUMIZI</option><option value=mkopo>MKOPO</option></select><button class=\"btn btn-green\">HIFADHI - Itaonekana Chini</button></form></div>", today)
 fmt.Fprint(w, "<div class=card><h3>👤 Mwanachama Mpya</h3><form method=POST action=/mwanachama><input name=jina placeholder=\"Jina Kamili\" required><input name=simu placeholder=\"Simu\"><input name=namba placeholder=\"Namba\"><button class=\"btn btn-green\">Sajili</button></form></div>")

 // HISTORIA NA EDIT HALISI
 fmt.Fprintf(w, "<div class=card><h3>📜 Historia (%d) - Bonyeza EDIT kuhariri</h3><table><tr><th>ID</th><th>Tarehe</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>EDIT</th></tr>", count)
 if db != nil {
  rows, err := db.Query("SELECT id,jina,kiasi,aina,tarehe FROM miamala ORDER BY id DESC LIMIT 300")
  if err != nil { fmt.Fprintf(w, "<tr><td colspan=6>DB Error: %s</td></tr>", err) } else {
   for rows.Next() {
    var id int
    var k int64
    var j, a string
    var t time.Time
    rows.Scan(&id, &j, &k, &a, &t)
    fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td><span class=badge>%s</span></td><td><b>%d</b></td><td><a class=editbtn href=\"/edit?id=%d\">EDIT</a></td></tr>", id, t.Format("02/01/2006"), j, a, k, id)
   }
  }
 } else {
  fmt.Fprint(w, "<tr><td colspan=6>No DB</td></tr>")
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
 ks := r.FormValue("kiasi")
 if aina == "fedha_dharura" || aina == "matumizi" || aina == "mkopo" { jina = "Washusema Nzima" }
 kiasi, _ := strconv.ParseInt(ks, 10, 64)
 t, _ := time.Parse("2006-01-02", r.FormValue("tarehe"))
 if t.IsZero() { t = time.Now() }
 if db != nil {
  _, err := db.Exec("INSERT INTO miamala(jina,kiasi,aina,tarehe)VALUES($1,$2,$3,$4)", jina, kiasi, aina, t)
  if err != nil { log.Println("INSERT FAIL:", err) }
 }
 http.Redirect(w, r, "/?add=ok", 302)
}
func del(w http.ResponseWriter, r *http.Request) {
 if db != nil { db.Exec("DELETE FROM miamala WHERE id=$1", r.URL.Query().Get("id")) }
 http.Redirect(w, r, "/", 302)
}
func edit(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Query().Get("id")
 if id == "" { id = r.FormValue("id") }
 if r.Method == "POST" {
  kiasi, _ := strconv.ParseInt(r.FormValue("kiasi"), 10, 64)
  t, _ := time.Parse("2006-01-02", r.FormValue("tarehe"))
  if t.IsZero() { t = time.Now() }
  if db != nil {
   _, err := db.Exec("UPDATE miamala SET jina=$1,kiasi=$2,aina=$3,tarehe=$4 WHERE id=$5", r.FormValue("jina"), kiasi, r.FormValue("aina"), t, r.FormValue("id"))
   if err != nil { log.Println("UPDATE FAIL", err) }
  }
  http.Redirect(w, r, "/?edit=ok", 302)
  return
 }
 var jina, aina string
 var kiasi int64
 var t time.Time
 if db != nil { db.QueryRow("SELECT jina,kiasi,aina,tarehe FROM miamala WHERE id=$1", id).Scan(&jina, &kiasi, &aina, &t) }
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 fmt.Fprintf(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;padding:12px}.card{background:#fff;border-radius:16px;padding:18px}input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:8px 0}.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700;background:#0f7a3a;color:#fff}</style></head><body><div class=card><h2>✏️ EDIT Muamala #%s</h2><form method=POST action=/edit><input type=hidden name=id value=\"%s\"><label>Jina</label><input name=jina value=\"%s\" required><label>Kiasi</label><input name=kiasi type=number value=\"%d\" required><label>Tarehe</label><input type=date name=tarehe value=\"%s\" required><label>Aina</label><select name=aina required><option value=\"%s\" selected>%s (current)</option><option value=hisa_anzia>Hisa Anzia</option><option value=ada_mwezi>Ada ya Mwezi</option><option value=marejesho>Marejesho</option><option value=faini>Faini</option><option value=fedha_fomu>Fomu</option><option value=michango_mingine>Mingine</option><option value=fedha_dharura>Dharura</option><option value=matumizi>Matumizi</option><option value=mkopo>Mkopo</option></select><button class=btn>💾 HIFADHI MABADILIKO</button></form><br><a href=\"/\">⬅️ Rudi</a> | <a href=\"/delete?id=%s\" onclick=\"return confirm('Futa huu?')\" style=color:red>Futa huu muamala</a></div></body></html>", id, id, jina, kiasi, t.Format("2006-01-02"), aina, aina, id)
}
func exportCSV(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "text/csv")
 w.Header().Set("Content-Disposition", "attachment; filename=Uwashusema_"+time.Now().Format("2006")+".csv")
 writer := csv.NewWriter(w)
 writer.Write([]string{"ID", "Tarehe", "Jina", "Aina", "Kiasi"})
 if db != nil {
  rows, _ := db.Query("SELECT id,tarehe,jina,aina,kiasi FROM miamala ORDER BY tarehe DESC")
  if rows != nil {
   for rows.Next() {
    var id int
    var t time.Time
    var j, a string
    var k int64
    rows.Scan(&id, &t, &j, &a, &k)
    writer.Write([]string{fmt.Sprintf("%d", id), t.Format("02/01/2006"), j, a, fmt.Sprintf("%d", k)})
   }
  }
 }
 writer.Flush()
}
