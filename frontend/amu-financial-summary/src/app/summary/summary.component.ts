import { Component, OnInit } from '@angular/core';
import {ReportService} from "../service/report.service";
import {Report} from "../model/report";
import {TransactionService} from "../service/transaction.service";
import {forkJoin} from "rxjs";
import {Transaction} from "../model/transaction";
import {CategoryService} from "../service/category.service";
import {Category} from "../model/category";

const LAST_VISITED_REPORT_YEAR = 'LAST_VISITED_REPORT_YEAR';
const LAST_VISITED_REPORT_MONTH = 'LAST_VISITED_REPORT_MONTH';

@Component({
  selector: 'app-summary',
  templateUrl: './summary.component.html',
  styleUrls: ['./summary.component.scss']
})
export class SummaryComponent implements OnInit {

  public year: number = new Date().getFullYear();
  public month: number = new Date().getMonth() + 1;

  public message: string = '';
  public report: Report;
  public transactions: Transaction[];
  public categories: Category[];

  constructor(
    private reportService: ReportService,
    private transactionService: TransactionService,
    private categoryService: CategoryService
  ) { }

  ngOnInit(): void {
    this.setLastVisitedReport();
    this.generateReport();
  }

  public generateReport(): void {
    localStorage.setItem(LAST_VISITED_REPORT_YEAR, `${this.year}`);
    localStorage.setItem(LAST_VISITED_REPORT_MONTH, `${this.month}`);
    forkJoin([
      this.reportService.generate(this.month, this.year),
      this.transactionService.getTransactionsInMonth(this.month, this.year),
      this.categoryService.getAll()
    ]).subscribe((res) => {
      this.report = res[0];
      this.transactions = res[1];
      this.categories = res[2];
    }, (err) => {
      console.error('Cannot generate report', err);
      this.message = 'Generate report failed';
    });
  }

  private setLastVisitedReport(): void {
    if(!localStorage.getItem(LAST_VISITED_REPORT_YEAR) || !localStorage.getItem(LAST_VISITED_REPORT_MONTH)) {
      return;
    }
    this.year = Number(localStorage.getItem(LAST_VISITED_REPORT_YEAR));
    this.month = Number(localStorage.getItem(LAST_VISITED_REPORT_MONTH));
  }
}
