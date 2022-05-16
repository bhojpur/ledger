# Bhojpur Ledger - Data Processing Engine

The `Bhojpur Ledger` is a double-entry financial accounting engine applied within
the [Bhojpur.NET Platform](https://github.com/bhojpur/platform/) ecosystem for
delivery of distributed `applications` or `services`.

## Simple Usage

You can run the Bhojpur Ledger `server` using the following command. It owns the
financial record keeping services.

```bash
ledgersvr
```

You can run the Bhojpur Ledger `client` using the following command. It invokes
the `server` using gRPC mechanism for data entry.

```bash
ledgerctl jsonjournal '{"Payee":"Yunica Retail","Date":"2022-05-15T00:00:00Z","AccountChanges":[{"Name":"Asset:Cash","Description":"Cash is better","Currency":"INR","Balance":"100"},{"Name":"Revenue:Sales","Description":"Income is good","Currency":"INR","Balance":"-100"}]}'
```

You can run the Bhojpur Ledger `reporter` using the following command. It invokes
the `server` using gRPC mechanism for report generation.

```bash
ledgerepo transactions
ledgerepo trialbalance
ledgerepo pdf -template profitandloss
```
