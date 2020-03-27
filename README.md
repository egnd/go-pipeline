# wpool

Go package to organize pool of workers

## Quick start:
1. Import package:
```
import "gitlab.com/egnd/wpool"
```
2. Create and start pool:
```
pool := wpool.NewPool().Start()
```
3. Run workers and attach them to pool:
```
for i := 1; i <= 5; i++ {
    worker := wpool.NewWorker(fmt.Sprintf("worker-%d", i), pool)
    pool.RegisterWorker(worker)
}
```
4. Add jobs to pool:
```
for i := 1; i <= 100; i++ {
    pool.AddJob(wpool.NewJob(fmt.Sprintf("job-%d", i), func(job JobInterface) (err error) {
        // job code
        return
    }))
}
```
5. Waiting while workers executes all jobs
```
pool.Wait()
```

### Hints:
1. Decorating jobs:
```
// Prepend job logic. Job will be executed only if error is nil.
worker.PrependJob(func(job JobInterface) (err error) {
    log.Print("prepend " + job.Name() + " at " + worker.Name())
    return
})

// Append job logic.
worker.AppendJob(func(job JobInterface) (err error) {
    log.Print("append " + job.Name() + " at " + worker.Name())
    return
})
```
