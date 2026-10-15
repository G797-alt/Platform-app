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
  db, _ = sql.Open("pgx", u)
  db.Exec("CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)")
  db.Exec("CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT)")
  db.Exec("ALTER TABLE miamala ADD COLUMN IF NOT EXISTS tarehe TIMESTAMP DEFAULT NOW()")
 }
 http.HandleFunc("/", home)
 http.HandleFunc("/ongeza", ongeza)
 http.HandleFunc("/mwanachama", mwanachama)
 http.HandleFunc("/export", exportCSV)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/delete", del)
 p := os.Getenv("PORT")
 if p == "" { p = "10000" }
 log.Println("V12 MIKOPO start", p)
 http.ListenAndServe("0.0.0.0:"+p, nil)
}

func getTotals() (mapato, dharura, matumizi, mikopo, marejesho int64, count int) {
 if db == nil { return }
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','faini','fedha_fomu','michango_mingine')").Scan(&mapato)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'").Scan(&dharura)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'").Scan(&matumizi)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'").Scan(&mikopo)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='marejesho'").Scan(&marejesho)
 db.QueryRow("SELECT COUNT(*) FROM miamala").Scan(&count)
 return
}

func home(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-cache")
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 mapato, dharura, matumizi, mikopo, marejesho, count := getTotals()
 bakaa := (mapato + dharura + marejesho) - (matumizi + mikopo)
 mwaka := time.Now().Year()
 waLink := "https://wa.me/?text=" + url.QueryEscape(fmt.Sprintf("UWASHUSEMA MAHIDA %d%0ABAKAA TZS %d%0AMikopo %d - Marejesho %d", mwaka, bakaa, mikopo, marejesho))
 msg := ""
 if r.URL.Query().Get("edit") == "ok" { msg = "<div class=card style=background:#e8f5e9;border-left:4px solid green>✅ Imehaririwa</div>" }
 if r.URL.Query().Get("add") == "ok" { msg = "<div class=card style=background:#e8f5e9;border-left:4px solid green>✅ Imehifadhiwa</div>" }
 today := time.Now().Format("2006-01-02")

 fmt.Fprint(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:10px;padding-bottom:90px}.header{padding:16px 6px;font-weight:800;font-size:20px}.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px;box-shadow:0 2px 6px #0001}.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}.big{font-size:28px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}.small{font-size:11px;opacity:.9}.big2{font-size:18px;font-weight:800}.btn{width:100%;padding:12px;border-radius:12px;border:0;font-weight:700;margin:6px 0}.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}input,select{width:100%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0}table{width:100%;border-collapse:collapse;font-size:12px}th{padding:8px;background:#222;color:#fff;text-align:left}td{padding:10px;border-bottom:1px solid #eee}.badge{background:#eef2ee;padding:3px 7px;border-radius:20px;font-size:10px}.editbtn{background:#1976d2;color:#fff;padding:6px 12px;border-radius:8px;text-decoration:none;font-weight:700;font-size:11px}.kopaji{background:#fff7ed;border-left:4px solid #ff9800;padding:12px;margin:8px 0;border-radius:10px}.progress{height:8px;background:#eee;border-radius:10px;overflow:hidden;margin:6px 0}.progress-bar{height:100%;background:#0f7a3a}.year{background:#0002;padding:3px 8px;border-radius:10px;font-size:12px}</style></head><body>")
 fmt.Fprintf(w, "<div class=header>Uwashusema Mahida <span class=year>%d | %d</span></div>%s", mwaka, count, msg)
 fmt.Fprintf(w, "<div class=\"card bakaa\"><div class=small>BAKAA %d (Mapato+Dharura+Marejesho - Matumizi-Mikopo)</div><div class=big>TZS %d</div><div class=small>Mapato %d + Dharura %d + Marejesho %d - Matumizi %d - Mikopo %d</div></div>", mwaka, bakaa, mapato, dharura, marejesho, matumizi, mikopo)
 fmt.Fprintf(w, "<div class=grid2><div class=card><div class=small>MAPATO</div><div class=big2>%d</div></div><div class=card><div class=small>DHARURA</div><div class=big2>%d</div></div><div class=card><div class=small>MATUMIZI</div><div class=big2 style=color:red>%d</div></div><div class=card><div class=small>MIKOPO</div><div class=big2 style=color:orange>%d</div></div></div>", mapato, dharura, matumizi, mikopo)

 // ====== SEHEMU MPYA YA WAKOPAJI ======
 fmt.Fprint(w, "<div class=card><h3>👥 Waliochukua Mkopo - Mwenendo wa Marejesho</h3>")
 if db != nil {
  rows, _ := db.Query(`
   SELECT jina,
   COALESCE(SUM(CASE WHEN aina='mkopo' THEN kiasi ELSE 0 END),0) as jumla_kop,
   COALESCE(SUM(CASE WHEN aina='marejesho' THEN kiasi ELSE 0 END),0) as jumla_rej
   FROM miamala WHERE jina!='Washusema Nzima' AND jina!='' GROUP BY jina HAVING SUM(CASE WHEN aina='mkopo' THEN kiasi ELSE 0 END) > 0 ORDER BY jumla_kop DESC`)
  if rows != nil {
   found := false
   for rows.Next() {
    found = true
    var jina string
    var kop, rej int64
    rows.Scan(&jina, &kop, &rej)
    baki := kop - rej
    perc := 0
    if kop > 0 { perc = int(rej * 100 / kop) }
    color := "#ff9800"
    status := fmt.Sprintf("Baki TZS %d", baki)
    if baki <= 0 { color = "#0f7a3a"; status = "✅ AMEMALIZA" } else if perc == 0 { color = "#d32f2f"; status = fmt.Sprintf("Baki %d - HAJAANZA", baki) }
    fmt.Fprintf(w, "<div class=kopaji style=border-left-color:%s><div style=display:flex;justify-content:space-between><b>%s</b><span style=font-size:11px>%s</span></div><div style=font-size:12px>Alikopa: <b>%d</b> | Rejesha: <b>%d</b> | %s</div><div class=progress><div class=progress-bar style=width:%d%%;background:%s></div></div><div style=font-size:11px>%d%% Amerejesha</div></div>", color, jina, status, kop, rej, status, perc, color, perc)
   }
   if !found { fmt.Fprint(w, "<div class=small>Hakuna mikopo bado. Ongeza mkopo uone hapa.</div>") }
  }
 }
 fmt.Fprint(w, "</div>")

 fmt.Fprintf(w, "<div class=card><a href=/export><button class=\"btn btn-white\">📥 Pakua Excel ya Mikopo %d</button></a><a href=\"%s\" target=_blank><button class=\"btn btn-wa\">💬 Tuma Ripoti WhatsApp</button></a></div>", mwaka, waLink)
 fmt.Fprintf(w, "<div class=card><h3>+ Ongeza Muamala (Mkopo / Marejesho)</h3><form method=POST action=/ongeza><input id=jina name=jina placeholder=\"Jina la Mkopaji / Mwanachama\" required><input name=kiasi type=number placeholder=\"Kiasi TZS\" required><input type=date name=tarehe value=\"%s\" required><select id=aina name=aina required><option value=\"\" disabled selected>Chagua Aina</option><option value=hisa_anzia>Hisa Anzia</option><option value=ada_mwezi>Ada ya Mwezi</option><option value=marejesho>MAREJESHO YA MKOPO</option><option value=faini>Faini</option><option value=fedha_fomu>Fedha ya Fomu</option><option value=michango_mingine>Michango Mingine</option><option value=fedha_dharura>FEDHA YA DHARURA</option><option value=matumizi>MATUMIZI</option><option value=mkopo>MKOPO - Mtu Amechukuwa</option></select><button class=\"btn btn-green\">HIFADHI</button></form><div class=small style=margin-top:8px>💡 Mkopo: weka jina la aliyekopa. Marejesho: weka jina lile lile la mkopaji anayerejesha.</div></div>", today)
 fmt.Fprint(w, "<div class=card><h3>👤 Mwanachama Mpya</h3><form method=POST action=/mwanachama><input name=jina placeholder=\"Jina Kamili\" required><input name=simu placeholder=\"Simu\"><input name=namba placeholder=\"Namba\"><button class=\"btn btn-green\">Sajili</button></form></div>")

 fmt.Fprintf(w, "<div class=card><h3>📜 Historia (%d) - EDIT</h3><table><tr><th>Tarehe</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>EDIT</th></tr>", count)
 if db != nil {
  rows, _ := db.Query("SELECT id,jina,kiasi,aina,tarehe FROM miamala ORDER BY id DESC LIMIT 300")
  if rows != nil {
   for rows.Next() {
    var id int
    var k int64
    var j, a string
    var t time.Time
    rows.Scan(&id, &j, &k, &a, &t)
    fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td><span class=badge>%s</span></td><td><b>%d</b></td><td><a class=editbtn href=\"/edit?id=%d\">EDIT</a></td></tr>", t.Format("02/01/06"), j, a, k, id)
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
 kiasi, _ := strconv.ParseInt(r.FormValue("kiasi"), 10, 64)
 t, _ := time.Parse("2006-01-02", r.FormValue("tarehe"))
 if t.IsZero() { t = time.Now() }
 // Hapa ndipo fix - Mkopo lazima awe na jina, sio Washusema Nzima
 if (aina == "fedha_dharura" || aina == "matumizi") && jina == "" { jina = "Washusema Nzima" }
 if db != nil { db.Exec("INSERT INTO miamala(jina,kiasi,aina,tarehe)VALUES($1,$2,$3,$4)", jina, kiasi, aina, t) }
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
  if db != nil { db.Exec("UPDATE miamala SET jina=$1,kiasi=$2,aina=$3,tarehe=$4 WHERE id=$5", r.FormValue("jina"), kiasi, r.FormValue("aina"), t, r.FormValue("id")) }
  http.Redirect(w, r, "/?edit=ok", 302)
  return
 }
 var jina, aina string
 var kiasi int64
 var t time.Time
 if db != nil { db.QueryRow("SELECT jina,kiasi,aina,tarehe FROM miamala WHERE id=$1", id).Scan(&jina, &kiasi, &aina, &t) }
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 fmt.Fprintf(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;padding:12px}.card{background:#fff;border-radius:16px;padding:18px}input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:8px 0}.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700;background:#0f7a3a;color:#fff}</style></head><body><div class=card><h2>✏️ EDIT #%s</h2><form method=POST action=/edit><input type=hidden name=id value=\"%s\"><label>Jina la Mkopaji</label><input name=jina value=\"%s\" required><label>Kiasi</label><input name=kiasi type=number value=\"%d\" required><label>Tarehe</label><input type=date name=tarehe value=\"%s\" required><label>Aina</label><select name=aina><option value=\"%s\" selected>%s</option><option value=mkopo>MKOPO</option><option value=marejesho>MAREJESHO</option><option value=hisa_anzia>Hisa</option><option value=ada_mwezi>Ada</option><option value=faini>Faini</option><option value=matumizi>Matumizi</option><option value=fedha_dharura>Dharura</option></select><button class=btn>💾 HIFADHI</button></form><br><a href=\"/\">Rudi</a> | <a href=\"/delete?id=%s\" style=color:red>Futa</a></div></body></html>", id, id, jina, kiasi, t.Format("2006-01-02"), aina, aina, id)
}
func exportCSV(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "text/csv")
 w.Header().Set("Content-Disposition", "attachment; filename=Mikopo_Mahida_"+time.Now().Format("2006")+".csv")
 writer := csv.NewWriter(w)
 writer.Write([]string{"Tarehe", "Jina la Mkopaji", "Aina", "Kiasi", "Mwaka"})
 if db != nil {
  rows, _ := db.Query("SELECT tarehe,jina,aina,kiasi FROM miamala ORDER BY jina, tarehe DESC")
  if rows != nil {
   for rows.Next() {
    var t time.Time
    var j, a string
    var k int64
    rows.Scan(&t, &j, &a, &k)
    writer.Write([]string{t.Format("02/01/2006"), j, a, fmt.Sprintf("%d", k), t.Format("2006")})
   }
  }
  writer.Write([]string{""})
  writer.Write([]string{"MUHTASARI WA WAKOPAJI"})
  rows2, _ := db.Query("SELECT jina, SUM(CASE WHEN aina='mkopo' THEN kiasi ELSE 0 END), SUM(CASE WHEN aina='marejesho' THEN kiasi ELSE 0 END) FROM miamala GROUP BY jina HAVING SUM(CASE WHEN aina='mkopo' THEN kiasi ELSE 0 END)>0")
  if rows2 != nil {
   writer.Write([]string{"Jina", "Jumla Alikopa", "Jumla Amerejesha", "Baki"})
   for rows2.Next() {
    var j string
    var kop, rej int64
    rows2.Scan(&j, &kop, &rej)
    writer.Write([]string{j, fmt.Sprintf("%d", kop), fmt.Sprintf("%d", rej), fmt.Sprintf("%d", kop-rej)})
   }
  }
 }
 writer.Flush()
}
