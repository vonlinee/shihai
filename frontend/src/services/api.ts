import { createHttpClient, type HttpClient, type HttpClientType } from './httpClient';

// Default client (Axios-based)
export const api: HttpClient = createHttpClient('axios');

// Alternative Fetch-based client
export const fetchApi: HttpClient = createHttpClient('fetch');

// For custom use-cases, create a new client
export { createHttpClient, type HttpClient, type HttpClientType };
