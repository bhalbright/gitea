package graph

import (
	"code.gitea.io/gitea/graph/generated"
	giteaCtx "code.gitea.io/gitea/modules/context"

	handler2 "github.com/99designs/gqlgen/graphql/handler"
)

//var (
//	schema graphql.Schema
//)

//type reqBody struct {
//	Query string `json:"query"`
//}

//type contextKeyType string

// Init initializes gql server
//func Init(cfg graphql.Schema) {
//	schema = cfg
//}

func GraphQL(ctx *giteaCtx.APIContext) {
	// NewExecutableSchema and Config are in the generated.go file
	c := generated.Config{
		Resolvers: &Resolver{},
	}

	//h := handler.GraphQL(gql.NewExecutableSchema(c))
	h := handler2.New(generated.NewExecutableSchema(c))

	h.ServeHTTP(ctx.Resp, ctx.Req.Request)
	//return func(c *gin.Context) {
	//	h.ServeHTTP(c.Writer, c.Request)
	//}
}
	/*
	// Check to ensure query was provided in the request body
	if ctx.Req.Body() == nil {
		ctx.Error(http.StatusBadRequest, "", "Must provide graphql query in request body")
		return
	}
	var rBody reqBody
	bodyString, err := ctx.Req.Body().String()
	if err != nil {
		ctx.Error(http.StatusBadRequest, "", "Error reading request body")
		return
	}
	// Decode the request body into rBody
	err = json.NewDecoder(strings.NewReader(bodyString)).Decode(&rBody)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "", "Error parsing JSON request body")
		return
	}

	// Execute graphql query
	result := ExecuteQuery(rBody.Query, schema, ctx)

	ctx.JSON(http.StatusOK, result)

}

	 */
/*
// ExecuteQuery runs our graphql queries
func ExecuteQuery(query string, schema graphql.Schema, ctx *giteaCtx.APIContext) *graphql.Result {
	apiContextKey := contextKeyType("giteaApiContext")
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context: context.WithValue(context.Background(), apiContextKey, ctx),
		RootObject: make(map[string]interface{}),
	})

	if len(result.Errors) > 0 {
		log.Error("Unexpected errors inside ExecuteQuery: %v", result.Errors)
	}

	return result
}
*/
