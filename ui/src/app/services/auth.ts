import {DOCUMENT, Inject, Injectable, Optional, PLATFORM_ID, REQUEST} from '@angular/core';
import {HttpBackend, HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {catchError, firstValueFrom, map, of} from 'rxjs';
import {isPlatformBrowser} from '@angular/common';

@Injectable({
  providedIn: 'root',
})
export class Auth {
  private readonly cookieName = 'apiKey';
  private readonly httpClient: HttpClient;
  private readonly isBrowser: boolean;

  public constructor(
    httpBackend: HttpBackend,
    @Inject(DOCUMENT)
    private readonly document: Document,
    @Optional() @Inject(REQUEST)
    private readonly request: Request | null,
    @Inject(PLATFORM_ID) platformId: string,
  ) {
    this.httpClient = new HttpClient(httpBackend);
    this.isBrowser = isPlatformBrowser(platformId);
  }

  public async login(apiKey: string): Promise<boolean> {
    const success = await firstValueFrom(this.httpClient.get(`${environment.apiUrl}/admin/packages`, {
      headers: {
        Authorization: `Bearer ${apiKey}`,
      },
    }).pipe(
      map(() => true),
      catchError(() => of(false)),
    ));

    if (!success) {
      return false;
    }

    this.document.cookie = [
      `${this.cookieName}=${encodeURIComponent(apiKey)}`,
      "Path=/",
      `Max-Age=${8 * 60 * 60}`,
      "Secure",
      `SameSite=${environment.sameOrigin ? 'Lax' : 'None'}`,
    ].join("; ");

    return true;
  }

  public getApiKey(): string | null {
    const prefix = `${encodeURIComponent(this.cookieName)}=`;
    let cookies: string;
    if (this.isBrowser) {
      cookies = this.document.cookie;
    } else {
      cookies = this.request!.headers.get('cookie') ?? '';
    }

    const cookie = cookies
      .split('; ')
      .find(cookie => cookie.startsWith(prefix));

    return cookie
      ? decodeURIComponent(cookie.substring(prefix.length))
      : null;
  }

  public isLoggedIn(): boolean {
    return this.getApiKey() !== null;
  }
}
