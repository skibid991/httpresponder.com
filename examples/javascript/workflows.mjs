import {HTTPResponder} from '../../sdk/javascript/index.js';
const client = new HTTPResponder();
const custom = await client.respond({status: 201, body: 'Created', headers: {'X-Example': 'yes'}});
console.log(custom.status, custom.headers.get('x-example'), await custom.text());
console.log(await client.echo({method: 'POST', body: {example: true}}));
const binary = await client.decode('AP+A', 'application/octet-stream');
console.log(new Uint8Array(await binary.arrayBuffer()));
