import {Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Transaction} from "../../model/transaction";
import {Category} from "../../model/category";
import {TransactionService} from "../../service/transaction.service";

@Component({
  selector: 'app-transactions-table',
  templateUrl: './transactions-table.component.html',
  styleUrls: ['./transactions-table.component.scss']
})
export class TransactionsTableComponent implements OnChanges {

  @Input() transactions: Transaction[] = [];
  @Input() categories: Category[] = [];
  public message: any = {};

  constructor(
    private transactionService: TransactionService
  ) {}

  ngOnChanges(changes: SimpleChanges) {
    if(changes.transactions || changes.categories) {
      this.transactions
        .filter(t => t.category)
        .forEach(t => t.category = this.categories.find(c => c.id === t.category.id));
    }
  }

  public save(transaction: Transaction): void {
    this.message[transaction.id] = 'saving...';
    this.transactionService.forceSetCategory(transaction, transaction.category).subscribe((res) => {
      this.message[transaction.id] = 'save succeed';
    }, (err) => {
      this.message[transaction.id] = 'save failed';
    });
  }

  public remove(transaction: Transaction): void {
    if(!confirm('Are you sure you want to remove this transaction?')) {
      return;
    }
    this.message[transaction.id] = 'removing...';
    this.transactionService.remove(transaction).subscribe(() => {
      this.message[transaction.id] = 'removed';
    }, (err) => {
      this.message[transaction.id] = 'removing failed';
    });
  }
}
