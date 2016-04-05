require('./less/style.less');
var Elm = require('./elm/Main');
Elm.embed(Elm.Main, document.getElementById('main'), {swap: false} );
