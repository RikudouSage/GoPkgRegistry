import {isPlatformBrowser} from '@angular/common';
import {Component, computed, DestroyRef, effect, inject, OnInit, PLATFORM_ID, signal} from '@angular/core';
import {emptyPackage, Package, Vcs} from '../../dto/package';
import {form, FormField, required} from '@angular/forms/signals';
import {Api} from '../../services/api';
import {ActivatedRoute, Router} from '@angular/router';
import {Title} from '@angular/platform-browser';
import {Translator} from '../../services/translator';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {FormsModule} from '@angular/forms';
import {TranslocoPipe} from '@jsverse/transloco';
import {firstValueFrom} from 'rxjs';
import {ToastrService} from 'ngx-toastr';

const MAIN_BRANCH_STORAGE_KEY = 'package-main-branch';

@Component({
  imports: [
    FormsModule,
    TranslocoPipe,
    FormField
  ],
  selector: 'app-package-detail',
  styleUrl: './package-detail.scss',
  templateUrl: './package-detail.html',
})
export class PackageDetail implements OnInit {
  protected readonly Vcs = Vcs;

  protected readonly formData = signal<Package<''>>(emptyPackage(''));
  protected readonly loading = signal(true);
  protected readonly form = form(this.formData, schemaPath => {
    required(schemaPath.import_path);
    required(schemaPath.vcs);
    required(schemaPath.repository_url);
  });
  protected readonly isNew = computed(() => this.formData().id === 0);
  protected readonly isGithub = computed(() => this.formData().repository_url.startsWith('https://github.com/'));
  private readonly platformId = inject(PLATFORM_ID);
  protected readonly mainBranch = signal(
    isPlatformBrowser(this.platformId) ? localStorage.getItem(MAIN_BRANCH_STORAGE_KEY) ?? 'master' : 'master',
  );

  public constructor(
    private readonly api: Api,
    private readonly activatedRoute: ActivatedRoute,
    private readonly title: Title,
    private readonly translator: Translator,
    private readonly destroyRef: DestroyRef,
    private readonly toastr: ToastrService,
    private readonly router: Router,
  ) {
    if (isPlatformBrowser(this.platformId)) {
      effect(() => localStorage.setItem(MAIN_BRANCH_STORAGE_KEY, this.mainBranch()));
    }

    effect(() => {
      if (!this.isGithub()) {
        return;
      }

      const repositoryUrl = this.formData().repository_url.replace(/\/+$/, '');
      const mainBranch = this.mainBranch();
      if (!this.form.source_url().touched()) {
        this.form.source_url().value.set(repositoryUrl);
      }
      if (!this.form.source_dir_url().touched()) {
        this.form.source_dir_url().value.set(`${repositoryUrl}/tree/${mainBranch}{/dir}`);
      }
      if (!this.form.source_file_url().touched()) {
        this.form.source_file_url().value.set(`${repositoryUrl}/blob/${mainBranch}{/dir}/{file}#L{line}`);
      }
    });
  }

  public ngOnInit(): void {
    this.activatedRoute.params.subscribe(params => {
      const id = params?.['packageId'] as string | undefined;
      if (!id) {
        this.translator.get('package.title.create').pipe(
          takeUntilDestroyed(this.destroyRef),
        ).subscribe(translated => this.title.setTitle(translated));

        this.loading.set(false);
        return;
      }

      this.api.findById(Number(id)).subscribe(pkg => {
        this.formData.set(this.packageToModel(pkg));
        this.loading.set(false);

        this.translator.get('package.title.edit', {pkgName: pkg.import_path}).pipe(
          takeUntilDestroyed(this.destroyRef),
        ).subscribe(
          translated => this.title.setTitle(translated),
        )
      });
    });
  }

  protected async onSubmit(): Promise<void> {
    if (!this.form().valid()) {
      this.toastr.error(
        await firstValueFrom(this.translator.get('package.error.invalid_form')),
        await firstValueFrom(this.translator.get('common.notification.error')),
      );
      return;
    }

    const pkg = this.modelToPackage(this.formData());
    try {
      const updated = await firstValueFrom(this.api.createOrUpdate(pkg));
      if (this.isNew()) {
        await this.router.navigateByUrl(`package/${updated.id}`);
        this.toastr.success(
          await firstValueFrom(this.translator.get('package.success.created')),
          await firstValueFrom(this.translator.get('common.notification.success')),
        )
        return;
      }

      this.formData.set(this.packageToModel(updated));
      this.toastr.success(
        await firstValueFrom(this.translator.get('package.success.updated')),
        await firstValueFrom(this.translator.get('common.notification.success')),
      )
    } catch (e) {
      this.toastr.error(
        await firstValueFrom(this.translator.get('package.error.failed_update')),
        await firstValueFrom(this.translator.get('common.notification.error')),
      );
      return;
    }
  }

  private packageToModel(pkg: Package): Package<''> {
    return {
      ...pkg,
      source_url: pkg.source_url ?? '',
      source_dir_url: pkg.source_dir_url ?? '',
      source_file_url: pkg.source_file_url ?? '',
    }
  }

  private modelToPackage(pkg: Package<''>): Package {
    return {
      ...pkg,
      source_url: pkg.source_url || null,
      source_dir_url: pkg.source_dir_url || null,
      source_file_url: pkg.source_file_url || null,
    }
  }
}
