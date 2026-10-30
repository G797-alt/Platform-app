func exportMikopo(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","text/csv")
 w.Header().Set("Content-Disposition","attachment; filename=Mikopo_Washusema.csv")
 writer:=csv.NewWriter(w)
 writer.Write([]string{"Jina","Mkopo Aliouchukua","Jumla Marejesho","Baki Anayodaiwa","Tarehe ya Mkopo","Status"})
 if db!=nil{
  rows,_:=db.Query(`SELECT m.jina, m.kiasi, COALESCE(SUM(rm.kiasi),0) as marejesho, m.baki, TO_CHAR(m.tarehe,'YYYY-MM-DD'), m.status FROM mikopo m LEFT JOIN marejesho_mikopo rm ON rm.jina=m.jina GROUP BY m.id ORDER BY m.id DESC`)
  if rows!=nil{for rows.Next(){var j,t,s string;var k,mm,b int64;rows.Scan(&j,&k,&mm,&b,&t,&s); writer.Write([]string{j,fmt.Sprintf("%d",k),fmt.Sprintf("%d",mm),fmt.Sprintf("%d",b),t,s})}}
 }
 writer.Flush()
}
