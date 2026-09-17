import { isPlatformServer } from '@angular/common';
import { HttpInterceptorFn } from '@angular/common/http';
import { inject, PLATFORM_ID } from '@angular/core';
import { environment } from '../../environments/environment';

export const SsrApiInterceptor: HttpInterceptorFn = (req, next) => {
  const platformId = inject(PLATFORM_ID);
  const apiBaseUrl = environment.apiUrl.replace(/\/+$/, '');
  const isApiRequest = apiBaseUrl === ''
    ? req.url === '/admin' || req.url.startsWith('/admin/')
    : req.url === apiBaseUrl || req.url.startsWith(`${apiBaseUrl}/`);

  if (environment.ssrApiUrl && isPlatformServer(platformId) && isApiRequest) {
    const ssrUrl = new URL(environment.ssrApiUrl);
    const url = new URL(req.url, ssrUrl.origin);

    url.host = ssrUrl.host;
    url.protocol = ssrUrl.protocol;

    req = req.clone({
      url: url.toString(),
    });
  }

  return next(req);
};
