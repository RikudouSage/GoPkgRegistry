import { inject, Injectable } from '@angular/core';
import { Translation, TranslocoLoader } from '@jsverse/transloco';
import { HttpClient } from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TranslocoHttpLoader implements TranslocoLoader {
  public constructor(
    private readonly httpClient: HttpClient,
  ) {
  }

  public getTranslation(lang: string): Observable<Translation> {
    return this.httpClient.get<Translation>(`/i18n/${lang}.json`);
  }
}
