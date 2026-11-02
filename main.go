package main
import ("log";"net/http";"os")
func main(){p:=os.Getenv("PORT");if p==""{p="10000"};http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html; charset=utf-8");w.Write([]byte(html))});log.Println(p);http.ListenAndServe(":"+p,nil)}
const html=
"<!doctype html><meta charset=utf-8><meta name=viewport content='width=device-width,initial-scale=1'><title>U</title><script src='https://unpkg.com/xlsx@0.18.5/dist/xlsx.full.min.js'></script><header s"+
"tyle='background:#1a6d2f;color:#fff;padding:12px;text-align:center;font-weight:900'>Uwashusema</header><div style='padding:8px'><select id=m><option value=1>Jan<option value=2>Feb<option value="+
"3>Mar<option value=4>Apr<option value=5>Mei<option value=6>Jun<option value=7>Jul<option value=8>Ago<option value=9>Sep<option value=10>Okt<option value=11>Nov<option value=12>Des</select><input id=y "+
"type=number value=2026><button onclick=render()>Sawa</button><hr><div id=wl></div><button onclick=addW()>+ Mwanachama</button><hr><select id=sw></select><select id=sa><option value=hisa>Hisa<option va"+
"lue=ada>Ada<option value=marejesho>Marejesho<option value=faini_ada>Riba Ada 30%<option value=faini_marejesho>Riba Rejesho 30%<option value=michango>Michango<option value=mkopo_toka>Mkopo</select><inp"+
"ut id=k type=number placeholder=Kiasi><input id=ma placeholder=Maelezo><button onclick=addM()>Hifadhi</button><div id=ca style=color:green;font-weight:900></div><hr><button onclick=doExcel() style='ba"+
"ckground:#1565c0;color:#fff;padding:10px;width:100%'>Pakua Excel</button><div id=mu></div><div id=hi></div></div><script>D=JSON.parse(localStorage.getItem('v2')||'null')||{w:[],m:[]};function S(){loca"+
"lStorage.setItem('v2',JSON.stringify(D));render()}function addW(){j=prompt('Jina:');if(!j)return;D.w.push({j:j});S()}function addM(){w=sw.value;a=sa.value;k=+document.getElementById('k').value;ma=docu"+
"ment.getElementById('ma').value;if(!w||!k)return alert('Jaza');D.m.push({w:w,a:a,k:k,ma:ma,mm:+m.value,yy:+y.value,t:new Date().toLocaleString()});if(a=='mkopo_toka'){f=Math.round(k*0.05);D.m.push({w:"+
"w,a:'fomu',k:f,ma:'Fomu',mm:+m.value,yy:+y.value,t:'AUTO'})}document.getElementById('k').value='';S()}function render(){mm=+m.value;yy=+y.value;mi=D.m.filter(x=>x.mm==mm&&x.yy==yy);wl.innerHTML='';sw."+
"innerHTML='';if(D.w.length==0)wl.innerHTML='<small>Hakuna mwanachama</small>';D.w.forEach((x,i)=>{wl.innerHTML+=x.j+' <button onclick=\"D.w.splice('+i+',1);S()\">X</button><br>';sw.innerHTML+='<option>'"+
"+x.j+'</option>'});mu.innerHTML='';s={};mi.forEach(t=>{s[t.a]=(s[t.a]||0)+t.k});for(k in s)mu.innerHTML+=k+':'+s[k]+'<br>';hi.innerHTML='';mi.slice().reverse().forEach(t=>{hi.innerHTML+=t.w+'-'+t.a+'-'"+
"'+t.k+'<br>'})}kEl=document.getElementById('k');sa.addEventListener('change',()=>{k=+kEl.value;ca.innerHTML=sa.value=='mkopo_toka'&&k>0?'Fomu AUTO='+Math.round(k*0.05)+' (1.5M=75k)':''});kEl.addEventL"+
"istener('input',()=>sa.dispatchEvent(new Event('change')));function doExcel(){try{mm=+m.value;yy=+y.value;all=D.m.filter(x=>x.mm==mm&&x.yy==yy);if(!all.length)return alert('Hakuna data');ada=[['Jina',"+
"'Ada','Riba Ada']];mk=[['Jina','Mkopo','Fomu','Marejesho','Riba']];hs=[['Jina','Hisa']];am={};amR={};dm={};hm={};D.w.forEach(x=>{am[x.j]=0;amR[x.j]=0;dm[x.j]={c:0,f:0,r:0,fr:0};hm[x.j]=0});all.forEach"+
"(t=>{if(t.a=='ada')am[t.w]+=t.k;if(t.a=='faini_ada')amR[t.w]+=t.k;if(t.a=='mkopo_toka')dm[t.w].c+=t.k;if(t.a=='fomu')dm[t.w].f+=t.k;if(t.a=='marejesho')dm[t.w].r+=t.k;if(t.a=='faini_marejesho')dm[t.w]"+
".fr+=t.k;if(t.a=='hisa')hm[t.w]+=t.k});for(j in am)ada.push([j,am[j],amR[j]]);for(j in dm){d=dm[j];if(d.c>0)mk.push([j,d.c,d.f,d.r,d.fr])}for(j in hm)hs.push([j,hm[j]]);wb=XLSX.utils.book_new();XLSX.u"+
"tils.book_append_sheet(wb,XLSX.utils.aoa_to_sheet(ada),'1-Ada');XLSX.utils.book_append_sheet(wb,XLSX.utils.aoa_to_sheet(mk),'2-Mikopo 75k');XLSX.utils.book_append_sheet(wb,XLSX.utils.aoa_to_sheet(hs),"+
"'3-Hisa');XLSX.writeFile(wb,'Uwashusema.xlsx')}catch(e){alert(e.message)}}m.value=new Date().getMonth()+1;render();</script>"
