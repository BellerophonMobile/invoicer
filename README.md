Invoicer
========

Generate invoices from Toggl data.

Quick Start
-----------

* Run `make` to build (requires [Go](https://golang.org/)).
* Log in to [Toggl](https://toggl.com/), and view your profile settings.  Copy
  your "API token" at the bottom to a new file named `token.txt` in the root of
  this repository.
* Run `./bin/invoicer` without arguments to get a list of workspace and client
  IDs.
* Re-run `./bin/invoicer -WorkspaceID=XX -ClientID=XX`, with the desired
  workspace and client IDs.
* You'll now have a `invoice.tex` file of last month's data to tweak.

Configuration
-------------

`invoicer` supports the following options:

* `-WorkspaceID` - **REQUIRED** The workspace ID to report on.  Run `invoicer`
                   without arguments to list available workspaces.
* `-ClientID` - **REQUIRED** The client ID to report on.  Run `invoicer` without
                arguments to list available workspaces.
* `-TokenFile` - Override the default file to read for your Toggl API token.
* `-Template` - Specify the input template file.
* `-Output` - Specify the output invoice file.
* `-Month` - The month to report on.  Format is "YYYY-MM".  Defaults to the
             previous month.
* `-Rate` - The rate to bill at.

Templates
---------

The invoice template can be customized with the `-Template` flag.  Templates are
[Go templates](https://golang.org/pkg/text/template) with the following
functions available:

* `time` - Reports a time as "YYYY-MM-DD"
* `texDuration` - Converts a duration to some number of hours.
* `texEscape` - Escapes "&" symbols for use in LaTeX tables.
* `texCash` - Formats a float as a dollar amount with an escaped $.
