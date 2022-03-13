import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Report} from "../model/report";
import {environment} from "../../environments/environment";
import {Transaction} from "../model/transaction";
import {Category} from "../model/category";

@Injectable({
  providedIn: 'root'
})
export class TransactionService {

  constructor(private http: HttpClient) { }

  public import(data: object): Observable<number[]> {
    return this.http.post<number[]>(`${this.getResourceUrl()}`, data);
  }

  public getTransactionsInMonth(month: number, year: number): Observable<Transaction[]> {
    return this.http.get<Transaction[]>(`${this.getResourceUrl()}/get-all/${month}/${year}`);
  }

  public forceSetCategory(transaction: Transaction, category: Category): Observable<Transaction> {
    return this.http.post<Transaction>(`${this.getResourceUrl()}/force-set-category`, {
      transactionId: transaction.id,
      categoryId: category ? category.id : null
    });
  }

  public remove(transaction: Transaction): Observable<void> {
    return this.http.delete<void>(`${this.getResourceUrl()}/${transaction.id}`);
  }

  private getResourceUrl(): string {
    return `${environment.apiUrl}/account-transaction`;
  }
}
