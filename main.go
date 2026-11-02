package main

import (
	"log"
	"net/http"
	"os"
)

const htmlPage = `<!DOCTYPE html><html lang="sw"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title><style>*{box-sizing:border-box}body{margin:0;font-family:system-ui;background:#f1f5f1;padding-bottom:20px}header{background:#1a6d2f;color:#fff;padding:16px;text-align:center;font-weight:900} .card{background:#fff;margin:10px;border-radius:12px;padding:12px;box-shadow:0 2px 8px #0001}h3{margin:0 0 10px 0;font-size:14px;color:#1a6d2f;border-bottom:1px solid #eee;padding-bottom:6px}.top3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:6px}.top3 div{padding:10px;border-radius:8px;text-align:center;color:#fff}.mapato{background:#2e7d32}.matumizi{background:#c62828}.salio{background:#1565c0}select,input{width:100%;padding:11px;border-radius:8px;border:1px solid #ccc;margin-bottom:8px}.btn{width:100%;padding:12px;border:0;border-radius:8px;font-weight:900;color:#fff;margin-top:6px}.green{background:#1a6d2f}.blue{background:#1565c0}.red{background:#c62828}.row{display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid #eee;font-size:13px}table{width:100%;border-collapse:collapse;font-size:12px}th{background:#e8f5e9;padding:6px;text-align:left}td{padding:7px;border-bottom:1px solid #eee}.badge{background:#ffecb3;padding:2px 6px;border-radius:8px;font-size:10px}</style></head><body>
<header>Uwashusema Mahida</header>
<div class="card"><div style="display:flex;gap:6px"><select id="mwezi"><option value="1">Jan</option><option value="2">Feb</option><option value="3">Mar</option><option value="4">Apr</option><option value="5">Mei</option><option value="6">Jun</option><option value="7">Jul</option><option value="8">Ago</option><option value="9">Sep</option><option value="10">Okt</option><option value="11">Nov</option><option value="12">Des</option></select><select id="mwaka"><option>2025</option><option selected>2026</option></select><button onclick="render()" style="background:#1a6d2f;color:#fff;border:0;border-radius:8px;padding:0 15px">Sawa</button></div><div class="top3" style="margin-top:10px"><div class="mapato"><small>MAPATO</small><b id="tMapato">0</b></div><div class="matumizi"><small>MATUMIZI</small><b id="tMatumizi">0</b></div><div class="salio"><small>SALIO</small><b id="tSalio">0</b></div></div><div style="margin-top:8px;display:flex;gap:6px"><small>Ada ya Mwezi inayotarajiwa: TZS </small><input id="adaSet" type="number" value="10000" style="width:100px;padding:4px" onchange="saveAda()"><small>(edit ukibadilisha)</small></div></div>
<div class="card"><h3>Wanachama <button onclick="addW()" style="float:right;background:#1a6d2f;color:#fff;border:0;padding:4px 10px;border-radius:6px">+ Mpya</button></h3><div id="wlist"></div></div>
<div class="card"><h3>Ingiza Data</h3><select id="selW"></select><select id="selAina"><optgroup label="MAPATO"><option value="hisa">Hisa</option><option value="ada">Ada ya Mwezi</option><option value="marejesho">Marejesho ya Mkopo</option><option value="faini">Faini / Riba 30%</option><option value="fomu">Fomu</option><option value="michango">Michango Mingine</option></optgroup><optgroup label="MATUMIZI"><option value="mkopo_toka">Mkopo kwa Mwanachama</option><option value="dharura">Dharura</option><option value="mkono_pole">Mkono wa Pole</option><option value="matumizi_mengine">Matumizi Mengine</option></optgroup></select><input id="kiasi" type="number" placeholder="Kiasi TZS"><input id="maelezo" placeholder="Maelezo"><button class="btn green" onclick="addM()">Hifadhi</button><small id="hint" style="color:#1a6d2f;display:none">Mkopo ukihifadhiwa, Fomu 5% itaongezeka yenyewe Mapato</small></div>
<div class="card"><h3>Ada ya Mwezi - Mwezi Huu (Riba 30% kwa wasioliwa)</h3><table><thead><tr><th>Jina</th><th>Ada Aliyolipa</th><th>Status</th><th>Riba 30%</th></tr></thead><tbody id="tblAda"></tbody></table></div>
<div class="card"><h3>Mikopo & Marejesho (Fomu 5% + Riba 30%)</h3><table><thead><tr><th>Jina</th><th>Mkopo</th><th>Marejesho</th><th>Baki</th><th>Riba</th></tr></thead><tbody id="tblDeni"></tbody></table></div>
<div class="card"><h3>Muhtasari</h3><div id="muhtasari"></div><button class="btn blue" onclick="backup()">📊 Pakua Excel (Sasa Inasoma Edits)</button><button class="btn green" onclick="wa()">📱 Tuma WhatsApp</button></div>
<div class="card"><h3>Historia</h3><div id="hist"></div></div>
<script>
let DATA=JSON.parse(localStorage.getItem('mahida_v4')||'null')||{wanachama:[{jina:"Charles Joseph Mgaya"}],miamala:[],adaTarajiwa:10000};
const curM=()=>+mwezi.value,curY=()=>+mwaka.value;
const save=()=>{localStorage.setItem('mahida_v4',JSON.stringify(DATA));render()};
function saveAda(){DATA.adaTarajiwa=+adaSet.value||10000;save();}
function addW(){let j=prompt('Jina:');if(!j)return;DATA.wanachama.push({jina:j});save()}
selAina.addEventListener('change',()=>{hint.style.display=selAina.value==='mkopo_toka'?'block':'none';});
function addM(){
 let w=selW.value,a=selAina.value,k=+kiasi.value,ma=maelezo.value;
 if(!w||!a||!k)return alert('Jaza');
 DATA.miamala.push({w,a,k,ma,mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()});
 if(a==='mkopo_toka'){
   let fomu=Math.round(k*0.05);
   DATA.miamala.push({w,a:'fomu',k:fomu,ma:'Fomu 5% ya mkopo '+k.toLocaleString(),mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()+' AUTO 5%'});
   alert('Mkopo '+k.toLocaleString()+' + Fomu 5% = '+fomu.toLocaleString()+' kwenye Mapato');
 }
 kiasi.value='';maelezo.value='';save();
}
function pigaRibaAda(jina){
 let ada=DATA.adaTarajiwa||10000;
 let riba=Math.round(ada*0.30);
 if(!confirm('Kumpiga '+jina+' Riba ya Ada 30% = '+riba.toLocaleString()+' ?'))return;
 DATA.miamala.push({w:jina,a:'faini',k:riba,ma:'Riba 30% ya kuchelewa Ada ya Mwezi TZS '+ada.toLocaleString(),mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()+' RIBA ADA'});
 save();
}
function pigaRibaMkopo(jina, baki){
 let rejesho=prompt('Mwanachama anatakiwa kulipa kiasi gani mwezi huu? (mfano 150000)', baki>0?Math.min(baki,150000):'');
 if(!rejesho)return; rejesho=+rejesho; if(!rejesho)return;
 let riba=Math.round(rejesho*0.30);
 if(!confirm(jina+' Riba 30% ya '+rejesho.toLocaleString()+' = '+riba.toLocaleString()+' ?'))return;
 DATA.miamala.push({w:jina,a:'faini',k:riba,ma:'Riba 30% ya kutolipa marejesho '+rejesho.toLocaleString(),mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()+' RIBA MKOPO'});
 save();
}
function render(){
 let m=curM(),y=curY(); adaSet.value=DATA.adaTarajiwa||10000;
 let mi=DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y);
 let mapA=['hisa','ada','marejesho','faini','fomu','michango'],matA=['mkopo_toka','dharura','mkono_pole','matumizi_mengine'];
 let map=0,mat=0; mi.forEach(t=>{if(mapA.includes(t.a))map+=t.k;if(matA.includes(t.a))mat+=t.k;});
 tMapato.textContent='TZS '+map.toLocaleString();tMatumizi.textContent='TZS '+mat.toLocaleString();tSalio.textContent='TZS '+(map-mat).toLocaleString();
 wlist.innerHTML='';selW.innerHTML='<option value="">Chagua Mwanachama</option>';
 DATA.wanachama.forEach((w,i)=>{wlist.innerHTML+='<div class="row"><b>'+w.jina+'</b><span><button onclick="let n=prompt(\'Jina jipya:\',DATA.wanachama['+i+'].jina);if(n){let old=DATA.wanachama['+i+'].jina;DATA.wanachama['+i+'].jina=n;DATA.miamala.forEach(x=>{if(x.w==old)x.w=n});save()}" style="border:0;background:#e3f2fd;border-radius:4px">✏️</button> <button onclick="if(confirm(\'Futa?\')){DATA.wanachama.splice('+i+',1);save()}" style="border:0;background:#ffebee;color:#c00;border-radius:4px">X</button></span></div>';selW.innerHTML+='<option>'+w.jina+'</option>'});
 // ADA TABLE
 let adaMap={}; DATA.wanachama.forEach(w=>adaMap[w.jina]=0); mi.filter(t=>t.a==='ada').forEach(t=>{adaMap[t.w]=(adaMap[t.w]||0)+t.k;});
 let tblAdaBody=document.getElementById('tblAda'); tblAdaBody.innerHTML='';
 for(let j in adaMap){
   let lipa=adaMap[j]; let tar=DATA.adaTarajiwa||10000; let status=lipa>=tar?'✅ Amelipa':'❌ Hajalipa';
   let color=lipa>=tar?'#2e7d32':'#c62828';
   tblAdaBody.innerHTML+='<tr><td>'+j+'</td><td style="font-weight:700;color:'+color+'">'+lipa.toLocaleString()+'</td><td>'+status+'</td><td>'+(lipa<tar?'<button onclick="pigaRibaAda(\''+j+'\')" class="red" style="border:0;padding:5px 8px;border-radius:5px;color:#fff;font-size:11px">Riba 30% = '+Math.round(tar*0.30).toLocaleString()+'</button>':'')+'</td></tr>';
 }
 // DENI MKOPO
 let dmap={}; DATA.miamala.forEach(t=>{if(!dmap[t.w])dmap[t.w]={c:0,r:0}; if(t.a=='mkopo_toka')dmap[t.w].c+=t.k; if(t.a=='marejesho')dmap[t.w].r+=t.k;});
 let td=document.getElementById('tblDeni'); td.innerHTML='';
 for(let j in dmap){let b=dmap[j].c-dmap[j].r; if(dmap[j].c>0||dmap[j].r>0){td.innerHTML+='<tr><td>'+j+'</td><td>'+dmap[j].c.toLocaleString()+'</td><td>'+dmap[j].r.toLocaleString()+'</td><td style="color:#c00;font-weight:800">'+b.toLocaleString()+'</td><td>'+(b>0?'<button onclick="pigaRibaMkopo(\''+j+'\','+b+')" style="border:0;background:#ef6c00;color:#fff;padding:5px 8px;border-radius:5px;font-size:11px">Riba 30%</button>':'')+'</td></tr>';}}
 if(!td.innerHTML)td.innerHTML='<tr><td colspan=5 style="color:#999">Hakuna mikopo</td></tr>';
 let mu=document.getElementById('muhtasari'); mu.innerHTML=''; let labels={hisa:'Hisa',ada:'Ada ya Mwezi',marejesho:'Marejesho',faini:'Faini / Riba 30% (Ada+Mkopo)',fomu:'Fomu 5%',michango:'Michango',mkopo_toka:'Mkopo Toka',dharura:'Dharura',mkono_pole:'Mkono wa Pole',matumizi_mengine:'Matumizi Mengine'}; let sum={}; document.querySelectorAll('#selAina option').forEach(o=>sum[o.value]=0); mi.forEach(t=>sum[t.a]+=t.k); for(let k in labels){if(sum[k]>0)mu.innerHTML+='<div class="row"><span>'+labels[k]+'</span><b>'+sum[k].toLocaleString()+'</b></div>'}
 let hi=document.getElementById('hist'); hi.innerHTML=''; [...mi].reverse().forEach(t=>{let idx=DATA.miamala.indexOf(t); hi.innerHTML+='<div class="row"><div><b>'+t.w+'</b> - '+t.a+' '+(t.a=='fomu'?'<span class=badge>5% AUTO</span>':'')+(t.a=='faini'?'<span class=badge style=background:#ffcdd2>30% RIBA</span>':'')+'<br><small>'+t.t+' - '+t.k.toLocaleString()+' '+(t.ma||'')+'</small></div><div><button onclick="let v=prompt(\'Kiasi jipya:\',DATA.miamala['+idx+'].k);if(v!==null){DATA.miamala['+idx+'].k=+v;save()}" style="border:0;background:#e3f2fd;border-radius:4px">✏️</button> <button onclick="if(confirm(\'Futa?\')){DATA.miamala.splice('+idx+',1);save()}" style="border:0;background:#ffebee;color:#c00;border-radius:4px">X</button></div></div>'});
}
function backup(){
 let m=curM(),y=curY();
 // FIX: Inasoma DATA halisi, si cache ya zamani
 let all=DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y);
 let dmap={};
 all.forEach(t=>{
   if(!dmap[t.w])dmap[t.w]={c:0,r:0,fomu:0,faini:0,ada:0};
   if(t.a=='mkopo_toka')dmap[t.w].c+=t.k;
   if(t.a=='marejesho')dmap[t.w].r+=t.k;
   if(t.a=='fomu')dmap[t.w].fomu+=t.k;
   if(t.a=='faini')dmap[t.w].faini+=t.k;
   if(t.a=='ada')dmap[t.w].ada+=t.k;
 });
 let csv='Jina,Mkopo Aliouchukua,Jumla Marejesho,Ada ya Mwezi,Fomu 5%,Riba 30% (Ada+Mkopo),Baki la Mkopo\n';
 for(let j in dmap){let b=dmap[j].c-dmap[j].r; csv+='"'+j+'",'+dmap[j].c+','+dmap[j].r+','+dmap[j].ada+','+dmap[j].fomu+','+dmap[j].faini+','+b+'\n';}
 if(Object.keys(dmap).length==0)csv+='"Hakuna data",0,0,0,0,0,0\n';
 let a=document.createElement('a');a.href=URL.createObjectURL(new Blob([csv],{type:'text/csv;charset=utf-8;'}));a.download='Uwashusema_Fomu5_Riba30_ExcelFix_'+m+'_'+y+'.csv';a.click();
 alert('Excel imepakuliwa na edits zako zote za mwisho!');
}
function wa(){let txt='*Uwashusema Mahida* '+curM()+'/'+curY()+'\nMAPATO: '+tMapato.textContent+'\nMATUMIZI: '+tMatumizi.textContent+'\nSALIO: '+tSalio.textContent+'\nAda tarajiwa: '+(DATA.adaTarajiwa||10000)+'\n';window.open('https://wa.me/?text='+encodeURIComponent(txt),'_blank')}
mwezi.value=new Date().getMonth()+1;render();
</script></body></html>`

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	})
	log.Println("Running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
