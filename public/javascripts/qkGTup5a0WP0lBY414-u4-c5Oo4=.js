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


var Infinite = function Infinite() {
}

Infinite.prototype.Constant_ = {
    __idContainer: 'infinite',
    __idIndicator: 'indicator',
    __templateArray: ["\u003Ca href=/a/",
        "1:{{admin.id}}",
        "\u003E\u003Cdiv class=medium-list__image\u003E\u003C/div\u003E\u003Cdiv class=medium-list__footer\u003E\u003Cdiv class=\u0022medium-list__footer-container mobile-flex-vc\u0022\u003E\u003Cdiv class=\u0022medium-list__title line-clamp--2\u0022\u003E\u003Cspan\u003E",
        "3:{{admin.title}}",
        "\u003C/span\u003E\u003C/div\u003E\u003C/div\u003E\u003C/div\u003E\u003C/a\u003E"],
    __distance: 10
}

Infinite.prototype.prepare = function () {
    this._scrolling = false;
    this._elementContainer = document.getElementById(this.Constant_.__idContainer);
    if (!this._elementContainer) {
        return;
    }
    this._elementIndicator = document.getElementById(this.Constant_.__idIndicator);
    if (!this._elementIndicator) return;

    this._handler = this.handleScroll.bind(this);

    window.addEventListener('scroll', this._handler);
    if ("touchmove" in document) {
        document.addEventListener('touchmove', this._handler);
    }

}
Infinite.prototype.getScrollPos = function () {
    return window.pageYOffset;
}
Infinite.prototype.changeStatus = function (isScrolling) {
    this._scrolling = isScrolling;
}
Infinite.prototype.handleScroll = function () {
    if (this._scrolling) return;
    var scrollPos = this.getScrollPos();
    var pageHeight = document.documentElement.scrollHeight;
    var clientHeight = document.documentElement.clientHeight;

    // Check if scroll bar position is just 50px above the max, if yes, initiate an update
    if (pageHeight - (scrollPos + clientHeight) < this.__distance) {

        this.changeStatus(true);

        this.changeStatus(false);
    }
}

var infinite=new Infinite();
infinite.prepare();