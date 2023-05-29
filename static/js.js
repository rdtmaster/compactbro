const baseURL = ''; //same domain
var moreURL = '';
var sorting = '';
var inProgress = false;
const doc = document.documentElement;
function form2json(form){
	const fd = new FormData(form);
	const data ={};
	for (let [name, val] of fd) {
		data[name] = val;
	}
	return JSON.stringify(data);
}

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

function eventsOfContainer(container){
	const things = container.getElementsByClassName('thing');
	for (thing of things){
		setEvents(thing);
	}
}

// set events for a thing
function setEvents(thing){
	const thing_id = thing.id;
	const usertext = thing.getElementsByClassName('usertext')[0];
	if (usertext){
		usertext.addEventListener('submit', e => {
			e.preventDefault();
			editSubmit(thing_id);
		});
	}
	//Expando options display button
	const optionsLink = thing.getElementsByClassName('options_link')[0]
	if (optionsLink){
		optionsLink.addEventListener('click', e => {
			e.preventDefault();
			optionsDisplay(thing_id);
		});
	}
	
	
	if (thing.classList.contains('link')){
		//Selftext expando
		const exButton = thing.getElementsByClassName('expando-button')[0];
		if (exButton){
			exButton.addEventListener('click', e => {
				e.preventDefault();
				postBodyDisplay(thing_id);
			});
		}
	}
	
	if (thing.classList.contains('comment')){
		const ut = thing.getElementsByClassName('usertext')[1];
		if (ut){
			ut.addEventListener('submit', e => {
			e.preventDefault();
			submitComment(thing_id);
		});
		}
		
	}
	//icon expando buttons
	const opts = thing.getElementsByClassName('options_expando')[0];
	const editButton = opts.getElementsByClassName('edit-icon')[0];
	if (editButton){
		editButton.parentNode.addEventListener('click', e => {
			e.preventDefault();
			editDisplay(thing_id);
		});
	}
	
	const replyButton = opts.getElementsByClassName('reply-icon')[0];
	if (replyButton){
		replyButton.parentNode.addEventListener('click', e => {
			e.preventDefault();
			replyFormDisplay(thing_id);
		});
	}
	const shareButton = opts.getElementsByClassName('email-icon')[0];
	shareButton.parentNode.addEventListener('click', e => {
		e.preventDefault();
		unimpl();
	});
	
	
	
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



//------- Event Listeners -------

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

function editDisplay(id){
	const thing = document.getElementById(id);
	const tcl = thing.classList;
	if (tcl.contains('link')){
		const expando = thing.getElementsByClassName('expando')[0];
		if (isHidden(expando)){
			show(expando);
		}
		togDisplay(expando.getElementsByClassName('usertext-edit')[0]);
		togDisplay(expando.getElementsByClassName('usertext-body')[0]);
	const exButton = thing.getElementsByClassName('expando-button')[0]
	if (exButton){
			if(exButton.classList.contains('expanded')){
				exButton.classList.remove('expanded');
				exButton.classList.add('collapsed');
			} else {
				exButton.classList.remove('collapsed');
				exButton.classList.add('expanded');
			}
		}
	} else if (tcl.contains('comment')){
		const uf = thing.getElementsByClassName('entry')[0].getElementsByClassName('usertext')[0];
		togDisplay(uf.getElementsByClassName('md')[0]);
		togDisplay(uf.getElementsByClassName('usertext-edit')[0]);
	}
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
	replytext.focus();
}

function replyFormDisplay(thing_id){
	const crf = document.getElementById(thing_id).getElementsByClassName('usertext')[1];
	togDisplay(crf);
	appendSelected(crf.getElementsByClassName('replytext')[0]);
	
}

function msgReplyFormDisplay(that){
	const mrf = that.parentNode.parentNode.getElementsByClassName('usertext')[0];
	console.log(mrf);
	togDisplay(mrf);
	return false;
}

function postBodyDisplay(id){
	const thing = document.getElementById(id);
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
	const exButton = thing.getElementsByClassName('expando-button')[0]
	if(exButton.classList.contains('expanded')){
		exButton.classList.remove('expanded');
		exButton.classList.add('collapsed');
	} else {
		exButton.classList.remove('collapsed');
		exButton.classList.add('expanded');
	}
}

function showErrors(status, errors, errDisplay){
	let s = 'Status: '+status+'. ';
		for (error of errors){
			s += error.Reason;
			s+='<br/>';
		}
		errDisplay.innerHTML = s;
}


function collapseComment(that){
	const comment = that.parentNode.parentNode.parentNode;
	comment.classList.add('collapsed');
	comment.title = 'double click to show collapsed comment';
	comment.addEventListener('dblclick', e => togCollapsed(comment));
	return false;
}

function editSubmit(id){
	const uf = document.getElementById(id).getElementsByClassName('usertext')[0];
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
}

function submitComment(id){
	const thing = document.getElementById(id);
	const pf = thing.getElementsByClassName('usertext')[1];
	const json = form2json(pf);
	const errDisplay = pf.getElementsByClassName('error')[0];
	errDisplay.textContent ='Loading...';
	pst('/comment/',json)
		.then(r => {
			if (r.status === 200) {
				errDisplay.textContent ='';
				pf.getElementsByTagName('textarea')[0].value='';
				r.text().then(t => {
					const pcl = thing.classList;
					const child = document.createElement('div');
					child.innerHTML = t;
					eventsOfContainer(child);
					if (id === 'commentArea'){
						const listing = document.getElementsByClassName('nestedlisting')[0];
						listing.insertBefore(child, listing.firstChild);
					} else if(pcl.contains('comment')||pcl.contains('was_comment')){
						child.classList.add('child');
						thing.appendChild(child);
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
	if (inProgress){
		return;
	}
	const sTop = doc.scrollTop;
	const sHeight = doc.scrollHeight;
	const cHeight = doc.clientHeight;
	const aftercontainers = document.getElementsByClassName('aftercontainer');
	const after = aftercontainers[aftercontainers.length-1].value;
	if (sTop + cHeight >= sHeight && after && after.length > 1){
		inProgress = true;
		const loader = document.getElementById('loaderimg');
		togDisplay(loader);
		get(moreURL+'?after='+after+'&sort='+sorting).then(r => {
			inProgress = false;
			hide(loader);
			if (r.status === 200) {
				r.text().then(t => {
					const addition = document.createElement('div');
					addition.innerHTML = t;
					eventsOfContainer(addition);
					document.getElementById('siteTable').appendChild(addition);
				});
			} else {
				alert('Error fetching more data: '+r.status);
			}
		});
	}
}
function docOnLoad(){
	backgroundUnread();
	commentarea = document.getElementById('commentArea');
	if (commentarea){
		ut = commentarea.getElementsByClassName('usertext')[1];
		if (ut){
			ut.addEventListener('submit', e => {
				e.preventDefault();
				submitComment('commentArea');
			});
		}
	}
	eventsOfContainer(document);
	
}

document.addEventListener("DOMContentLoaded", docOnLoad);