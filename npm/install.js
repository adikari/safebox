#!/usr/bin/env node

'use strict';

const request = require('request'),
  os = require('os'),
  fs = require('fs'),
  path = require('path'),
  cp = require('child_process'),
  constants = require('./constants');

const { name, platform, arch, binaryName, bin, tarUrl } = constants;

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
  const req = request({ uri: tarUrl });

  const tarFile = `${tmpdir}/${name}-${Date.now()}.tar.gz`;
  const download = fs.createWriteStream(tarFile);

  if (retries > 0) {
    console.log(`retrying to install safebox - retry ${retries} out of ${MAX_RETRIES}`)
  }

  console.log(`downloading safebox binary`);

  req.on('response', res => {
    if (res.statusCode !== 200) {
      error(`Error downloading safebox binary. HTTP Status Code: ${res.statusCode}`);
    }

    req.pipe(download);
  });

  req.on('complete', () => {
    console.log('download complete. installing safebox.')

    try {
      if (!fs.existsSync(tarFile)) {
        throw new Error(`${tarFile} does not exist`)
      }

      cp.execSync(`tar -xf ${tarFile} -C ${tmpdir}`);
      fs.copyFileSync(path.join(tmpdir, binaryName), path.join(bin, binaryName));
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
