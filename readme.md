Add a simple progress bar or spinner to golang CLI.

```bash
go get github.com/ermineaweb/pbar
```

```golang
// progress bar with total tasks known
pb := pbar.NewPbar(50)
pb.Add(1)

// spinner with total tasks unknown
sp := pbar.NewSpinner(pbar.SPINNER_POINTS)
sp.Start()
defer sp.Stop()
```