import { inject, Injectable } from '@angular/core';
import { Translation, TranslocoLoader } from '@jsverse/transloco';
import { HttpClient } from '@angular/common/http';
import {Observable} from 'rxjs';
import {environment} from '../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class TranslocoHttpLoader implements TranslocoLoader {
  public constructor(
    private readonly httpClient: HttpClient,
  ) {
  }

  public getTranslation(lang: string): Observable<Translation> {
    const baseHref = environment.baseHref.replace(/\/+$/, '');

    return this.httpClient.get<Translation>(`${baseHref}/i18n/${lang}.json`);
  }
}
