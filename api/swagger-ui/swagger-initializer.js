window.onload = function () {
  //<editor-fold desc="Changeable Configuration Block">

  const DisableAuthorizePlugin = function () {
    return {
      wrapComponents: {
        authorizeBtn: () => () => null
      }
    }
  }
  const DisableTryItOutPlugin = function () {
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
  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  param = {
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl,
      DisableAuthorizePlugin,
      DisableTryItOutPlugin
    ],
    layout: "StandaloneLayout",
  }

  if (window.location.href.startsWith("file:")) {
    param.spec = localSpec;
  } else {
    param.url = "../openapi/v2/swagger.json";
  }

  window.ui = SwaggerUIBundle(param);

  //</editor-fold>
};
