package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	InstanceTypePostgres = "POSTGRES"
	InstanceTypeMySQL    = "MYSQL"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"instance_type": {
					Description:  "TODO",
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{InstanceTypePostgres, InstanceTypeMySQL}, true),
					DefaultFunc:  schema.EnvDefaultFunc("CLOUDSQL_INSTANCE_TYPE", InstanceTypePostgres),
				},
				"instance_name": {
					Description: "TODO",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("CLOUDSQL_INSTANCE_NAME", ""),
				},
				"database": {
					Description: "TODO",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("CLOUDSQL_DATABASE", ""),
				},
				"username": {
					Description: "TODO",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("CLOUDSQL_USERNAME", ""),
				},
				"password": {
					Description: "TODO",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("CLOUDSQL_PASSWORD", ""),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{},
			ResourcesMap: map[string]*schema.Resource{
				"cloudsql_migration": resourceCloudSQLMigration(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type providerConfig struct {
	InstanceType string
	InstanceName string
	Database     string
	Username     string
	Password     string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-cloudsql", version)
		// TODO: myClient.UserAgent = userAgent

		return &providerConfig{
			InstanceType: d.Get("instance_type").(string),
			InstanceName: d.Get("instance_name").(string),
			Database:     d.Get("database").(string),
			Username:     d.Get("username").(string),
			Password:     d.Get("password").(string),
		}, nil
	}
}

func leveledLog(level string) func(format string, v ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", strings.ToUpper(level))
	return func(format string, v ...interface{}) {
		log.Printf(prefix+format, v...)
	}
}

var traceLog = leveledLog("trace")
var debugLog = leveledLog("debug")
var infoLog = leveledLog("info")
var warnLog = leveledLog("warn")
var errorLog = leveledLog("error")
