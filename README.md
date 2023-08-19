# **indi** is Not a DI

## **indi** is a simple as hell lazy initializer

The real beauty of this tiny module is that it doesn't let you stray too far from the Go style into the dark world of
Java.

The problem it solves is only initialization. It doesn't do anything else.

## The problem

Imagine you have a giant application that access same stateful resources in different places. You write something like
this (go style!):

```go
func main() {
    a := a.New()
    b := b.New(a)
    c := c.New(a, b)
}
```

Now, it's not hard to see that the way initialization process goes is kind of the best already. You can't start
with initializing `c` because it depends on `a` and `b`. And `b` depends on `a`. And you can't parallelize anything.

But as your code base grows, you have more and more such dependencies. And it becomes harder to keep track of them and
initialize them the optimal way. You try to initialize several things in parallel.

```go
func main() {
    a := a.New()
    b := b.New(a)
    c := c.New(a, b)
    go initCache(a, b, c)
    d := d.New(a, c)
    go initDB(d)
    e := e.New(b, d)
    go initRPCConns(e)
    f := f.New(b, c, d)
    // ...
    z := z.New(d, x)
}

func initCache(...) {
    // ...
}

func initDB(...) {
    // ...
}

func initRPCConns(...) {
    // ...
}
```

You're not sure anymore that you do it in the best order. So you need a sort of dependency graph.

## Here comes the **indi**

```go
// Dependency tree:
// A -> B -> C
// A -> D

func main() {
    indi.Set("A", func(r *indi.Registry) (a *services.A, err error) {
        var (
            b *services.B
            d *services.D
        )
        if b, err = indi.GetFromRegistry[*services.B](r, "B"); err != nil {
            return nil, err
        }
        if d, err = indi.GetFromRegistry[*services.D](r, "D"); err != nil {
            return nil, err
        }

        return services.NewServiceA(b, d) // 10 seconds
    })
    indi.Set("B", func(r *indi.Registry) (b *services.B, err error) {
        var c *services.C
        if c, err = indi.GetFromRegistry[*services.C](r, "C"); err != nil {
            return nil, err
        }

        return services.NewServiceB(c) // 15 seconds
    })
    indi.Set("C", func(r *indi.Registry) (*services.C, error) {
        return services.NewServiceC() // 5 seconds
    })
    indi.Set("D", func(r *indi.Registry) (*services.D, error) {
        return services.NewServiceD() // 10 seconds
    })

    now := time.Now()
    if err := indi.Init(); err != nil {
        panic(err)
    }

    println(time.Since(now).Seconds()) // should be 30 seconds, not 40
}
```

You casually declare all the stuff regardless of the order, and **indi** takes care of it.

P.S. Generics are used instead of `interface{}` on purpose:
1. To avoid type assertions (well, there are still some under the hood).
2. To force you to place all the dirty initialization with implemented types in `main` package, or `init` package, or
wherever but far from the actual code that use those dependencies.
3. To support the principle "accept interfaces, pass actual types".
4. To make it clear what type is used to initialize what.
5. Just for fun!

Enjoy!

## TODO:

- [X] Add tests
- [X] Add default global registry
- [ ] Support parallelism limitation (blocking variation)
- [ ] Support cleanup
