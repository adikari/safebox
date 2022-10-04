#!/usr/bin/env node

'use strict';

const request = require('request'),
  os = require('os'),
  fs = require('fs'),
  path = require('path'),
  cp = require('child_process'),
  packageJson = require('./package.json');

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
  ia32: '386',
  x64: 'amd64',
  arm: 'arm',
  arm64: 'arm64'
};

// Mapping between Node's `process.platform` to Golang's
const PLATFORM_MAPPING = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows',
  freebsd: 'freebsd'
};

const name = 'safebox';
const version = packageJson.version;
const platform = process.platform;
const arch = process.arch;
const binaryName = platform === 'win32' ? `${name}.ext` : name;
const tarUrl = `https://github.com/adikari/safebox/releases/download/v${version}/safebox_${version}_${platform}_${arch}.tar.gz`;

const nodeBin = cp.execSync("npm bin").toString().replace(/\r?\n|\r/g, "");

const error = msg => {
  console.error(msg);
  process.exit(1);
};

if (!(arch in ARCH_MAPPING)) {
  error(`${name} is not supported for this architecture: ${arch}`);
  return;
}

if (!(platform in PLATFORM_MAPPING)) {
  error(`${name} is not supported for this platform: ${platform}`);
  return;
}

const install = () => {
  const tmpdir = os.tmpdir();
  const req = request({ uri: tarUrl });

  const tarFile = `${tmpdir}/${name}.tar.gz`;
  const download = fs.createWriteStream(tarFile);

  req.on('response', res => {
    if (res.statusCode !== 200) {
      return callback(`Error downloading safebox binary. HTTP Status Code: ${res.statusCode}`);
    }

    req.pipe(download);
  });

  req.on('complete', () => {
    cp.execSync(`tar -xf ${tarFile} -C ${tmpdir}`);
    fs.copyFileSync(path.join(tmpdir, binaryName), path.join(nodeBin, binaryName));
  });
};

install();
