const baseURL=''; //same domain
function unimpl(){
	alert('Unimplemented!!!');
}
function hide(elem){
	elem.style.display = 'none';
}

function show(elem){
	elem.style.display = 'block';
}
function togDisplay(elem){
	if (elem.style.display === 'none') {
		elem.style.display = 'block';
	} else {
		elem.style.display = 'none';
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

function editDisplay(that){
	const thing = that.parentNode.parentNode;
	const expando = thing.getElementsByClassName('expando')[0];
	if (expando.style.display == 'none'){
		show(expando);
	}
	togDisplay(expando.getElementsByClassName('usertext-edit')[0]);
	
	
	/*
	
	
	togDisplay(userTextBody.getElementsByClassName('md')[1]);
	togDisplay(userTextBody);
	
	
	togDisplay(userTextBody.getElementsByClassName('md')[0]);
	const btn = subm.getElementsByClassName('selftext')[0];
	if(btn.classList.contains('expanded')){
		btn.classList.remove('expanded');
		btn.classList.add('collapsed');
	} else {
		btn.classList.remove('collapsed');
		btn.classList.add('expanded');
	}
	*/
	
}
function userTextDisplay(thingId){
	const subm = document.getElementsByClassName('id-'+thingId)[0];
	const userTextBody = subm.getElementsByClassName('expando')[0];
	const btn = subm.getElementsByClassName('selftext')[0];
	togDisplay(userTextBody);
	if(btn.classList.contains('expanded')){
		btn.classList.remove('expanded');
		btn.classList.add('collapsed');
	} else {
		btn.classList.remove('collapsed');
		btn.classList.add('expanded');
	}
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
	for (let i = 0;i < expandoButtons.length;i++){
		const b = expandoButtons[i];
		const id = b.dataset.id;
		b.addEventListener("click",e => {
			e.preventDefault();
			userTextDisplay(id);
			
		});
	}
	
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
		const editButton = link.getElementsByClassName('editButton')[0];
		if(editButton){
			editButton.addEventListener("click",e => {
				e.preventDefault();
				editDisplay(editButton);
			});
		}
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