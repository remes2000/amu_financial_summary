import {Regexp} from "./regexp";

export interface Category {
  id?: number;
  name?: string;
  regexps?: Regexp[];
}
