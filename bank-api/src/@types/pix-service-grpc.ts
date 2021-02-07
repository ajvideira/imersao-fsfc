import { Observable } from 'rxjs';

interface Bank {
  id: string;
  code: string;
  name: string;
}

interface Account {
  id: string;
  number: string;
  ownerName: string;
  createdAt: string;
  bank: Bank;
}

interface PixKeyFindResult {
  id: string;
  kind: string;
  key: string;
  createdAt: string;
  account: Account;
}

export interface PixService {
  registerPixKey: (data: {
    kind: string;
    key: string;
    accountID: string;
  }) => Observable<{ id: string; status: string; error: string }>;
  find: (data: { kind: string; key: string }) => Observable<PixKeyFindResult>;
}
