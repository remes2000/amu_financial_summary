import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {CategoryComponent} from "./category/category.component";
import {SummaryComponent} from "./summary/summary.component";
import {ImportComponent} from "./import/import.component";

const routes: Routes = [
  {
    path: 'category',
    component: CategoryComponent
  },
  {
    path: '',
    component: SummaryComponent
  },
  {
    path: 'import',
    component: ImportComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
