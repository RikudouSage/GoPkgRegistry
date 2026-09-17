import { Routes } from '@angular/router';
import {LoggedInGuard} from './guards/logged-in-guard';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./pages/package-list/package-list').then(m => m.PackageList),
    canActivate: [LoggedInGuard],
  },
  {
    path: 'package/create',
    loadComponent: () => import('./pages/package-detail/package-detail').then(m => m.PackageDetail),
    canActivate: [LoggedInGuard],
  },
  {
    path: 'package/:packageId',
    loadComponent: () => import('./pages/package-detail/package-detail').then(m => m.PackageDetail),
    canActivate: [LoggedInGuard],
  },
  {
    path: 'login',
    loadComponent: () => import('./pages/login-page/login-page').then(m => m.LoginPage),
  },
];
