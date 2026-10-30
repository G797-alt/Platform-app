
package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func init() {
	u := os.Getenv("DATABASE_URL")
	if u == "" {
		return
	}
	db, _ = sql.Open("pgx", u)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.PingContext(ctx)
	db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)
	db.Exec(`CREATE TABLE IF NOT EXISTS mikopo(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,baki BIGINT,tarehe TIMESTAMP DEFAULT NOW(),status TEXT DEFAULT 'ipo')`)
	db.Exec(`CREATE TABLE IF NOT EXISTS marejesho_mikopo(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,tarehe TIMESTAMP DEFAULT NOW())`)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ongeza", ongeza)
	http.HandleFunc("/mwanachama", mwanachama)
	http.HandleFunc("/export", exportCSV)
	http.HandleFunc("/export-mikopo", exportMikopo)
	http.HandleFunc("/delete", del)
	p := os.Getenv("PORT")
	if p == "" {
		p = "10000"
	}
	http.ListenAndServe(":"+p, nil)
}

func getTotals() (mapato, dharura, matumizi, mikopo int64) {
	if db == nil {
		return
	}
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')`).Scan(&mapato)
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'`).Scan(&dharura)
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'`).Scan(&matumizi)
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'`).Scan(&mikopo)
	return
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	mapato, dharura, matumizi, mikopo := getTotals()
	bakaa := (mapato + dharura) - (matumizi + mikopo)
	waText := fmt.Sprintf("UWASHUSEMA MAHIDA RIPOTI\nBAKAA: TZS %d", bakaa)
	waLink := "https://wa.me/?text=" + url.QueryEscape(waText)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<!DOCTYPE html><html><head><meta name=viewport content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title>
<style>
body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:10px 10px 80px}
.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px;box-shadow:0 4px 12px #0000000d}
.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}
.big{font-size:28px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.small{font-size:11px;opacity:.8}.big2{font-size:18px;font-weight:800}
.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700;margin:6px 0}
.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}
input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0;background:#fafafa}
table{width:100%%;border-collapse:collapse;font-size:12px}th{text-align:left;padding:8px;background:#f6f6f6}td{padding:8px;border-bottom:1px solid #f0f0f0}
.badge{padding:3px 7px;border-radius:20px;background:#eef2ee;font-size:10px}
</style></head><body>
<div class="card bakaa"><div class=small>BAKAA</div><div class=big>TZS %d</div><div class=small>Mapato %d + Dharura %d - Matumizi %d - Mikopo %d</div></div>
<div class=grid2>
<div class=card><div class=small>MAPATO</div><div class=big2>%d</div></div>
<div class=card><div class=small>DHARURA</div><div class=big2>%d</div></div>
<div class=card><div class=small>MATUMIZI</div><div class=big2 style="color:red">%d</div></div>
<div class=card><div class=small>MIKOPO</div><div class=big2 style="color:orange">%d</div></div>
</div>
<div class=card>
<a href=/export><button class="btn btn-white">📥 Pakua Excel ya Miamala Yote (Inasoma Sasa)</button></a>
<a href=/export-mikopo><button class="btn btn-green">📊 Pakua Excel ya Mikopo Inayopungua</button></a>
<a href="%s"><button class="btn btn-wa">WhatsApp</button></a>
</div>
<div class=card><h3>+ Ongeza Muamala</h3>
<form method=POST action=/ongeza>
<input name=jina placeholder="Jina la Mwanachama" required>
<input name=kiasi type=number placeholder="Kiasi TZS" required>
<select name=aina required>
<option value="" disabled selected>Chagua Aina</option>
<option value=hisa_anzia>Hisa Anzia</option>
<option value=ada_mwezi>Ada ya Mwezi</option>
<option value=marejesho>Marejesho ya Mkopo - Itapunguza Deni</option>
<option value=faini>Faini</option>
<option value=fedha_fomu>Fedha ya Fomu</option>
<option value=michango_mingine>Michango Mingine</option>
<option value=fedha_dharura>Fedha ya Dharura (Kikundi)</option>
<option value=matumizi>Matumizi (Kikundi)</option>
<option value=mkopo>Mkopo Aliyochukua (Mwanachama)</option>
</select>
<button class="btn btn-green">Ongeza</button></form></div>
`, bakaa, mapato, dharura, matumizi, mikopo, mapato, dharura, matumizi, mikopo, waLink)

	if db != nil {
		fmt.Fprint(w, `<div class=card><h3>Mikopo - Baki Inayopungua Live</h3><table><tr><th>Jina</th><th>Kiasi</th><th>Baki</th><th>Status</th></tr>`)
		rows, _ := db.Query(`SELECT jina,kiasi,baki,status FROM mikopo ORDER BY id DESC LIMIT 100`)
		if rows != nil {
			for rows.Next() {
				var j, s string
				var k, b int64
				rows.Scan(&j, &k, &b, &s)
				c := "green"
				if b > 0 {
					c = "red"
				}
				fmt.Fprintf(w, `<tr><td>%s</td><td>%d</td><td style="color:%s;font-weight:800">%d</td><td>%s</td></tr>`, j, k, c, b, s)
			}
			rows.Close()
		}
		fmt.Fprint(w, `</table></div>`)

		fmt.Fprint(w, `<div class=card><h3>Miamala 20 ya Mwisho</h3><table><tr><th>Jina</th><th>Aina</th><th>Kiasi</th></tr>`)
		rows2, _ := db.Query(`SELECT jina,aina,kiasi FROM miamala ORDER BY id DESC LIMIT 20`)
		if rows2 != nil {
			for rows2.Next() {
				var j, a string
				var k int64
				rows2.Scan(&j, &a, &k)
				fmt.Fprintf(w, `<tr><td>%s</td><td>%s</td><td>%d</td></tr>`, j, a, k)
			}
			rows2.Close()
		}
		fmt.Fprint(w, `</table></div>`)
	}
	fmt.Fprint(w, `</body></html>`)
}

func ongeza(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	jinaOrig := r.FormValue("jina")
	jina := jinaOrig
	aina := r.FormValue("aina")
	kiasiStr := r.FormValue("kiasi")
	var kiasi int64
	fmt.Sscan(kiasiStr, &kiasi)
	if aina == "fedha_dharura" || aina == "matumizi" {
		jina = "Washusema Nzima"
	}
	if db != nil {
		db.Exec("INSERT INTO miamala(jina,kiasi,aina)VALUES($1,$2,$3)", jina, kiasi, aina)
		if aina == "mkopo" {
			db.Exec("INSERT INTO mikopo(jina,kiasi,baki)VALUES($1,$2,$2)", jinaOrig, kiasi)
		}
		if aina == "marejesho" {
			db.Exec("UPDATE mikopo SET baki = baki - $1 WHERE id = (SELECT id FROM mikopo WHERE jina=$2 AND baki>0 ORDER BY id DESC LIMIT 1)", kiasi, jinaOrig)
			db.Exec("UPDATE mikopo SET status='imelipwa' WHERE baki<=0")
			db.Exec("INSERT INTO marejesho_mikopo(jina,kiasi)VALUES($1,$2)", jinaOrig, kiasi)
		}
	}
	http.Redirect(w, r, "/?added=1", 302)
}

func mwanachama(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if db != nil {
		db.Exec("INSERT INTO wanachama(jina,simu,namba)VALUES($1,$2,$3)", r.FormValue("jina"), r.FormValue("simu"), r.FormValue("namba"))
	}
	http.Redirect(w, r, "/", 302)
}
func del(w http.ResponseWriter, r *http.Request) {
	if db != nil {
		db.Exec("DELETE FROM miamala WHERE id=$1", r.URL.Query().Get("id"))
	}
	http.Redirect(w, r, "/", 302)
}

// HAPA NDIPO NILIPOFIX - EXCEL SASA INASOMA
func exportCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=Miamala_Yote.csv")
	writer := csv.NewWriter(w)
	writer.Write([]string{"Tarehe", "Jina", "Aina", "Kiasi"})
	if db != nil {
		rows, err := db.Query(`SELECT tarehe, jina, aina, kiasi FROM miamala ORDER BY id DESC`)
		if err == nil && rows != nil {
			for rows.Next() {
				var t time.Time
				var j, a string
				var k int64
				rows.Scan(&t, &j, &a, &k)
				writer.Write([]string{t.Format("2006-01-02 15:04:05"), j, a, fmt.Sprintf("%d", k)})
			}
			rows.Close()
		}
	}
	writer.Flush()
}

func exportMikopo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=Mikopo_Inayopungua.csv")
	writer := csv.NewWriter(w)
	writer.Write([]string{"Jina", "Mkopo Aliouchukua", "Jumla Marejesho", "Baki Anayodaiwa", "Tarehe", "Status"})
	if db != nil {
		rows, err := db.Query(`SELECT jina, kiasi, baki, tarehe, status FROM mikopo ORDER BY id DESC`)
		if err == nil && rows != nil {
			for rows.Next() {
				var j, s string
				var k, b int64
				var t time.Time
				rows.Scan(&j, &k, &b, &t, &s)
				marejesho := k - b
				if marejesho < 0 {
					marejesho = 0
				}
				writer.Write([]string{j, fmt.Sprintf("%d", k), fmt.Sprintf("%d", marejesho), fmt.Sprintf("%d", b), t.Format("2006-01-02"), s})
			}
			rows.Close()
		}
	}
	writer.Flush()
}
