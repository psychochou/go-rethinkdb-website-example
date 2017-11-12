var Application = function Application() {

}

Application.prototype.ajustMenuStatus = function () {

    var navItems = document.querySelectorAll('.nav__item');
    if (!navItems.length) {
        return;
    }
    var pathname = window.location.pathname
    var match = /\/t\/([^\/]*)/.exec(pathname)
    if (match) {
        for (var index = 0; index < navItems.length; index++) {
            var element = navItems[index];
            if (element.getAttribute('href').endsWith(match[1])) {
                element.classList.add('active')
            } else {
                element.classList.remove('active')
            }
        }
    }

}

var app = new Application();
app.ajustMenuStatus();

; (function () {


    document.querySelector('.article-nav') && window.addEventListener('hashchange', function () {
        var n = document.querySelector(":target");
        n && window.scrollTo(0, n.offsetTop - 45)
    })


})()