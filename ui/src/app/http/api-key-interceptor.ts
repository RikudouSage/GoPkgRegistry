import {HttpInterceptorFn} from '@angular/common/http';
import {inject} from '@angular/core';
import {Auth} from '../services/auth';

export const ApiKeyInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(Auth);

  if (auth.isLoggedIn()) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${auth.getApiKey()}`,
      },
    });
  }

  return next(req);
}
