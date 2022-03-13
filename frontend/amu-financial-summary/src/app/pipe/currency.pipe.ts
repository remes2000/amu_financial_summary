import { Pipe, PipeTransform } from '@angular/core';
import { CurrencyPipe as AngularCurrencyPipe } from '@angular/common';

@Pipe({
  name: 'currency'
})
export class CurrencyPipe implements PipeTransform {
  constructor(
    private currencyPipe: AngularCurrencyPipe
  ) {}

  transform(value: number|string): unknown {
    return this.currencyPipe.transform(value, 'PLN', 'symbol-narrow', '', 'pl');
  }

}
