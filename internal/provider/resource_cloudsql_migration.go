package provider

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

func resourceCloudSQLMigration() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",

		CreateContext: resourceCloudSQLMigrationCreate,
		ReadContext:   resourceCloudSQLMigrationRead,
		UpdateContext: resourceCloudSQLMigrationUpdate,
		DeleteContext: resourceCloudSQLMigrationDelete,

		Schema: map[string]*schema.Schema{
			"instance_type": {
				Description:  "TODO",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{InstanceTypePostgres, InstanceTypeMySQL}, true),
			},
			"instance_name": {
				Description: "TODO",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"database": {
				Description: "TODO",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"username": {
				Description: "TODO",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "TODO",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"up": {
				Description: "TODO",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"down": {
				Description: "TODO",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCloudSQLMigrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db, err := sqlOpen(ctx, d, meta)
	if err != nil {
		return diag.Errorf("Failed to open connect to database instance: %v", err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return diag.Errorf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	query := d.Get("up").(string)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return diag.Errorf("Failed to exec query `%s`, returned: %v", query, err)
	}

	if err := tx.Commit(); err != nil {
		return diag.Errorf("Failed to commit query: %v", err)
	}

	d.SetId(query)
	return nil
}

func resourceCloudSQLMigrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("Can not read resource without ID")
	}
	return nil
}

func resourceCloudSQLMigrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("Can not update resource without ID")
	}
	return nil // Also all relevant operations force new, so nothing to do here
}

func resourceCloudSQLMigrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db, err := sqlOpen(ctx, d, meta)
	if err != nil {
		return diag.Errorf("Failed to open connect to database instance: %v", err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return diag.Errorf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	query := d.Get("down").(string)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return diag.Errorf("Failed to exec query `%s`, returned: %v", query, err)
	}

	if err := tx.Commit(); err != nil {
		return diag.Errorf("Failed to commit query: %v", err)
	}

	return nil
}

func sqlOpen(ctx context.Context, d *schema.ResourceData, meta interface{}) (*sql.DB, error) {
	p := meta.(*providerConfig)
	instanceType := d.Get("instance_type").(string)
	if instanceType == "" {
		instanceType = p.InstanceType
	}
	instanceName := d.Get("instance_name").(string)
	if instanceName == "" {
		instanceName = p.InstanceName
	}
	database := d.Get("database").(string)
	if database == "" {
		database = p.Database
	}
	username := d.Get("username").(string)
	if username == "" {
		username = p.Username
	}
	password := d.Get("password").(string)
	if password == "" {
		password = p.Password
	}
	if strings.ToUpper(instanceType) == InstanceTypePostgres {
		return sql.Open("cloudsqlpostgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", instanceName, database, username, password))
	} else {
		cfg := mysql.Cfg(instanceName, username, password)
		cfg.DBName = database
		cfg.ParseTime = true
		const timeout = 60 * time.Second
		cfg.Timeout = timeout
		cfg.ReadTimeout = timeout
		cfg.WriteTimeout = timeout
		return mysql.DialCfg(cfg)
	}
}
