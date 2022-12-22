"use strict";

String.prototype.interpolate = function (params) {
  const names = Object.keys(params);
  const vals = Object.values(params);
  const func = new Function(...names, `return \`${this}\`;`);
  return func(...vals);
};

class ParseVarPlugin {
  constructor(serverless, options) {
    this.serverless = serverless;
    this.options = options;
    this.serverless.service.provider.environment =
      this.serverless.service.provider.environment || {};
    this.custom = this.serverless.service.custom;
    this.serviceProvider = this.serverless.service.provider;
    this.resources = this.serverless.service.resources.Resources;
    this.environment = this.serverless.service.provider.environment;

    this.hooks = {
      initialize: () => this.parse(),
    };
  }

  /**
   * entrance function
   */
  parse() {
    console.log("----------------------------------");
    console.log("running plugin: parse variables");
    const vars = this.custom.parseVar || [];

    for (let customVar of vars) {
      const varName = customVar.name;
      const parsed = this.parseVar(
        customVar.name,
        customVar.source,
        customVar.skipFields || []
      );
      if (customVar.addToEnv) {
        this.addToEnv(varName, parsed);
      }
    }
  }

  /**
   *
   * @param {string} varName
   * @param {string} varSource
   * @param {Array.<string>} skipFields
   * @returns
   */
  parseVar(varName, varSource, skipFields) {
    console.log(`  - parse variables under ${varName}`);
    let parent = undefined;
    if (varSource == "custom") {
      parent = this.custom;
    } else if (varSource == "resource") {
      parent = this.resources;
    } else if (varSource == "provider") {
      parent = this.serviceProvider;
    } else {
      console.log(`  - unspported source type: ${varSource}`);
      return "";
    }
    let customVarMap = {
      stage: this.custom.stage,
      variant: this.custom.variant,
      appCode: this.custom.appCode,
    };
    let value = parent[varName];
    let fvalue = value;
    switch (typeof value) {
      case "string":
        fvalue = this.parseString("", varName, value, skipFields, customVarMap);
        break;
      case "object":
        if (value instanceof Array) {
          fvalue = this.parseArray("", value, skipFields, customVarMap);
        } else {
          fvalue = this.parseMap("", value, skipFields, customVarMap);
        }
        break;
      default: {
        console.log(`    - unsupported type for variable: ${varName}`);
      }
    }
    parent[varName] = fvalue;
    return fvalue;
  }

  /**
   *
   * @param {string} path
   * @param {Array.<string>} valArray
   * @param {Object.<string,string>} skipFields
   * @param {Object.<string,string>} customVarMap
   * @returns
   */
  parseArray(path, valArray, skipFields, customVarMap) {
    console.log(`    - array type not supported: ${path}`);
    return valArray;
  }

  /**
   *
   * @param {string} path
   * @param {Object.<string,string>} valMap
   * @param {Array.<string>} skipFields
   * @param {Object.<string,string>} customVarMap
   * @returns
   */
  parseMap(path, valMap, skipFields, customVarMap) {
    let fvalMap = {};
    for (let key in valMap) {
      const val = valMap[key];
      const keyPath = !!path ? `${path}.${key}` : key;
      switch (typeof val) {
        case "string":
          const fval1 = this.parseString(
            path,
            key,
            val,
            skipFields,
            customVarMap
          );
          fvalMap[key] = fval1;
          customVarMap[keyPath] = fval1;
          break;
        case "object":
          let fval2 = val;
          if (val instanceof Array) {
            fval2 = this.parseArray(keyPath, val, skipFields, customVarMap);
          } else {
            fval2 = this.parseMap(keyPath, val, skipFields, customVarMap);
          }
          fvalMap[key] = fval2;
          break;
        default: {
          fvalMap[key] = val;
          console.log(`    - skip val: ${keyPath}`);
        }
      }
    }
    return fvalMap;
  }

  /**
   *
   * @param {string} path
   * @param {string} key
   * @param {string} val
   * @param {Array.<string>} skipFields
   * @param {Object.<string, string>} customVarMap
   * @returns
   */
  parseString(path, key, val, skipFields, customVarMap) {
    let fval = val;
    if (!skipFields.includes(key)) {
      fval = String(val).interpolate(customVarMap);
      const keyPath = !!path ? `${path}.${key}` : key;
      console.log(`    - format ${keyPath} from ${val} to ${fval}`);
    }
    return fval;
  }

  /**
   *
   * @param {string} varName
   * @param {string} varValue
   */
  addToEnv(varName, varValue) {
    if (typeof varValue == "object") {
      for (let attrKey in varValue) {
        this.environment[attrKey] = varValue[attrKey];
      }
    } else {
      this.environment[varName] = varValue;
    }
  }
}

module.exports = ParseVarPlugin;
