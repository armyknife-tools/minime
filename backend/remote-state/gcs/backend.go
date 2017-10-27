// Package gcs implements remote storage of state on Google Cloud Storage (GCS).
package gcs

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/option"
)

// gcsBackend implements "backend".Backend for GCS.
// Input(), Validate() and Configure() are implemented by embedding *schema.Backend.
// State(), DeleteState() and States() are implemented explicitly.
type gcsBackend struct {
	*schema.Backend

	storageClient  *storage.Client
	storageContext context.Context

	bucketName       string
	prefix           string
	defaultStateFile string
}

func New() backend.Backend {
	be := &gcsBackend{}
	be.Backend = &schema.Backend{
		ConfigureFunc: be.configure,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Google Cloud Storage bucket",
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the default state file",
				Deprecated:  "Use the \"prefix\" option instead",
			},

			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The directory where state files will be saved inside the bucket",
			},

			"credentials": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Google Cloud JSON Account Key",
				Default:     "",
			},
		},
	}

	return be
}

func (b *gcsBackend) configure(ctx context.Context) error {
	if b.storageClient != nil {
		return nil
	}

	// ctx is a background context with the backend config added.
	// Since no context is passed to remoteClient.Get(), .Lock(), etc. but
	// one is required for calling the GCP API, we're holding on to this
	// context here and re-use it later.
	b.storageContext = ctx

	data := schema.FromContextBackendConfig(b.storageContext)

	b.bucketName = data.Get("bucket").(string)
	b.prefix = strings.TrimLeft(data.Get("prefix").(string), "/")

	b.defaultStateFile = strings.TrimLeft(data.Get("path").(string), "/")

	opts := []option.ClientOption{
		option.WithScopes(storage.ScopeReadWrite),
		option.WithUserAgent(terraform.UserAgentString()),
	}
	if credentialsFile := data.Get("credentials").(string); credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	} else if credentialsFile := os.Getenv("GOOGLE_CREDENTIALS"); credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	}

	client, err := storage.NewClient(b.storageContext, opts...)
	if err != nil {
		return fmt.Errorf("storage.NewClient() failed: %v", err)
	}

	b.storageClient = client

	return nil
}
