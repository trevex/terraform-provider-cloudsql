# Terraform Provider CloudSQL

```
This provider is a test with the main purpose to see how to improve automation of CloudSQL users with restricted permissions, e.g.:

resource "cloudsql_migration" "example" {
  instance_name = google_sql_database_instance.example.connection_name
  database      = google_sql_database.example.name
  username      = google_sql_user.admin.name
  password      = google_sql_user.admin.password

  up = <<-EOF
    CREATE USER readonly PASSWORD 'test';
    GRANT CONNECT ON DATABASE mydatabase TO readonly;
    GRANT USAGE ON SCHEMA "public" TO readonly;
    GRANT SELECT ON ALL TABLES IN SCHEMA "public" TO readonly;
    ALTER DEFAULT PRIVILEGES IN SCHEMA "public" GRANT SELECT ON TABLES TO readonly;
  EOF
  down = <<-EOF
    ALTER DEFAULT PRIVILEGES IN SCHEMA "public" REVOKE SELECT ON TABLES FROM readonly;
    REVOKE SELECT ON ALL TABLES IN SCHEMA "public" FROM readonly;
    REVOKE USAGE ON SCHEMA "public" FROM readonly;
    REVOKE CONNECT ON DATABASE mydatabase FROM readonly;
    DROP USER readonly;
  EOF
}
```


This repository is a *template* for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

 - A resource, and a data source (`internal/provider/`),
 - Examples (`examples/`) and generated documentation (`docs/`),
 - Miscellaneous meta files.
 
These files contain boilerplate code that you will need to edit to create your own Terraform provider. A full guide to creating Terraform providers can be found at [Writing Custom Providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://www.terraform.io/docs/registry/providers/publishing.html) so that others can use it.


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
