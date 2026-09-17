import {Component, DestroyRef, OnInit, signal} from '@angular/core';
import {Title} from '@angular/platform-browser';
import {Translator} from '../../services/translator';
import {form, FormField, required} from '@angular/forms/signals';
import {FormsModule} from '@angular/forms';
import {TranslocoPipe} from '@jsverse/transloco';
import {ToastrService} from 'ngx-toastr';
import {firstValueFrom} from 'rxjs';
import {Auth} from '../../services/auth';
import {Router} from '@angular/router';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';

interface LoginModel {
  apiKey: string;
}

@Component({
  imports: [
    FormField,
    FormsModule,
    TranslocoPipe
  ],
  selector: 'app-login-page',
  styleUrl: './login-page.scss',
  templateUrl: './login-page.html',
})
export class LoginPage implements OnInit {
  private readonly formData = signal<LoginModel>({apiKey: ''});
  protected readonly form = form(this.formData, schemaPath => {
    required(schemaPath.apiKey);
  });

  public constructor(
    private readonly title: Title,
    private readonly translator: Translator,
    private readonly toastr: ToastrService,
    private readonly auth: Auth,
    private readonly router: Router,
    private readonly destroyRef: DestroyRef,
  ) {
  }

  public ngOnInit(): void {
    this.translator.get('login.title').pipe(
      takeUntilDestroyed(this.destroyRef),
    ).subscribe(title => {
      this.title.setTitle(title);
    });
  }

  protected async onSubmit(): Promise<void> {
    if (!this.form().valid()) {
      this.toastr.error(
        await firstValueFrom(this.translator.get('login.error.empty_api_key')),
        await firstValueFrom(this.translator.get('common.notification.error')),
      );
      return;
    }

    if (await this.auth.login(this.formData().apiKey)) {
      await this.router.navigateByUrl('/');
      return;
    }

    this.toastr.error(
      await firstValueFrom(this.translator.get('login.error.invalid_api_key')),
      await firstValueFrom(this.translator.get('common.notification.error')),
    );
  }
}
