import { createHttpClient, type HttpClient, type HttpClientType } from './httpClient';

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || '/api';

// Default client (Axios-based)
export const api: HttpClient = createHttpClient('axios', apiBaseUrl);

// Alternative Fetch-based client
export const fetchApi: HttpClient = createHttpClient('fetch', apiBaseUrl);

// For custom use-cases, create a new client
export { createHttpClient, type HttpClient, type HttpClientType };
