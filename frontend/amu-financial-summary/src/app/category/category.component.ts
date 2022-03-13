import { Component, OnInit } from '@angular/core';
import {CategoryService} from "../service/category.service";
import {Category} from "../model/category";

@Component({
  selector: 'app-category',
  templateUrl: './category.component.html',
  styleUrls: ['./category.component.scss']
})
export class CategoryComponent implements OnInit {

  public isLoading: boolean = false;
  public categories: Category[] = [];

  constructor(
    private categoryService: CategoryService
  ) { }

  ngOnInit(): void {
    this.loadCategories();
  }

  public addNew(): void {
    this.categories = [{name: '', regexps: []} as Category, ...this.categories];
  }

  public remove(category: Category): void {
    this.categories.splice(this.categories.findIndex(c => c.id === category.id), 1);
  }

  private loadCategories(): void {
    this.isLoading = true;
    this.categoryService.getAll().subscribe((res) => {
      this.isLoading = false;
      this.categories = res;
    }, (err) => {
      this.isLoading = false;
      console.error('Cannot load categories', err);
    });
  }
}
