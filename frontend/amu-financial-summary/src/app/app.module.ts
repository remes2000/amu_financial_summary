import {APP_INITIALIZER, NgModule} from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CategoryComponent } from './category/category.component';
import {HTTP_INTERCEPTORS, HttpClientModule} from "@angular/common/http";
import { LoaderComponent } from './loader/loader.component';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { CategoryCardComponent } from './category/category-card/category-card.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { SummaryComponent } from './summary/summary.component';
import { SummaryCardComponent } from './summary/summary-card/summary-card.component';
import { CurrencyPipe } from './pipe/currency.pipe';
import { CurrencyPipe as AngularCurrencyPipe } from '@angular/common';
import '@angular/common/locales/global/pl';
import {ChartsModule} from "ng2-charts";
import { ImportComponent } from './import/import.component';
import { TransactionsTableComponent } from './summary/transactions-table/transactions-table.component';
import * as crypto from 'crypto-js';
import {showPasswordPrompt} from "./password-prompt";
import {AuthenticationInterceptor} from "./interceptor/authentication.interceptor";
import {ErrorHandlerInterceptor} from "./interceptor/error-handler.interceptor";

const setCredentials = async () => {
  let apiKey = localStorage.getItem('API_KEY');
  if (!apiKey) {
    apiKey = crypto.SHA512(await showPasswordPrompt()).toString();
    localStorage.setItem('API_KEY', apiKey);
  }
};

@NgModule({
  declarations: [
    AppComponent,
    CategoryComponent,
    LoaderComponent,
    CategoryCardComponent,
    SummaryComponent,
    SummaryCardComponent,
    CurrencyPipe,
    ImportComponent,
    TransactionsTableComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    NgbModule,
    FormsModule,
    ReactiveFormsModule,
    ChartsModule
  ],
  providers: [
    CurrencyPipe,
    AngularCurrencyPipe,
    { provide: APP_INITIALIZER, useFactory: () => setCredentials, multi: true},
    { provide: HTTP_INTERCEPTORS, useClass: AuthenticationInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorHandlerInterceptor, multi: true }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
