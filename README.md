tempfile
========

An alternative implementation of ioutil.TempFile that provides a little more control over the temporary files being created. It allows the caller to provide a file suffix and/or file mode to set the permissions to something other than the default 0600.


## Examples

```go
// Make a temp file in /var/tmp with the default file mode 0600. This is
// equivalent to what ioutil.TempFile does.
f, err := tempfile.New("/var/tmp", "myfile")
if err != nil {
  panic(err)
}
defer f.Close()
defer os.Remove(f.Name())

// Make a temporaty file with the given prefix and suffix
f, err = tempfile.NewSuffix("", "myfile", ".tmp")
...

// Create a world-readable temporary file in the OS' temp dir by providing
// a file mode.
f, err = tempfile.NewMode("", "myfile", 0644)
...

// Combining all of the above
f, err = tempfile.NewSuffixAndMode("", "myfile", ".ext", 0644)
...
```

## Links
- ioutil package - https://golang.org/pkg/io/ioutil/
- GoDoc for tempfile - https://godoc.org/github.com/folbricht/tempfile
