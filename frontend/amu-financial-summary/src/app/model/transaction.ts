import {Category} from "./category";

export interface Transaction {
  id?: number;
  title?: number;
  date?: number;
  amount?: number;
  category?: Category;
}
