"use strict";

const fs = require("fs");
const mime = require("mime");

// interface S3UploadConfig {
//   bucket: {
//     name: string,
//     websiteCfgs: S3UploadWebsiteConfig,
//   },
//   resources: Array.<S3UploadResource>,
// }

// interface S3UploadWebsiteConfig {
//   index: string,
//   error: string,
// }

// interface S3UploadResource {
//   name: string,
//   source: string,
//   dest: string,
//   versionFileName: string,
// }

// refs
// - https://github.com/fernando-mc/serverless-finch/blob/master/lib/plugin.js
// - https://github.com/serverless/serverless/blob/main/lib/plugins/aws/lib/check-if-bucket-exists.js
class S3UploadPlugin {
  constructor(serverless, options) {
    this.serverless = serverless;
    this.options = options;
    this.custom = this.serverless.service.custom;
    this.provider = this.serverless.getProvider("aws");

    this.commands = {
      upload: {
        lifecycleEvents: ["upload"],
      },
    };
    this.hooks = {
      "upload:upload": async () => await this.upload(),
    };
  }

  /**
   * get configs from custom and upload files
   */
  async upload() {
    console.log("----------------------------------");
    console.log("running plugin: s3 upload");
    const s3upload = this.custom.s3upload;
    if (!this.validate(s3upload)) {
      throw error("invalid s3 upload config");
    }
    const bucketName = s3upload.bucket.name;
    const bucketWebCfgs = s3upload.bucket.websiteCfgs;
    const resources = s3upload.resources || [];
    const hasBucket = await this.checkIfBucketExists(bucketName);
    if (!hasBucket) {
      console.log(`- create new bucket: ${bucketName}`);
      await this.createBucket(bucketName, bucketWebCfgs);
    } else {
      console.log(`- bucket already exists: ${bucketName}`);
    }
    for (let resource of resources) {
      console.log(`- handle resource ${resource.name}`);
      let isSameVersion = false;
      if (hasBucket && !!resource.versionFileName) {
        console.log("  - check version");
        isSameVersion = await this.isSameVersion(
          bucketName,
          resource.source,
          resource.dest,
          resource.versionFileName
        );
      }
      if (isSameVersion) {
        console.log(
          `  - version doesn't change, skip resource ${resource.name}`
        );
        continue;
      }
      console.log("  - new version found");
      await this.emptyFolder(bucketName, resource.dest);
      await this.uploadFolder(bucketName, resource.source, resource.dest);
    }
  }

  /**
   *
   * @param {S3UploadConfig} s3upload
   * @returns
   */
  validate(s3upload) {
    if (!s3upload || !s3upload.bucket) {
      console.log("no bucket found");
      return false;
    }
    const bucketName = s3upload.bucket.name;
    if (!bucketName) {
      console.log("no bucket name found");
      return false;
    }
    const resources = s3upload.resources || [];
    if (resources.length == 0) {
      console.log("no resource found");
      return false;
    }
    for (let resource of resources) {
      if (!resource.name || !resource.source || !resource.dest) {
        return false;
      }
    }
    return true;
  }

  /**
   *
   * @param {string} bucketName
   * @returns
   */
  async checkIfBucketExists(bucketName) {
    try {
      await this.provider.request("S3", "headBucket", {
        Bucket: bucketName,
      });
      return true;
    } catch (err) {
      if (err.code === "AWS_S3_HEAD_BUCKET_NOT_FOUND") {
        return false;
      }
      throw err;
    }
  }

  /**
   *
   * @param {string} bucketName
   * @param {S3UploadWebsiteConfig} websiteCfgs
   * @returns
   */
  async createBucket(bucketName, websiteCfgs) {
    console.log(`  - create bucket: ${bucketName}`);
    await this.provider.request("S3", "createBucket", {
      Bucket: bucketName,
      // ACL: acl || "public",
    });
    console.log(`  - put bucket policy: ${bucketName}`);
    await this.provider.request("S3", "putBucketPolicy", {
      Bucket: bucketName,
      Policy: JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Effect: "Allow",
            Principal: { AWS: "*" },
            Action: "s3:GetObject",
            Resource: `arn:aws:s3:::${bucketName}/*`,
          },
        ],
      }),
    });
    if (!websiteCfgs) {
      return;
    }
    // enable website
    console.log(`  - enable website: ${bucketName}`);
    const params = {
      Bucket: bucketName,
      // https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketWebsite.html#API_PutBucketWebsite_RequestSyntax
      WebsiteConfiguration: {
        IndexDocument: {
          Suffix: websiteCfgs.index || "index.html",
        },
        ErrorDocument: {
          Key: websiteCfgs.error || "error.html",
        },
      },
    };
    await this.provider.request("S3", "putBucketWebsite", params);
  }

  /**
   *
   * @param {string} bucketName
   * @param {string} source
   * @param {string} dest
   * @param {string} versionFilename
   * @returns
   */
  async isSameVersion(bucketName, source, dest, versionFilename) {
    let oldver = await this.getObject(bucketName, `${dest}/${versionFilename}`);
    let newver = await this.readFile(`${source}/${versionFilename}`);
    oldver = oldver || "";
    newver = newver || "";
    return oldver == newver;
  }

  /**
   *
   * @param {string} bucketName
   * @param {string} objectKey
   * @returns
   */
  async getObject(bucketName, objectKey) {
    try {
      const result = await this.provider.request("S3", "getObject", {
        Bucket: bucketName,
        Key: objectKey,
      });
      return !result ? null : String(result.Body);
    } catch (error) {
      if (error.code === "AWS_S3_GET_OBJECT_NO_SUCH_KEY") {
        return null;
      }
      throw error;
    }
  }

  /**
   *
   * @param {string} bucketName
   * @param {string} folder
   */
  async emptyFolder(bucketName, folder) {
    console.log(`  - empty folder ${folder} in bucket ${bucketName}`);
    const result = await this.provider.request("S3", "listObjectsV2", {
      Bucket: bucketName,
      Prefix: folder,
    });
    if (!result || !result.Contents || result.Contents.length == 0) {
      return;
    }
    let objects = [];
    for (let content of result.Contents) {
      objects.push({ Key: content.Key });
    }
    await this.provider.request("S3", "deleteObjects", {
      Bucket: bucketName,
      Delete: {
        Objects: objects,
        Quiet: true,
      },
    });
  }

  /**
   *
   * @param {string} bucketName
   * @param {string} folder
   */
  async uploadFolder(bucketName, folder, dest) {
    console.log(`  - read folder ${folder}`);
    const filenames = (await this.readFolder(folder)) || [];
    const fileMap = (await this.readFiles(filenames)) || {};
    if (fileMap.size == 0) {
      return;
    }
    console.log(`  - upload folder to bucket ${bucketName}/${dest}`);
    for (let filename in fileMap) {
      const path = filename.substring(folder.length + 1);
      console.log(`    - upload ${filename} to ${bucketName}/${dest}/${path}`);
      await this.provider.request("S3", "putObject", {
        Bucket: bucketName,
        Key: !!dest ? `${dest}/${path}` : path,
        Body: fileMap[filename],
        ContentType: mime.getType(filename),
      });
    }
  }

  /**
   *
   * @param {string} folder
   * @returns
   */
  async readFolder(folder) {
    const dirents = fs.readdirSync(`./${folder}`, {
      encoding: "utf-8",
      withFileTypes: true,
    });
    if (!dirents) {
      return {};
    }
    let filenames = [];
    for (let dirent of dirents) {
      if (dirent.isDirectory()) {
        const subfiles = await this.readFolder(`${folder}/${dirent.name}`);
        filenames.push(...subfiles);
      } else {
        filenames.push(`${folder}/${dirent.name}`);
      }
    }
    return filenames;
  }

  /**
   *
   * @param {Arrray.<string>} filenames
   * @returns
   */
  async readFiles(filenames) {
    let filemap = {};
    for (let filename of filenames) {
      const content = await this.readFile(filename);
      filemap[filename] = content;
    }
    return filemap;
  }

  /**
   *
   * @param {string} filename
   * @returns
   */
  async readFile(filename) {
    try {
      return fs.readFileSync(`./${filename}`, { encoding: "utf-8" });
    } catch (error) {
      console.log(error);
      return null;
    }
  }
}

module.exports = S3UploadPlugin;
