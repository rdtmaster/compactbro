const baseURL = ''; //same domain
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
	fetch(baseURL+'/checkunread/', {method: 'HEAD'}).then(r =>{
		if (r.status === 200){
			const mailcl = document.getElementById('mail').classList;
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

function optionsDisplay(that){
	togDisplay(that.parentNode.getElementsByClassName('options_expando')[0]);
	return false;
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
function replyFormDisplay(that){
	const thing = that.parentNode.parentNode.parentNode;
	togDisplay(thing.getElementsByClassName('commentreply')[0]);
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
	pst('/comment',json)
		.then(r => {
			if (r.status === 200) {
				errDisplay.textContent ='';
				pf.getElementsByTagName('textarea')[0].value='';
				r.text().then(t => {
					document.getElementsByClassName('nestedlisting')[0]
						.insertAdjacentHTML('afterbegin',t);
					
				});
			} else {
				r.json().then(j => {
					showErrors(r.status, j.json.errors, errDisplay)
				});
			}
		});
}

function vote(that){
	const cl = that.classList
	const midcol = that.parentNode;
	const thing = midcol.parentNode;
	var direction = '';
	if (cl.contains('downmod') || cl.contains('upmod')){
		direction ="remove";
	} else if (cl.contains('up')){
		direction = 'up';
	} else if (cl.contains('down')){
		direction = 'down';
	}
	get('/vote/'+direction+'/'+thing.dataset.id+'/').then(r => {
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
function docOnLoad(){
	backgroundUnread();
	const links = document.getElementsByClassName('link');
	
	
	const postCommentForm = document.getElementById('postComment');
	if (postCommentForm){
		console.log('triggered');
		postCommentForm.addEventListener("submit",e => {
			e.preventDefault();
			submitComment(postCommentForm);
		});
	}
}

document.addEventListener("DOMContentLoaded", docOnLoad);