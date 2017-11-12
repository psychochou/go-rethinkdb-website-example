function copyTextToClipboard(text) {
  var copyFrom = document.createElement("textarea");
  copyFrom.textContent = text;
  var body = document.getElementsByTagName('body')[0];
  body.appendChild(copyFrom);
  copyFrom.select();
  document.execCommand('copy');
  body.removeChild(copyFrom);
}


function prepare() {

  window.addEventListener('keyup', function (event) {
    var element;
	
	console.log(event)
	 //if (event.key == "z" && (element = document.body)) {
   if (event.keyCode==90 && (element = document.querySelector(".detail-toc"))) {
      
      copyTextToClipboard(element.outerHTML);
    }else if(event.keyCode==88 && (element=document.querySelectorAll('.tab_set_info a'))){
		var linkString='';
		element.forEach((ele)=>{
			linkString+=ele.getAttribute('href')+'\r\n'
			
		});
		copyTextToClipboard(linkString);
	}
  })
}
prepare();