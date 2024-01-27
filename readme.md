### Pbar

Add a (very) simple progress bar or spinner to golang CLI. Does not support multilines.

```bash
go get github.com/ermineaweb/pbar
```

### Quick start

#### Progress Bar

```golang
// progress bar with total tasks known
pb := pbar.NewPbar(pbar.ConfigPbar{TotalTasks: uint16(len(tasks))})

for _, task := range tasks {
    longWork(task)
    pb.Add(1)
}
```

#### Spinner

```golang
// spinner with total tasks unknown
sp := pbar.NewSpinner(pbar.ConfigSpinner{})

sp.Start()

// ... some work ...

sp.Stop()
```

### Custom configurations

```golang
pb := pbar.NewPbar(
    pbar.ConfigPbar{
        TotalTasks: uint16(tasks),
        CharDone:   '#',
        CharTodo:   '-',
    },
)
```

```golang
sp := pbar.NewSpinner(
    pbar.ConfigSpinner{
        Spinner:      pbar.SPINNER_ARROW,
        StartMessage: "Let's work!",
        StopMessage:  "Job's done!",
    },
)
```
