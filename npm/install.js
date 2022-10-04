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
const platform = PLATFORM_MAPPING[process.platform];
const arch = ARCH_MAPPING[process.arch];
const binaryName = platform === 'win32' ? `${name}.ext` : name;
const tarUrl = `https://github.com/adikari/safebox/releases/download/v${version}/safebox_${version}_${platform}_${arch}.tar.gz`;

if (!arch) {
  error(`${name} is not supported for this architecture: ${arch}`);
}

if (!platform) {
  error(`${name} is not supported for this platform: ${platform}`);
}

const bin = path.join(__dirname, "bin");

if (!fs.existsSync(bin)){
    fs.mkdirSync(bin);
}

const error = msg => {
  console.error(msg);
  process.exit(1);
};

const install = () => {
  const tmpdir = os.tmpdir();
  const req = request({ uri: tarUrl });

  const tarFile = `${tmpdir}/${name}.tar.gz`;
  const download = fs.createWriteStream(tarFile);

  console.log(`downloading safebox from ${tarUrl}`);

  req.on('response', res => {
    if (res.statusCode !== 200) {
      error(`Error downloading safebox binary from ${tarUrl}. HTTP Status Code: ${res.statusCode}`);
    }

    req.pipe(download);
  });

  req.on('complete', () => {
    cp.execSync(`tar -xf ${tarFile} -C ${tmpdir}`);
    fs.copyFileSync(path.join(tmpdir, binaryName), path.join(bin, binaryName));
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
