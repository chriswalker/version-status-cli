# version-status-cli
Displays version age of staging vs production contexts in Kubernetes

## TODO
* Fix text/tabwriter - fatih/color issues (latter's colour escape codes confuse tabwriter a bit)
* Pass environments in on command line
  * basically working, but needs a little polish (not happy with names)
* General tidy up
* ~~Spinner UI to indicate working; already libs to do this (https://github.com/briandowns/spinner)~~
* Comments
* ~~Handler errors from goroutines~~
* Finish README!
* ~~diffsonly flag, that shows only services with a version difference (cut down on output)~~
