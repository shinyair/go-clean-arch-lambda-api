"use strict";

// refs
// - https://www.serverless.com/framework/docs/guides/plugins/creating-plugins
// - https://github.com/mvila/serverless-plugin-scripts/blob/master/src/index.js
// - https://github.com/sean9keenan/serverless-go-build/blob/master/index.js

class PrintInfoPlugin {
  constructor(serverless, options) {
    this.serverless = serverless;
    this.options = options;

    this.hooks = {
      initialize: () => this.init(),
    };
  }

  /**
   * example
   */
  init() {
    console.log("----------------------------------");
    console.log("running plugin: print info");
    // console.log('Serverless instance: ', this.serverless);
    // `serverless.service` contains the (resolved) serverless.yml config
    const service = this.serverless.service;
    console.log("Provider: ", service.provider);
    // console.log('Custom: ', service.custom);
    console.log("Functions: ", service.functions);
    console.log("Resources: ", service.resources);
  }
}

module.exports = PrintInfoPlugin;
