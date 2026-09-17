import {Component, DOCUMENT, effect, Inject, PLATFORM_ID, signal} from '@angular/core';
import {RouterLink, RouterOutlet} from '@angular/router';
import {isPlatformBrowser} from '@angular/common';

@Component({
  imports: [RouterLink, RouterOutlet],
  selector: 'app-root',
  styleUrl: './app.scss',
  templateUrl: './app.html',
})
export class App {
  protected readonly theme = signal<'light' | 'dark'>('light');

  public constructor(
    @Inject(DOCUMENT) document: Document,
    @Inject(PLATFORM_ID) platformId: object,
  ) {
    const isBrowser = isPlatformBrowser(platformId);
    if (isBrowser) {
      const savedTheme = localStorage.getItem('theme');
      this.theme.set(
        savedTheme === 'dark' || savedTheme === 'light'
          ? savedTheme
          : window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light',
      );
    }

    effect(() => {
      if (!isBrowser) {
        return;
      }

      const theme = this.theme();
      document.documentElement.dataset['theme'] = theme;
      localStorage.setItem('theme', theme);
    });
  }

  protected toggleTheme(): void {
    this.theme.update(theme => theme === 'dark' ? 'light' : 'dark');
  }
}
