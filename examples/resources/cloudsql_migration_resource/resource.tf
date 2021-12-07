resource "cloudsql_migration" "example" {
  instance_type = "POSTGRES" # optional defaults to POSTGRES
  instance_name = "project:region:name"
  database      = "mydb"
  username      = "myuser"
  password      = "mypw"
  # All above attributes are optional and can also be set on provider level and are inherited

  up   = "GRANT ..."
  down = "..."
}
