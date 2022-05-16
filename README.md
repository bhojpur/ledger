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
ledgerctl --help
```

You can run the Bhojpur Ledger `reporter` using the following command. It invokes
the `server` using gRPC mechanism for report generation.

```bash
ledgerepo --help
```
