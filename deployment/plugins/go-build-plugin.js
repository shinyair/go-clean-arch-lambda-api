"use strict";

const { execSync } = require("child_process");
const zip = require("cross-zip");
const rimraf = require("rimraf");

class GoBuildPlugin {
  constructor(serverless, options) {
    this.serverless = serverless;
    this.options = options;
    this.functions = this.serverless.service.functions;

    this.hooks = {
      initialize: () => this.buildPackage(),
    };
  }

  /**
   * build and package executable file for each lambda function
   */
  buildPackage() {
    console.log("----------------------------------");
    console.log("running plugin: go build");
    const currDir = process.cwd();
    execSync("set GOARCH=amd64");
    execSync("set GOOS=linux");
    for (let lambdaKey in this.functions) {
      console.log(`- build for lambda: ${lambdaKey}`);
      const lambdaFunc = this.functions[lambdaKey];
      const buildCmd = `go build -o ./bin/${lambdaKey}/main ${lambdaFunc.handler}`;
      console.log(`  - ${buildCmd}`);
      execSync(buildCmd);
      console.log("  - build done");
      const zipName = lambdaFunc.package.artifact;
      console.log(`  - zip executable ./bin/${lambdaKey}/main as ${zipName}`);
      zip.zipSync(`${currDir}/bin/${lambdaKey}/main`, `${currDir}/${zipName}`);
      console.log("  - zip done");
      lambdaFunc.handler = "main";
      console.log("  - update handler as main");
    }
    rimraf.sync(`${currDir}/bin`);
    console.log("- remove bin");
  }
}

module.exports = GoBuildPlugin;
