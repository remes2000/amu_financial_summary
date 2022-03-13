import { Injectable } from '@angular/core';
import {Category} from "../model/category";
import {Observable} from "rxjs";
import {environment} from "../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {Report} from "../model/report";

@Injectable({
  providedIn: 'root'
})
export class ReportService {

  constructor(private http: HttpClient) { }

  public generate(month: number, year: number): Observable<Report> {
    return this.http.get<Report>(`${this.getResourceUrl()}/${month}/${year}`);
  }

  private getResourceUrl(): string {
    return `${environment.apiUrl}/report`;
  }
}
