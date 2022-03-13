import {ReportDetail} from "./report-details";

export interface Report {
  year?: number;
  month?: number;
  total?: string;
  totalIncome?: string;
  totalOutcome?: string;
  details?: ReportDetail[];
}
