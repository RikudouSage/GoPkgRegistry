import { environment } from '../../environments/environment';

export function apiUrl(path: string): string {
  const baseUrl = environment.apiUrl.replace(/\/+$/, '');
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  return `${baseUrl}${normalizedPath}`;
}
