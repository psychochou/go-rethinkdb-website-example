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
    this._count = 0;
    this._disable = false;
    this._elementContainer = document.getElementById(this.Constant_.__idContainer);
    if (!this._elementContainer) {
        return;
    }
    this._elementIndicator = document.getElementById(this.Constant_.__idIndicator);
    if (!this._elementIndicator) return;

    this._url = window.location.pathname.startsWith("/a/") ? "/api/list/" : "/api/lt?tag=" + window.location.pathname.split('/')[window.location.pathname.split('/').length - 1] + "&skip="
    this._handler = this.handleScroll.bind(this);

    window.addEventListener('scroll', this._handler);
    if ("ontouchmove" in document) {
        document.addEventListener('touchmove', this._handler);
    }

}
Infinite.prototype.getScrollPos = function () {
    return window.pageYOffset;
}
Infinite.prototype.changeStatus = function (isScrolling) {
    this._scrolling = isScrolling;
    if (isScrolling) {
        this._elementIndicator.style.display = "block";
    } else {
        this._elementIndicator.style.display = "none";
    }
}
Infinite.prototype.append = function (obj) {
    for (var index = 0; index < obj.length; index++) {
        var element = obj[index];
        var div = document.createElement('div');
        div.className = "medium-list__item";
        this.Constant_.__templateArray[1] = element["id"];
        this.Constant_.__templateArray[3] = element["title"];
        this._elementContainer.appendChild(div);
        div.innerHTML = this.Constant_.__templateArray.join('');
    }

}
Infinite.prototype.dispose = function () {
    this._disable = true;
    window.removeEventListener('scroll', this._handler);
    if ("ontouchmove" in document) {
        document.removeEventListener('touchmove', this._handler);
    }
    this._elementIndicator.style.display = "none";

}
Infinite.prototype.handleScroll = function () {
    if (this._disable || this._scrolling) return;
    var scrollPos = this.getScrollPos();
    var pageHeight = document.documentElement.scrollHeight;
    var clientHeight = document.documentElement.clientHeight;

    // Check if scroll bar position is just 50px above the max, if yes, initiate an update
    if (pageHeight - (scrollPos + clientHeight) < this.Constant_.__distance) {
        this.changeStatus(true);
        this._count++;
        var _this = this;
        fetch(this._url + (this._count * 10)).then(function (res) {
            return res.text();
        }).then(function (value) {
            if (value == null) {
                _this.dispose();
            } else {
                var obj = JSON.parse(value);
                _this.append(obj);
                _this.changeStatus(false);
            }


        }).catch(function () {
            _this.dispose();
            console.log(arguments)
        })
    }
}

var infinite = new Infinite();
infinite.prepare();