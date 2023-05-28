const baseURL = ''; //same domain
var moreURL = '';
var sorting = '';
var inProgress = false;
const doc = document.documentElement;
function unimpl(){
	alert('Unimplemented!!!');
	return false;
}
function hide(elem){
	elem.classList.add('hidden');
}

function show(elem){
	elem.classList.remove('hidden');
}
function togDisplay(elem){
	if (elem.classList.contains('hidden')) {
		show(elem);
	} else {
		hide(elem);
	}
}
function togTopMenu(){
	togDisplay(document.getElementById('top_menu'));
	return false;
}

function backgroundUnread(){
	if (document.getElementById('msgTrigger')){
		return;
	}
	fetch(baseURL+'/checkunread/', {method: 'HEAD'}).then(r =>{
		if (r.status === 200){
			const mail = document.getElementById('mail');
			const mailcl = mail.classList;
			mail.href='/message/unread/';
			mailcl.remove('nohavemail');
			mailcl.add('havemail');
		}
	});
}

function get(url){
	return fetch(baseURL+url, {method: 'GET'});
}

function pst(url,json){
	return fetch(baseURL+url, {
		method: 'POST',
		headers: {
			'Accept': 'application/json',
			'Content-Type': 'application/json',
		},
		body: json
	
	});
}

function isHidden(elem){
	return elem.classList.contains('hidden');
}

function isShown(elem){
	return !isHidden(elem);
}

function optionsDisplay(id){
	togDisplay(document.getElementById(id).getElementsByClassName('options_expando')[0]);
}

function editDisplay(that){
	const thing = that.parentNode.parentNode;
	const expando = thing.getElementsByClassName('expando')[0];
	
	if (isHidden(expando)){
		show(expando);
	}
	togDisplay(expando.getElementsByClassName('usertext-edit')[0]);
	togDisplay(expando.getElementsByClassName('usertext-body')[0]);
	return false;
	
}
function editFormDisplay(that){
	const uf = that.parentNode.parentNode.getElementsByClassName('usertext')[0];
	togDisplay(uf.getElementsByClassName('md')[0]);
	togDisplay(uf.getElementsByClassName('usertext-edit')[0]);
	return false;
}
function appendSelected(replytext){
	const strs = window.getSelection().toString().split("\r\n\r\n");
	for (str of strs) {
		if (str && str.length>0){
			console.log(replytext);
			replytext.value += '> '+str.trim();
			replytext.value += "\r\n\r\n";
		}
	}
}

function replyFormDisplay(that){
	const crf = that.parentNode.parentNode.parentNode.getElementsByClassName('usertext')[1];
	
	togDisplay(crf);
	appendSelected(crf.getElementsByClassName('replytext')[0]);
	return false;
}

function msgReplyFormDisplay(that){
	const mrf = that.parentNode.parentNode.getElementsByClassName('usertext')[0];
	console.log(mrf);
	togDisplay(mrf);
	return false;
}

function postBodyDisplay(that){
	const thing = that.parentNode.parentNode;
	const expando = thing.getElementsByClassName('expando')[0];
	if (expando){
		
		const utEdit = expando.getElementsByClassName('usertext-edit')[0];
		if (utEdit){
			hide(utEdit);
		}
		const utBody = expando.getElementsByClassName('usertext-body')[0];
		if (utBody){
			show(utBody);
		}
		togDisplay(expando);
		
	}	
	
	if(that.classList.contains('expanded')){
		that.classList.remove('expanded');
		that.classList.add('collapsed');
	} else {
		that.classList.remove('collapsed');
		that.classList.add('expanded');
	}
	return false;
	
}

function form2json(form){
	const fd = new FormData(form);
	const data ={};
	for (let [name, val] of fd) {
		data[name] = val;
	}
	return JSON.stringify(data);
}

function showErrors(status, errors, errDisplay){
	let s = 'Status: '+status+'. ';
		for (error of errors){
			s += error.Reason;
			s+='<br/>';
		}
		errDisplay.innerHTML = s;
}

function togCollapsed(comment){
	
	const ccl = comment.classList;
	if (ccl.contains('collapsed')){
		
		ccl.remove('collapsed');
		comment.title = '';
		console.log(comment);
	} else {
		
		ccl.add('collapsed');
		comment.title = 'Double click to show';
		comment.addEventListener('dblclick', e => {
			togCollapsed(comment);

		});
	}
}

function collapseComment(that){
	const comment = that.parentNode.parentNode.parentNode;
	comment.classList.add('collapsed');
	comment.title = 'double click to show collapsed comment';
	comment.addEventListener('dblclick', e => togCollapsed(comment));
	return false;
}

function editSubmit(uf){
	const json = form2json(uf);
	const errDisplay = uf.getElementsByClassName('error')[0];
	errDisplay.textContent ='Loading...';
	pst('/edit/',json)
		.then(r => r.json().then(j =>{
			if (r.status === 200) {
				errDisplay.textContent ='';
				const tcl = uf.parentNode.parentNode.classList;
				console.log(j);
					uf.getElementsByClassName('usertext-body')[0].innerHTML = j.body;
				
				hide(uf.getElementsByClassName('usertext-edit')[0]);
				show(uf.getElementsByClassName('usertext-body')[0]);
				
			} else {
				showErrors(r.status, j.json.errors, errDisplay)
			}
		}));
		return false;
}

function submitComment(pf){
	const json = form2json(pf);
	const errDisplay = pf.getElementsByClassName('error')[0];
	errDisplay.textContent ='Loading...';
	pst('/comment/',json)
		.then(r => {
			if (r.status === 200) {
				errDisplay.textContent ='';
				pf.getElementsByTagName('textarea')[0].value='';
				r.text().then(t => {
					const par = pf.parentNode;
					console.log(par);
					const pcl = par.classList;
					if (pcl.contains('commentarea')){
					document.getElementsByClassName('nestedlisting')[0]
						.insertAdjacentHTML('afterbegin',t);
					} else if(pcl.contains('comment')||pcl.contains('was_comment')){
						const child = document.createElement('div');
						child.classList.add('child');
						child.innerHTML = t;
						par.appendChild(child);
						togDisplay(pf);
					} else if (pcl.contains('was_comment')){
						
					} else {
						//reply to dm
					}
					
				});
			} else {
				r.json().then(j => {
					showErrors(r.status, j.json.errors, errDisplay)
				});
			}
		});
		return false;
}

function vote(id, direction){
	const thing = document.getElementById(id);
	const midcol = thing.getElementsByClassName('midcol')[0];
	get('/vote/'+direction+'/'+id+'/').then(r => {
		if (r.status === 200){
			const arrows = midcol.getElementsByClassName('arrow');
			const ecl = thing.getElementsByClassName('entry')[0].classList;
			switch(direction){
				case 'remove':
					ecl.remove('likes');
					ecl.remove('dislikes');
					ecl.add('unvoted');
					arrows[0].classList.remove('upmod');
					arrows[1].classList.remove('downmod');
					break;
				case 'up':
					ecl.remove('dislikes');
					ecl.remove('unvoted');
					ecl.add('likes');
					arrows[0].classList.add('upmod');
					arrows[1].classList.remove('downmod');
					break;
				case 'down':
					ecl.remove('unvoted');
					ecl.remove('likes');
					ecl.add('dislikes');
					arrows[0].classList.remove('upmod');
					arrows[1].classList.add('downmod');
					break;
			}
		} else {
			alert('Voting error '+r.status);
		}
	});
	return false;
}

function scrolling(){
	const sTop = doc.scrollTop;
	const sHeight = doc.scrollHeight;
	const cHeight = doc.clientHeight;
	const aftercontainers = document.getElementsByClassName('aftercontainer');
	const after = aftercontainers[aftercontainers.length-1].value;
	if (!inProgress && sTop + cHeight >= sHeight && after && after.length > 1){
		inProgress = true;
		const loader = document.getElementById('loaderimg');
		togDisplay(loader);
		get(moreURL+'?after='+after+'&sort='+sorting).then(r => {
			inProgress = false;
			hide(loader);
			if (r.status === 200) {
				r.text().then(t => {
					console.log('here');
					document.getElementById('siteTable')
						.insertAdjacentHTML('beforeend', t);
				});
			} else {
				alert('Error fetching more data: '+r.status);
			}
		});
	}
}
function docOnLoad(){
	backgroundUnread();
	const things = document.getElementsByClassName('thing');
	for (thing of things){
		const thing_id = thing.id;
		
		//Expando options display button
		const optionsLink = thing.getElementsByClassName('options_link')[0]
		if (optionsLink){
			optionsLink.addEventListener('click', e => {
				e.preventDefault();
				optionsDisplay(thing_id);
			});
		}
		
		//Up and downvote buttons
		const midcol = thing.getElementsByClassName('midcol')[0];
		if (midcol){
			const up = midcol.getElementsByClassName('up')[0];
			if (up.classList.contains('upmod')){
				up.addEventListener('click', e => {
					e.preventDefault();
					vote(thing_id,'remove');
				});
			} else {
				up.addEventListener('click', e => {
					e.preventDefault();
					vote(thing_id,'up');
				});
			}
			
			const down = midcol.getElementsByClassName('down')[0];
			if (down.classList.contains('downmod')){
				down.addEventListener('click', e => {
					e.preventDefault();
					vote(thing_id,'remove');
				});
			} else {
				down.addEventListener('click', e => {
					e.preventDefault();
					vote(thing_id,'down');
				});
			}
			
		}
	}
	
}

document.addEventListener("DOMContentLoaded", docOnLoad);