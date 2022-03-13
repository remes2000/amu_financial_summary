import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Category} from "../model/category";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class CategoryService {
  constructor(private http: HttpClient) { }

  public getAll(): Observable<Category[]> {
    return this.http.get<Category[]>(`${this.getResourceUrl()}`);
  }

  public save(category: Category): Observable<Category> {
    if(category.id) {
      return this.http.put<Category>(`${this.getResourceUrl()}`, category);
    }
    return this.http.post<Category>(`${this.getResourceUrl()}`, category);
  }

  public remove(category: Category): Observable<void> {
    return this.http.delete<void>(`${this.getResourceUrl()}/${category.id}`);
  }

  private getResourceUrl(): string {
    return `${environment.apiUrl}/category`;
  }
}
