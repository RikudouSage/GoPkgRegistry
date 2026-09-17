import { isPlatformServer } from '@angular/common';
import { HttpInterceptorFn } from '@angular/common/http';
import { inject, PLATFORM_ID } from '@angular/core';
import { environment } from '../../environments/environment';

export const SsrApiInterceptor: HttpInterceptorFn = (req, next) => {
  const platformId = inject(PLATFORM_ID);
  const apiUrl = environment.apiUrl.replace(/\/$/, '');
  const isApiRequest = req.url === apiUrl || req.url.startsWith(`${apiUrl}/`);

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
