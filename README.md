# **indi**: simple lazy initializer 

Parallelize initialization of your dependencies with ease.

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

Every step requires some time. You want to spend as little time as possible.

Now, it's not hard to see that the way initialization process goes is kind of the best already. You can't start
with initializing `c` because it depends on `a` and `b`. And `b` depends on `a`. And you can't parallelize anything.

But as your code base grows, you have more and more such dependencies. And it becomes harder to keep track of them and
initialize them the optimal way. You try to initialize several things in parallel.

```go
func main() {
	var (
        a *A
        b *B
        ...
        z *Z
    )

    aReady := make(chan struct{})
    var wg sync.WaitGroup
    wg.Add(3)
	
    go func() {
        defer wg.Done()
        a = a.New()
        close(aReady)
        b = b.New(a)
        c = c.New(a, b)
    }()

    go func() {
        defer wg.Done()
        <-aReady
        d = d.New(a)
        e = e.New(d)
    }()
    
    go func() {
        ...
    }
}
```

You're not sure anymore that you are doing it the optimal way.

So you need a sort of dependency graph...

## Here comes the **indi**

```go
// Dependency tree:
// A -> B -> C
// A -> D

func main() {
    var (
        a A
        b B
        c C
        d D
    )

    indi.Declare(&a, func () (*A, error) { return NewA(&b, &d) }, &b, &d)
    indi.Declare(&b, func () (*B, error) { return NewB(&c) }, &c)
    indi.Declare(&c, NewC)
    indi.Declare(&d, NewD)

    err := indi.Init()
    ...
}
```

See `example/main.go` to get the idea.

You casually declare all the stuff regardless of the order, and **indi** takes care of it.

P.S. Generics are used instead of `interface{}` on purpose:
1. To avoid type assertions (well, there are still some under the hood).
2. To force you to place all the dirty initialization with actual types in `main` package, or `init` package, or
wherever but far from the actual code that uses those dependencies. `main` is usually the place where you pass
implementations as interface functions parameters.
3. To support the principle "accept interfaces, return actual types".
4. Just for fun!

Enjoy!
