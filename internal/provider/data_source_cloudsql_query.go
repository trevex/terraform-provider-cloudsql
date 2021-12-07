package provider

import (
	"context"
	"encoding/json"

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
			// TODO: well, we can't output results directly without
			//       https://github.com/hashicorp/terraform-plugin-go
			//       so should we should use it? or will there be support in sdk going forward?
			"result_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TODO",
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

	data := [][]map[string]interface{}{}
	for {
		cols, err := rows.Columns()
		if err != nil {
			return diag.Errorf("Failed to retrieve columns: %s", err)
		}
		debugLog("columns of result set: %v", cols)

		resultSetData := []map[string]interface{}{}
		for rows.Next() {
			rowData := map[string]interface{}{}
			ptrs := make([]interface{}, len(cols))
			vals := make([]interface{}, len(cols))

			for i := range ptrs {
				ptrs[i] = &vals[i]
			}

			if err := rows.Scan(ptrs...); err != nil {
				return diag.Errorf("Failed scan row: %s", err)
			}

			for i, col := range cols {
				var v interface{}
				val := vals[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				rowData[col] = v
				debugLog("setting row data `%s` = `%v`", col, v)
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

	debugLog("final computed data: %v", data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return diag.Errorf("Failed to marshal data into json: %v", err)
	}
	d.SetId(query)
	d.Set("data_json", string(jsonData))
	return nil
}
