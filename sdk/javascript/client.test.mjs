import test from 'node:test';
import assert from 'node:assert/strict';
import {HTTPResponder, APIError} from './index.js';
test('query encoding preserves nested URLs, Unicode, base64 and JSON objects', () => {
 const c = new HTTPResponder(); const target='https://example.com/a?x=1&y=two#part';
 const u = new URL(c.url('/test',{url:target,data:'+/= 🌎',enabled:false,headers:{'X-Test':'yes'},skip:undefined}));
 assert.equal(u.searchParams.get('url'),target); assert.equal(u.searchParams.get('data'),'+/= 🌎');
 assert.equal(u.searchParams.get('enabled'),'false'); assert.deepEqual(JSON.parse(u.searchParams.get('headers')),{'X-Test':'yes'}); assert(!u.searchParams.has('skip'));
});
test('JSON bodies, authorization, timeout and empty 204 responses', async () => {
 const c=new HTTPResponder({fetch:async (url,o)=>{assert.equal(o.method,'POST');assert.equal(o.headers.get('authorization'),'Bearer sample');assert.equal(o.headers.get('content-type'),'application/json');assert.deepEqual(JSON.parse(o.body),{text:'🌎'});assert(o.signal);return new Response(null,{status:204});}});
 assert.equal((await c.request('/test',{method:'POST',body:{text:'🌎'},headers:{Authorization:'Bearer sample'}})).status,204);
});
test('binary bytes are preserved and error responses remain inspectable', async () => {
 const c=new HTTPResponder({fetch:async()=>new Response(new Uint8Array([0,255,128]),{status:418})});
 await assert.rejects(c.request('/test'),e=>e instanceof APIError&&e.status===418);
 const r=await c.request('/test',{check:false});assert.deepEqual([...new Uint8Array(await r.arrayBuffer())],[0,255,128]);
});
test('redirects stay visible and protocol-relative paths are rejected', async () => {
 const c=new HTTPResponder({fetch:async(u,o)=>{assert.equal(o.redirect,'manual');return new Response(null,{status:302,headers:{location:'https://example.com/'}});}});
 assert.equal((await c.request('/test',{check:false})).status,302);
 assert.throws(()=>c.url('//example.com/')); assert.throws(()=>c.url('https://example.com/'));
});
