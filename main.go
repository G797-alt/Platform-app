rm ~/mbaochap/index.html
cat > ~/mbaochap/index.html <<'HTML'
<!DOCTYPE html><html lang="sw"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title><style>*{box-sizing:border-box}body{margin:0;font-family:system-ui;background:#f1f5f1;padding-bottom:20px}header{background:#1a6d2f;color:#fff;padding:16px;text-align:center;font-weight:900;font-size:20px}.card{background:#fff;margin:10px;border-radius:12px;padding:12px;box-shadow:0 2px 8px #0001}h3{margin:0 0 10px 0;font-size:14px;color:#1a6d2f;border-bottom:1px solid #eee;padding-bottom:6px}.top3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:6px}.top3 div{padding:10px;border-radius:8px;text-align:center;color:#fff}.mapato{background:#2e7d32}.matumizi{background:#c62828}.salio{background:#1565c0}select,input{width:100%;padding:11px;border-radius:8px;border:1px solid #ccc;margin-bottom:8px}.btn{width:100%;padding:12px;border:0;border-radius:8px;font-weight:900;color:#fff;margin-top:6px}.green{background:#1a6d2f}.blue{background:#1565c0}table{width:100%;border-collapse:collapse;font-size:13px}th{background:#e8f5e9;padding:6px;text-align:left}td{padding:7px;border-bottom:1px solid #eee}[contenteditable]{background:#fffde7;outline:1px solid #ffca28;border-radius:4px}.row{display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid #eee;font-size:13px}</style></head><body>
<header>Uwashusema Mahida</header>
<div class="card"><div style="display:flex;gap:6px"><select id="mwezi"><option value="1">Jan</option><option value="2">Feb</option><option value="3">Mar</option><option value="4">Apr</option><option value="5">Mei</option><option value="6">Jun</option><option value="7">Jul</option><option value="8">Ago</option><option value="9">Sep</option><option value="10">Okt</option><option value="11">Nov</option><option value="12">Des</option></select><select id="mwaka"><option>2025</option><option selected>2026</option></select><button onclick="render()" style="background:#1a6d2f;color:#fff;border:0;border-radius:8px;padding:0 15px">Sawa</button></div><div class="top3" style="margin-top:10px"><div class="mapato"><small>MAPATO</small><b id="tMapato">0</b></div><div class="matumizi"><small>MATUMIZI</small><b id="tMatumizi">0</b></div><div class="salio"><small>SALIO</small><b id="tSalio">0</b></div></div></div>
<div class="card"><h3>Wanachama <button onclick="addW()" style="float:right;background:#1a6d2f;color:#fff;border:0;padding:4px 10px;border-radius:6px">+ Mpya</button></h3><div id="wlist"></div></div>
<div class="card"><h3>Ingiza Data</h3><select id="selW"></select><select id="selAina"><optgroup label="MAPATO"><option value="hisa">Hisa</option><option value="ada">Ada ya Mwezi</option><option value="marejesho">Marejesho ya Mkopo</option><option value="faini">Faini</option><option value="fomu">Fomu</option><option value="michango">Michango Mingine</option></optgroup><optgroup label="MATUMIZI"><option value="mkopo_toka">Mkopo kwa Mwanachama</option><option value="dharura">Dharura</option><option value="mkono_pole">Mkono wa Pole</option><option value="matumizi_mengine">Matumizi Mengine</option></optgroup></select><input id="kiasi" type="number" placeholder="Kiasi TZS"><input id="maelezo" placeholder="Maelezo"><button class="btn green" onclick="addM()">Hifadhi</button></div>
<div class="card"><h3>Wanaorudisha Mikopo - EDITABLE</h3><table><thead><tr><th>Jina</th><th>Mkopo Aliouchukua</th><th>Marejesho</th></tr></thead><tbody id="tblRudisha"></tbody></table><small style="color:#888">Gusa namba kuhariri moja kwa moja</small></div>
<div class="card"><h3>Wenye Deni la Mkopo - EDITABLE</h3><table><thead><tr><th>Jina</th><th>Chukua</th><th>Rudisha</th><th>Baki</th></tr></thead><tbody id="tblDeni"></tbody></table></div>
<div class="card"><h3>Muhtasari + Tuma</h3><div id="muhtasari"></div><button class="btn blue" onclick="backup()">📊 Pakua Excel (Inasoma Edits)</button><button class="btn green" onclick="wa()">📱 Tuma WhatsApp Ripoti</button></div>
<div class="card"><h3>Historia</h3><div id="hist"></div></div>
<script>
let DATA=JSON.parse(localStorage.getItem('mahida_vFinal')||'null')||{wanachama:[{jina:"Charles Joseph Mgaya"}],miamala:[]};
const curM=()=>+mwezi.value,curY=()=>+mwaka.value;
const save=()=>{localStorage.setItem('mahida_vFinal',JSON.stringify(DATA));render()};
function addW(){let j=prompt('Jina:');if(!j)return;DATA.wanachama.push({jina:j});save()}
function addM(){let w=selW.value,a=selAina.value,k=+kiasi.value,ma=maelezo.value;if(!w||!a||!k)return alert('Jaza');DATA.miamala.push({w,a,k,ma,mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()});kiasi.value='';save()}
// Hii ndio fix ya edit inayoandika moja kwa moja kwenye DATA
function editDirect(type, idx, field, el){
  let val = el.innerText.replace(/,/g,'').trim();
  let num = +val;
  if(isNaN(num)) return;
  if(type=='miamala'){ DATA.miamala[idx][field]=num; }
  if(type=='agg'){ // inahariri jumla
    // tunatafuta miamala ya huyo mtu na ku-adjust
  }
  localStorage.setItem('mahida_vFinal',JSON.stringify(DATA));
  renderCalcOnly();
}
function renderCalcOnly(){
  let m=curM(),y=curY(),mi=DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y);
  let mapA=['hisa','ada','marejesho','faini','fomu','michango'],matA=['mkopo_toka','dharura','mkono_pole','matumizi_mengine'];
  let map=0,mat=0; mi.forEach(t=>{if(mapA.includes(t.a))map+=t.k;if(matA.includes(t.a))mat+=t.k;});
  tMapato.textContent='TZS '+map.toLocaleString();tMatumizi.textContent='TZS '+mat.toLocaleString();tSalio.textContent='TZS '+(map-mat).toLocaleString();
}
function render(){
let m=curM(),y=curY(),mi=DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y);
renderCalcOnly();
wlist.innerHTML='';selW.innerHTML='<option value="">Chagua Mwanachama</option>';
DATA.wanachama.forEach((w,i)=>{wlist.innerHTML+=`<div class="row"><b>${w.jina}</b><span><button onclick="let n=prompt('Jina jipya:',DATA.wanachama[${i}].jina);if(n){let old=DATA.wanachama[${i}].jina;DATA.wanachama[${i}].jina=n;DATA.miamala.forEach(x=>{if(x.w==old)x.w=n});save()}" style="border:0;background:#e3f2fd;border-radius:4px">✏️</button> <button onclick="if(confirm('Futa?')){DATA.wanachama.splice(${i},1);save()}" style="border:0;background:#ffebee;color:#c00;border-radius:4px">X</button></span></div>`;selW.innerHTML+=`<option>${w.jina}</option>`});

// RUDISHA - editable moja kwa moja
let tr=tblRudisha;tr.innerHTML='';
mi.forEach(t=>{
  let idx=DATA.miamala.indexOf(t);
  if(t.a=='marejesho' || t.a=='mkopo_toka'){
    tr.innerHTML+=`<tr><td>${t.w}</td><td contenteditable="true" onblur="let v=this.innerText.replace(/,/g,''); DATA.miamala[${idx}].k=+v; save()">${t.a=='mkopo_toka'?t.k:''}</td><td contenteditable="true" onblur="let v=this.innerText.replace(/,/g,''); DATA.miamala[${idx}].k=+v; save()">${t.a=='marejesho'?t.k:''}</td></tr>`;
  }
});

// DENI AGGREGATED - hii ndio itakayokwenda Excel
let dmap={}; DATA.miamala.forEach(t=>{if(!dmap[t.w])dmap[t.w]={c:0,r:0}; if(t.a=='mkopo_toka')dmap[t.w].c+=t.k; if(t.a=='marejesho')dmap[t.w].r+=t.k; });
let td=tblDeni;td.innerHTML='';
for(let j in dmap){let b=dmap[j].c-dmap[j].r; if(dmap[j].c>0||dmap[j].r>0){ td.innerHTML+=`<tr><td>${j}</td><td>${dmap[j].c.toLocaleString()}</td><td>${dmap[j].r.toLocaleString()}</td><td style="color:#c00;font-weight:800">${b.toLocaleString()}</td></tr>`; }}
let mu=muhtasari;mu.innerHTML='';let labels={hisa:'Hisa',ada:'Ada',marejesho:'Marejesho',faini:'Faini',fomu:'Fomu',michango:'Michango',mkopo_toka:'Mkopo Toka',dharura:'Dharura',mkono_pole:'Mkono wa Pole',matumizi_mengine:'Matumizi Mengine'};let sum={};document.querySelectorAll('#selAina option').forEach(o=>sum[o.value]=0);mi.forEach(t=>sum[t.a]+=t.k);for(let k in labels){if(sum[k]>0)mu.innerHTML+=`<div class="row"><span>${labels[k]}</span><b>${sum[k].toLocaleString()}</b></div>`}
let hi=hist;hi.innerHTML='';[...mi].reverse().forEach(t=>{let idx=DATA.miamala.indexOf(t);hi.innerHTML+=`<div class="row"><div><b>${t.w}</b> - ${t.a}<br><small>${t.t} - ${t.k.toLocaleString()}</small></div><div><button onclick="let v=prompt('Kiasi kipya:',DATA.miamala[${idx}].k);if(v!==null){DATA.miamala[${idx}].k=+v;save()}" style="border:0;background:#e3f2fd;border-radius:4px">✏️</button> <button onclick="if(confirm('Futa?')){DATA.miamala.splice(${idx},1);save()}" style="border:0;background:#ffebee;border-radius:4px;color:#c00">X</button></div></div>`});
}
// Hii ndio Excel inayosoma edits zako
function backup(){
  let m=curM(),y=curY();
  let dmap={}; DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y).forEach(t=>{if(!dmap[t.w])dmap[t.w]={c:0,r:0}; if(t.a=='mkopo_toka')dmap[t.w].c+=t.k; if(t.a=='marejesho')dmap[t.w].r+=t.k; });
  let csv='Jina,Mkopo Aliouchukua,Jumla Marejesho,Baki\n';
  for(let j in dmap){let b=dmap[j].c-dmap[j].r; csv+=`"${j}",${dmap[j].c},${dmap[j].r},${b}\n`; }
  if(Object.keys(dmap).length==0) csv+='"Hakuna data mwezi huu",0,0,0\n';
  let a=document.createElement('a');a.href=URL.createObjectURL(new Blob([csv],{type:'text/csv'}));a.download=`Uwashusema_Mahida_${m}_${y}.csv`;a.click();
}
function wa(){let txt=`*Uwashusema Mahida* ${curM()}/${curY()}\nMAPATO: ${tMapato.textContent}\nMATUMIZI: ${tMatumizi.textContent}\nSALIO: ${tSalio.textContent}\n\nWenye Deni:\n`;for(let r of tblDeni.rows){txt+=r.innerText+'\n'}window.open('https://wa.me/?text='+encodeURIComponent(txt),'_blank')}
mwezi.value=new Date().getMonth()+1;render();
</script></body></html>
HTML
cd ~/mbaochap && python -m http.server 8080
