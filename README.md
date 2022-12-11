# adepttech

Simplified API implementation for adapttech

## Installation

go get github.com/oiime/adept-tech-golang

## Usage example


#### Initiating instance with configuration 
```golang
package main

import (
	"context"
	"fmt"

	adepttech "github.com/oiime/adept-tech-golang"
)

func main() {
    api, err := adepttech.NewInstance(adepttech.Config{
		Instance:     "demo",
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		RedirectURL:  "https://localhost/oauth_continue",
	})
	if err != nil {
		panic(err)
	}
    // Returns an auth url with foobar as the state and RedirectURL
    fmt.Println(api.AuthURL("foobar"))
}

```

#### Exchanging code recieved in https://localhost/oauth_continue for access token and reloading an existing token

```golang
package main

import (
	"context"
	"fmt"

	adepttech "github.com/oiime/adept-tech-golang"
)

func main() {
    api, err := adepttech.NewInstance(adepttech.Config{
		Instance:     "demo",
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		RedirectURL:  "https://localhost/oauth_continue",
	})
	if err != nil {
		panic(err)
	}
    if err := api.ExchangeCode(context.Background(), "code"); err != nil {
		panic(err)
	}
    // At this point api contains the access token and can be used to make requests
    // You can also "save" the token and refresh token for future use along with the user using the following:
    token, err := api.Token()
	if err != nil {
		panic(err)
	}
	marshalledToken := adepttech.MarshalToken(token)
	b, err := marshalledToken.Bytes()

	// unmarshal the token
	mt, err := adepttech.UnmarshalToken(b)
	if err != nil {
		panic(err)
	}
	// use the token with the api (as an alternative to ExchangeCode)
	if err := api.AssignTokenSource(context.Background(), mt); err != nil {
		panic(err)
	}
}

```

#### Making a request using a defined type

```golang
package main

import (
	"context"
	"fmt"

	adepttech "github.com/oiime/adept-tech-golang"
)

type DatasetRequest struct {
	Network string               `json:"network"`
	Object  string               `json:"campaign"`
	Filter  DatasetFilterRequest `json:"filter"`
	Date    string               `json:"date"`
}

type DatasetFilterRequest struct {
	Condition string                     `json:"condition"`
	Rules     []DatasetFilterRuleRequest `json:"rules"`
}

type DatasetFilterRuleRequest struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    int    `json:"value"`
}

type DatasetResponse struct {
	State bool `json:"state"`
	View  bool `json:"view"`
}

func main() {
    api, err := adepttech.NewInstance(adepttech.Config{
		Instance:     "demo",
		ClientID:     "clientid",
		ClientSecret: "clientsecret",
		RedirectURL:  "https://localhost/oauth_continue",
	})
	if err != nil {
		panic(err)
	}
    // Either exchange code or unmarshal a token to api at this point
    // As GetInto and Get expect url.Values we can either pass it as-is or use the helper encode function EncodeStructAsParams
    payload := DatasetRequest{
		Network: "google",
		Object:  "campaign",
		Date:    "lifetime",
		Filter: DatasetFilterRequest{
			Condition: "and",
			Rules: []DatasetFilterRuleRequest{
				{
					Name:     "status",
					Operator: "equal",
					Value:    1,
				},
			},
		},
	}
	target := DatasetResponse{}
	if err := api.GetInto(ctx, "dataset", adepttech.EncodeStructAsParams(payload), &target); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", target)
}

```
