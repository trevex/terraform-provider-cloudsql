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
			"data": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
					},
				},
				Computed:    true,
				Description: "List of lists of map of strings",
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
	debugLog("running sql query: %s", query)
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return diag.Errorf("Failed to exec query `%s`, returned: %v", query, err)
	}
	defer rows.Close()

	data := [][]map[string]string{}
	for {
		cols, err := rows.Columns()
		if err != nil {
			return diag.Errorf("Failed to retrieve columns: %s", err)
		}
		debugLog("columns of result set: %v", cols)

		resultSetData := []map[string]string{}
		for rows.Next() {
			rowData := map[string]string{}
			ptrs := make([]interface{}, len(cols))
			vals := make([]string, len(cols))

			for i := range ptrs {
				ptrs[i] = &vals[i]
			}

			if err := rows.Scan(ptrs...); err != nil {
				return diag.Errorf("Failed scan row: %s", err)
			}

			for i, col := range cols {
				rowData[col] = vals[i]
				debugLog("setting row data `%s` = `%s`", col, vals[i])
			}
			resultSetData = append(resultSetData, rowData)
		}
		data = append(data, resultSetData)
		if !rows.NextResultSet() {
			break
		}
	}

	if err := tx.Commit(); err != nil {
		return diag.Errorf("Failed to commit query: %v", err)
	}

	d.SetId(query)
	debugLog("final computed data: %v", data)
	d.Set("data", data)
	return nil
}
