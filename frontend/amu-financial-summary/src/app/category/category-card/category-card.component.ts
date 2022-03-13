import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Category} from "../../model/category";
import {FormArray, FormControl, FormGroup, Validators} from "@angular/forms";
import {Regexp} from "../../model/regexp";
import {CategoryService} from "../../service/category.service";

@Component({
  selector: 'app-category-card',
  templateUrl: './category-card.component.html',
  styleUrls: ['./category-card.component.scss']
})
export class CategoryCardComponent implements OnInit {

  @Input() public category: Category;

  @Output() public removeCategory: EventEmitter<Category> = new EventEmitter();

  public form: FormGroup;
  public regexps: FormArray;
  public message: string = '';

  constructor(
    private categoryService: CategoryService
  ) {}

  ngOnInit(): void {
    this.initForm();
  }

  public save(): void {
    if (this.form.invalid) {
      this.message = 'form is not valid';
      return;
    }
    this.message = 'saving';
    this.categoryService.save(this.getCategoryFromForm()).subscribe((res) => {
      this.category = res;
      this.initForm();
      this.message = 'save success';
    }, (err) => {
      this.message = 'save failed';
      console.error(err);
    });
  }

  public remove(): void {
    if (!confirm('Are you sure you want to remove this category?')) {
      return;
    }
    if (!this.category.id) {
      this.removeCategory.emit(this.category);
      return;
    }
    this.message = 'removing';
    this.categoryService.remove(this.category).subscribe(() => {
      this.removeCategory.emit(this.category);
    }, (err) => {
      this.message = 'remove failed';
      console.error(err);
    });
  }

  public addNewRegexp(): void {
    this.regexps.push(this.getRegexpFormGroup());
  }

  private getCategoryFromForm(): Category {
    return {
      id: this.category.id,
      name: this.form.get('categoryName').value,
      regexps: this.regexps.controls.map(regexpForm => {
        return {
          id: regexpForm.get('id').value,
          content: regexpForm.get('content').value
        } as Regexp;
      })
    } as Category;
  }

  private initForm(): void {
    this.form = new FormGroup({
      categoryName: new FormControl(this.category.name, [Validators.required]),
      regexps: new FormArray(this.category.regexps.map(regexp => this.getRegexpFormGroup(regexp)))
    });
    this.regexps = this.form.get('regexps') as FormArray;
  }

  private getRegexpFormGroup(regexp: Regexp = {content: ''} as Regexp): FormGroup {
    return new FormGroup({
      id: new FormControl(regexp.id),
      content: new FormControl(regexp.content, [Validators.required])
    });
  }
}
