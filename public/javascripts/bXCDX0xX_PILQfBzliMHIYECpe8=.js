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

function applyInifinite() {
    var infinite = document.getElementById('infinite');

    if (!infinite) return;

    var isIE = /msie/gi.test(navigator.userAgent); // http://pipwerks.com/2011/05/18/sniffing-internet-explorer-via-javascript/

    var infiniteScroll = function (options) {
        var defaults = {
            callback: function () { },
            distance: 50
        }
        // Populate defaults
        for (var key in defaults) {
            if (typeof options[key] == 'undefined') options[key] = defaults[key];
        }

        var scroller = {
            options: options,
            updateInitiated: false
        }

        window.onscroll = function (event) {
            handleScroll(scroller, event);
        }
        // For touch devices, try to detect scrolling by touching
        document.ontouchmove = function (event) {
            handleScroll(scroller, event);
        }
    }

    function getScrollPos() {
        // Handle scroll position in case of IE differently
        if (isIE) {
            return document.documentElement.scrollTop;
        } else {
            return window.pageYOffset;
        }
    }

    var prevScrollPos = getScrollPos();

    // Respond to scroll events
    function handleScroll(scroller, event) {
        if (scroller.updateInitiated) {
            return;
        }
        var scrollPos = getScrollPos();
        if (scrollPos == prevScrollPos) {
            return; // nothing to do
        }

        // Find the pageHeight and clientHeight(the no. of pixels to scroll to make the scrollbar reach max pos)
        var pageHeight = document.documentElement.scrollHeight;
        var clientHeight = document.documentElement.clientHeight;

        // Check if scroll bar position is just 50px above the max, if yes, initiate an update
        if (pageHeight - (scrollPos + clientHeight) < scroller.options.distance) {
            scroller.updateInitiated = true;

            scroller.options.callback(function () {
                scroller.updateInitiated = false;
            });
        }

        prevScrollPos = scrollPos;
    }

    var options = {
        distance: 10,
        callback: function (done) {
            // 1. fetch data from the server
            // 2. insert it into the document
            // 3. call done when we are done
            console.log(1);
            done();
        }
    }

    // setup infinite scroll
    infiniteScroll(options);
}

applyInifinite();
