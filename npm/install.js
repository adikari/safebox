#!/usr/bin/env node

'use strict';

const request = require('request'),
  os = require('os'),
  fs = require('fs'),
  path = require('path'),
  constants = require('./constants');

const { name, platform, arch, binaryName, bin, binaryUrl } = constants;

if (!arch) {
  error(`${name} is not supported for this architecture: ${arch}`);
}

if (!platform) {
  error(`${name} is not supported for this platform: ${platform}`);
}

if (!fs.existsSync(bin)){
    fs.mkdirSync(bin);
}

let retries = 0;
const MAX_RETRIES = 3;

const install = () => {
  const tmpdir = os.tmpdir();
  const binary = path.join(tmpdir, binaryName);

  const copyBinary = () => fs.copyFileSync(binary, path.join(bin, binaryName));

  if (fs.existsSync(binary)) {
    console.log('safebox binary already downloaded');
    copyBinary();
    return;
  }

  const req = request({ uri: binaryUrl });
  if (retries > 0) {
    console.log(`retrying to install safebox - retry ${retries} out of ${MAX_RETRIES}`)
  }

  console.log(`downloading safebox binary`);

  const download = fs.createWriteStream(binary);

  req.on('response', res => {
    if (res.statusCode !== 200) {
      error(`Error downloading safebox binary. HTTP Status Code: ${res.statusCode}`);
    }

    req.pipe(download);
  });

  req.on('complete', () => {
    console.log('download complete. installing safebox.')

    try {
      if (!fs.existsSync(binary)) {
        throw new Error(`${binary} does not exist`)
      }

      copyBinary();
    } catch (error) {
      console.error('failed to extract binary.', error.message)

      retries += 1;
      if (retries <= MAX_RETRIES) {
        install();
      }
    }
  });
};

const uninstall = () => {
  fs.unlinkSync(path.join(bin, binaryName));
}

let actions = {
    "install": install,
    "uninstall": uninstall
};

let argv = process.argv;
if (argv && argv.length > 2) {
    let cmd = process.argv[2];
    if (!actions[cmd]) {
        error("Invalid command. `install` and `uninstall` are the only supported commands");
    }

    actions[cmd]();
}
