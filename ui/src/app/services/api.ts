import {Injectable, Service} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Package} from '../dto/package';
import {apiUrl} from '../http/api-url';

@Injectable({
  providedIn: 'root',
})
export class Api {
  public constructor(
    private readonly httpClient: HttpClient,
  ) {
  }

  public getAll(): Observable<Package[]> {
    return this.httpClient.get<Package[]>(apiUrl('/admin/packages'));
  }

  public findById(id: number): Observable<Package> {
    return this.httpClient.get<Package>(apiUrl(`/admin/packages/by-id/${id}`));
  }

  public createOrUpdate(pkg: Package): Observable<Package> {
    return this.httpClient.post<Package>(apiUrl('/admin/packages'), pkg);
  }
}
