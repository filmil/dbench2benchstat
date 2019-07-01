# dbench2benchstat

A converter from the dart benchmark format into the go benchstat format.

This allows using benchstat for comparing benchmark results.


Conversion goes from text lines such as:

```
basic/increment -> avg 0.00927782162588819ms out of 2534 samples. (std dev 0.00161196365508829, min 0.008, max 0.082)
```

into:

```
basic/increment      2532         123456 ns/op               56.65 MB/s
```
