import { Component, OnInit } from '@angular/core';
import {TransactionService} from "../service/transaction.service";
import {parse} from "@angular/compiler/src/render3/view/style_parser";

@Component({
  selector: 'app-import',
  templateUrl: './import.component.html',
  styleUrls: ['./import.component.scss']
})
export class ImportComponent implements OnInit {

  public isLoading: boolean = false;
  public json: string = undefined;
  public message: string = '';

  constructor(
    private transactionService: TransactionService
  ) {}

  ngOnInit(): void {
  }

  public performImport(): void {
    let parsedJson: object = {};
    try {
      parsedJson = JSON.parse(this.json);
    } catch (err) {
      console.error('Invalid json', err);
      this.message = 'Invalid json';
      return;
    }
    this.isLoading = true;
    this.transactionService.import(parsedJson).subscribe((res) => {
      this.message = `Import success`;
      this.isLoading = false;
    }, (err) => {
      console.error('cannot perform import', err);
      this.message = 'Import failed';
      this.isLoading = false;
    });
  }

  public async generateJson(event: any): Promise<void> {
    const input = event.target;
    const content = await this.getFilesContent(input.files);
    this.json = this.santander2Json(content);
  }

  public getFilesContent(files: any): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      let loadedFiles = 0;
      let content = '';
      for (const file of files) {
        const reader = new FileReader();
        reader.onload = () => {
          const result = reader.result as string;
          const parts = result.split('\n');
          parts.shift();
          parts.pop();
          content += parts.join('\n');
          loadedFiles++;
          if (loadedFiles === files.length) {
            resolve(content);
          }
        };
        reader.readAsText(file);
      }
    });
  }

  public santander2Json(content: string): string {
    let result = '[\n';
    let i = 0;
    for (const line of content.split('\n')) {
      const values = line.split('\t', 8);
      if (values.length !== 8) {
        throw new Error(`Line should contain at least 8 values, got ${values.length} ` + line);
      }
      result += '\t' + (i==0?'':',') + this.santanderLine2Json(values) + '\n';
      i++;
    }
    return result + ']';
  }

  private santanderLine2Json(values: string[]): string {
    const object = {
      date: values[1],
      title: `${values[2]} ${values[3]} ${values[4]}`,
      amount: values[5].replace(/,/g, '.')
    };
    return JSON.stringify(object);
  }
}
