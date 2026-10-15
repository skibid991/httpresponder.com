
export interface RequestOptions {query?: Record<string, unknown>; method?: string; body?: BodyInit | Record<string, unknown>; headers?: HeadersInit; check?: boolean; redirect?: RequestRedirect}
export class APIError extends Error {status: number; body: string}
export class HTTPResponder {
  constructor(options?: {baseUrl?: string; timeout?: number; fetch?: typeof globalThis.fetch});
  url(path: string, query?: Record<string, unknown>): string;
  request(path: string, options?: RequestOptions): Promise<Response>;
  endpoint(path: string, options?: RequestOptions): Promise<Response>;
  decode(data: string, contentType?: string): Promise<Response>;
  encode(data: string): Promise<string>;
  respond(options?: {status?: number; body?: string; base64?: string; content_type?: string; headers?: Record<string, string>; delay?: number}): Promise<Response>;
  status(code: number | string): Promise<Response>;
  echo(options?: RequestOptions): Promise<Record<string, unknown>>;
}
