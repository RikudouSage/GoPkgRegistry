import {Component, DestroyRef, OnInit, Signal, WritableSignal} from '@angular/core';
import {Api} from '../../services/api';
import {Package} from '../../dto/package';
import {takeUntilDestroyed, toSignal} from '@angular/core/rxjs-interop';
import {Title} from '@angular/platform-browser';
import {Translator} from '../../services/translator';
import {TranslocoPipe} from '@jsverse/transloco';
import {RouterLink} from '@angular/router';

@Component({
  imports: [
    TranslocoPipe,
    RouterLink
  ],
  selector: 'app-package-list',
  styleUrl: './package-list.scss',
  templateUrl: './package-list.html',
})
export class PackageList implements OnInit {
  protected readonly packages: Signal<Package[] | null>;

  public constructor(
    api: Api,
    private readonly title: Title,
    private readonly translator: Translator,
    private readonly destroyRef: DestroyRef,
  ) {
    this.packages = toSignal(api.getAll(), {
      initialValue: null,
    });
  }

  public ngOnInit(): void {
    this.translator.get('packages.title').pipe(
      takeUntilDestroyed(this.destroyRef),
    ).subscribe(title => this.title.setTitle(title));
  }
}
