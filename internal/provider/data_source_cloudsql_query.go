package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCloudSQLQuery() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",
		ReadContext: dataSourceCloudSQLQueryRead,

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
			"query": {
				Description: "TODO",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceCloudSQLQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	query := d.Get("query").(string)
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return diag.Errorf("Failed to exec query `%s`, returned: %v", query, err)
	}
	defer rows.Close()

	data := []map[string]string{}
	resultSet := 0
	for {
		cols, err := rows.Columns()
		if err != nil {
			return diag.Errorf("Failed to retrieve columns of resultSet#%d: %s", resultSet, err)
		}

		data[resultSet] = map[string]string{}
		for rows.Next() {
			ptrs := make([]interface{}, len(cols))
			vals := make([]string, len(cols))

			for i := range ptrs {
				ptrs[i] = &vals[i]
			}

			if err := rows.Scan(ptrs...); err != nil {
				return diag.Errorf("Failed scan row of resultSet#%d: %s", resultSet, err)
			}

			for i, col := range cols {
				data[resultSet][col] = vals[i]
			}
		}
		if !rows.NextResultSet() {
			break
		}
		resultSet += 1
	}

	if err := tx.Commit(); err != nil {
		return diag.Errorf("Failed to commit query: %v", err)
	}

	d.SetId("")
	d.Set("data", data)
	return nil
}
