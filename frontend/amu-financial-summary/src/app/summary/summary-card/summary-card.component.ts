import {Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Report} from "../../model/report";
import {ChartDataSets, ChartOptions} from "chart.js";

@Component({
  selector: 'app-summary-card',
  templateUrl: './summary-card.component.html',
  styleUrls: ['./summary-card.component.scss']
})
export class SummaryCardComponent implements OnChanges {

  @Input() public report: Report;
  public pieChartData: number[];
  public barChartData: ChartDataSets[];
  public chartLabels: string[];
  public pieChartOptions: ChartOptions = {
    legend: {
      position: 'right'
    }
  };

  constructor() { }

  ngOnChanges(changes: SimpleChanges) {
    if (changes.report) {
      this.setPieChart();
    }
  }

  private setPieChart(): void {
    this.chartLabels = [];
    this.pieChartData = [];
    this.report.details
      .filter(d => Number(d.amount) < 0)
      .forEach(d => {
       this.chartLabels.push(d.category);
       this.pieChartData.push(-Number(d.amount));
    });
    this.barChartData = [{
      data: [...this.pieChartData],
      label: ''
    }];
    console.log(this.barChartData);
  }
}
