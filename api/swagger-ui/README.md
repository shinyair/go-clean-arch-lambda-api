# Swagger UI
## Source
- github: https://github.com/swagger-api/swagger-ui
- static resources: https://github.com/swagger-api/swagger-ui/tree/master/dist
## Customizations
- disable try-it-out feature
- disable authorize feature
```
  const DisableAuthorizePlugin = function() {
    return {
      wrapComponents: {
        authorizeBtn: () => () => null
      }
    }
  }
  const DisableTryItOutPlugin = function() {
    return {
      statePlugins: {
        spec: {
          wrapSelectors: {
            allowTryItOutFor: () => () => false
          }
        }
      }
    }
  }
  ...
  window.ui = SwaggerUIBundle({
    ...
    plugins: [
	  DisableAuthorizePlugin,
	  DisableTryItOutPlugin
    ],
    ...
  });
```