# Contributing

Contributions are always welcome, both reporting issues and submitting pull requests!

### Code of Conduct

This project adheres to the Contributor Covenant [Code of
Conduct](/CODE_OF_CONDUCT.md). By participating, you are expected to uphold
this code. Please report unacceptable behavior to
code-of-conduct@zeebe.io.

### Reporting issues

Please make sure to include any potentially useful information in the issue, so we can pinpoint the issue faster without going back and forth.

- What SHA of zbc-go are you running? If this is not the latest SHA on the master branch, please try if the problem persists with the latest version.
- You can set `zbc.Logger` to a [log.Logger](http://golang.org/pkg/log/#Logger) instance to capture debug output. Please include it in your issue description.
- Also look at the logs of the Zeebe broker you are connected to. If you see anything out of the ordinary, please include it.

Also, please include the following information about your environment, so we can help you faster:

- What version of Zeebe are you using?
- What version of Go are you using?
- What are the values of your Client/Broker configuration?


### Submitting pull requests

We will gladly accept bug fixes, or additions to this library. Please fork this library, commit & push your changes, and open a pull request. Because this library is in production use by many people and applications, we code review all additions. To make the review process go as smooth as possible, please consider the following.

- If you plan to work on something major, please open an issue to discuss the design first.
- Don't break backwards compatibility. If you really have to, open an issue to discuss this first.
- Make sure to use the `go fmt` command to format your code according to the standards. Even better, set up your editor to do this for you when saving.
- Run [go vet](https://godoc.org/golang.org/x/tools/cmd/vet) to detect any suspicious constructs in your code that could be bugs.
- You may also want to run [golint](https://github.com/golang/lint) as well to detect style problems.
- Add tests that cover the changes you made. Make sure to run `go test` with the `-race` argument to test for race conditions.
- Make sure your code is supported by all the Go versions we support. You can rely on [Travis CI](https://travis-ci.org/jsam/zbc-go) for testing older Go versions
- Make sure that you don't commit any of development dependencies such as resursive printers, debugers, etc
