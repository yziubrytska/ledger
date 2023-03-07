# ledger

Ledger service executes transactions, returns balance and transactions history

## Dependencies

- Postgres DB (`accounts`, `transactions` tables)
- Log Level 
- Listen Address

##Tests 

- Unit tests are provided for service package
- E2E tests are provided, but they don't cover all the cases

##Additional notes

- Security measures are not implemented (no token parsing, no user - password check)