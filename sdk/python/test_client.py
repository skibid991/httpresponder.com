import unittest
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from threading import Thread
from urllib.parse import urlsplit, parse_qs
from httpresponder import HTTPResponder, APIError

class Handler(BaseHTTPRequestHandler):
    def log_message(self, *args): pass
    def do_GET(self):
        path=urlsplit(self.path).path
        if path=='/redirect':
            self.send_response(302); self.send_header('Location','/ok'); self.end_headers(); return
        self.send_response(418 if path=='/error' else 200); self.end_headers(); self.wfile.write(bytes([0,255,128]))
    def do_POST(self):
        body=self.rfile.read(int(self.headers.get('Content-Length','0')))
        self.send_response(200); self.send_header('X-Auth',self.headers.get('Authorization','')); self.end_headers(); self.wfile.write(body)

class ClientTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.server=ThreadingHTTPServer(('127.0.0.1',0), Handler)
        cls.thread=Thread(target=cls.server.serve_forever,daemon=True);cls.thread.start()
    @classmethod
    def tearDownClass(cls):
        cls.server.shutdown();cls.server.server_close();cls.thread.join()
    def client(self): return HTTPResponder(base_url=f'http://127.0.0.1:{self.server.server_port}')
    def test_encoding(self):
        target='https://example.com/a?x=1&y=two#part'
        q=parse_qs(urlsplit(self.client().url('/test',{'url':target,'data':'+/= 🌎','enabled':False,'skip':None})).query)
        self.assertEqual(q,{'url':[target],'data':['+/= 🌎'],'enabled':['false']})
    def test_binary(self): self.assertEqual(self.client().request('/ok').body,bytes([0,255,128]))
    def test_json_and_auth(self):
        r=self.client().request('/test',method='POST',body={'text':'🌎'},headers={'Authorization':'Bearer sample'})
        self.assertEqual(r.json(),{'text':'🌎'});self.assertEqual(r.headers['X-Auth'],'Bearer sample')
    def test_errors_and_redirects(self):
        with self.assertRaises(APIError) as caught:self.client().request('/error')
        self.assertEqual(caught.exception.status,418)
        self.assertEqual(self.client().request('/error',check=False).status,418)
        self.assertEqual(self.client().request('/redirect',check=False).status,302)
    def test_bad_path(self):
        with self.assertRaises(ValueError):self.client().url('//example.com/')

if __name__=='__main__':unittest.main()
