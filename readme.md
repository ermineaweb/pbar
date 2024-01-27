## Pbar

Add a (very) simple progress bar or spinner to golang CLI. Does not support multilines.

```bash
go get github.com/ermineaweb/pbar
```

## Quick start

### Progress Bar

```golang
// progress bar with total tasks known
pb := pbar.NewPbar(pbar.ConfigPbar{TotalTasks: uint16(len(tasks))})

for _, task := range tasks {
    longWork(task)
    pb.Add(1)
}
```

#### Custom configuration

```golang
pbar := pbar.NewPbar(
    pbar.ConfigPbar{
        TotalTasks:           uint16(tasks),
        CharDone:             '-',
        CharTodo:             '-',
        ColorPercentWorking:  pbar.RED_BRIGHT,
        ColorPercentFinished: pbar.GREEN,
        ColorCharDone:        pbar.RED_BRIGHT,
        ColorCharTodo:        pbar.BLACK_BRIGHT,
    },
)
```

### Spinner

```golang
// spinner with total tasks unknown
sp := pbar.NewSpinner(pbar.ConfigSpinner{})

sp.Start()

// ... some work ...

sp.Stop()
```

#### Custom configuration

```golang
spinner := pbar.NewSpinner(
    pbar.ConfigSpinner{
        Spinner:          pbar.SPINNER_ARROW,
        StartMessage:     "Let's work!",
        StopMessage:      "Job's done!",
        ColorSpinner:     pbar.RED_BRIGHT,
        ColorTimer:       pbar.BLUE_BRIGHT,
        AnimationDelayMs: 130,
    },
)
```
