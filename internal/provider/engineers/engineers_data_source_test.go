package engineers_test

import (
	"testing"
	"fmt"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"terraform-provider-devops/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Local acceptance test helpers for this package. We cannot rely on unexported
// variables from provider/*_test.go across packages, since those are not
// compiled into the provider package for import. Define local equivalents here.


var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
    "dob": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccEngineersDataSource(t *testing.T) {

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch {
        case r.Method == http.MethodGet && r.URL.Path == "/engineers":
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            _ = json.NewEncoder(w).Encode([]map[string]string{
                {"id": "CJ123", "name": "Colin",  "email": "Colin@liatrio.com"},
                {"id": "AN123", "name": "Austin", "email": "Austin@liatrio.com"},
                {"id": "JM123", "name": "Jack",   "email": "Jack@liatrio.com"},
                {"id": "MW123", "name": "Madi",   "email": "Madi@liatrio.com"},
                {"id": "AG123", "name": "Angel",  "email": "Angel@liatrio.com"},
            })
            return
        default:
            http.NotFound(w, r)
            return
        }
    }))
    defer server.Close()

    cfg := fmt.Sprintf(`
provider "dob" {
  endpoint = %q
}
data "dob_engineer" "test" {}
`, server.URL)

    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Read testing
            {
                Config: cfg,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of engineers returned
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.#", "5"),
                    // Verify the first engineer to ensure all attributes are set
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.0.name", "Colin"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.0.id", "CJ123"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.0.email", "Colin@liatrio.com"),

                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.1.name", "Austin"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.1.id", "AN123"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.1.email", "Austin@liatrio.com"),

                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.2.name", "Jack"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.2.id", "JM123"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.2.email", "Jack@liatrio.com"),

                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.3.name", "Madi"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.3.id", "MW123"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.3.email", "Madi@liatrio.com"),

                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.4.name", "Angel"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.4.id", "AG123"),
                    resource.TestCheckResourceAttr("data.dob_engineer.test", "engineers.4.email", "Angel@liatrio.com"),
                ),
            },
        },
    })
}
