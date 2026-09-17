import {Injectable} from '@angular/core';
import {TranslocoService} from '@jsverse/transloco';
import {map, Observable, switchMap} from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class Translator {
  public constructor(
    private readonly transloco: TranslocoService,
  ) {
  }

  public get(key: string, params: Record<string, any> = {}): Observable<string> {
    return this.transloco.langChanges$.pipe(
      switchMap(lang => this.transloco.load(lang).pipe(
        map (() => this.transloco.translate(key, params)),
      )),
    );
  }
}
