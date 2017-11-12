const r = require('rethinkdb');
const util = require('util');

var query =
   r.db('psycho').table('articles').getAll("health",{index:"tags"}).orderBy(r.desc("createAt")).pluck("id","title","image").skip(0).limit(10).run();

console.log(JSON.stringify(query));