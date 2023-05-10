const baseURL=''; //same domain
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

function optionsDisplay(thingId){
	const subm = document.getElementsByClassName('id-'+thingId)[0];
	togDisplay(subm.getElementsByClassName('options_expando')[0]);
			
}
function isHidden(elem){
	return elem.classList.contains('hidden');
}

function isShown(elem){
	return !isHidden(elem);
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

function userTextSubm(uf){
	
	const json = form2json(uf);
	const errDisplay = uf.getElementsByClassName('error')[0];
	errDisplay.textContent ='Loading...';
	pst('/post/edit',json)
		.then(r => r.json().then(j =>{
			if (r.status === 200) {
				errDisplay.textContent ='';
				uf.getElementsByClassName('md')[0].innerHTML = j.selftext_html;
				const entry = uf.parentNode.parentNode;
				hide(entry.getElementsByClassName('md-container')[0]);
				
			} else {
				showErrors(r.status, j.json.errors, errDisplay)
			}
		}));
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
function docOnLoad(){
	const expandoButtons = document.getElementsByClassName('expando-button');
	
	
	const optionButtons = document.getElementsByClassName('options_link');
	
	for (let i = 0;i < optionButtons.length;i++){
		const b = optionButtons[i];
		const id = b.dataset.id;
		b.addEventListener("click",e => {
			e.preventDefault();
			optionsDisplay(id);
		});
	}
	

	
	const links = document.getElementsByClassName('link');
	for (link of links){
		
		const uf = link.getElementsByClassName('usertext')[0];
		uf.addEventListener("submit",e => {
			e.preventDefault();
			userTextSubm(uf);
			
		});
	
	}
	
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