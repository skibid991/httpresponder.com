import {HTTPResponder} from '../../sdk/javascript/index.js';
const client = new HTTPResponder();
console.log(await (await client.decode('SGVsbG8sIHdvcmxkIQ==')).text());
const response = await client.status(418);
console.log('Requested status:', response.status);
