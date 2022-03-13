import { Component } from '@angular/core';
import {logout} from "./interceptor/error-handler.interceptor";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'amu-financial-summary';

  public logout = logout;
}
