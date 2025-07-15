import "node:fs";
import { readFileSync, readFile } from 'node:fs';
import { WASI } from 'node:wasi';

const wasi = new WASI({
    version: 'preview1',
    preopens: {
        '/data': process.cwd()
    },
    stderr: 1,
    stdout: 1
});

const importObject = { wasi_snapshot_preview1: wasi.wasiImport };

var module = null;

function insertText(text) {
   // Get the address of the writable memory.
   let addr = module.exports.MkBuffer()
   let buffer = module.exports.memory.buffer

   let mem = new Int8Array(buffer)
   let view = mem.subarray(addr, addr + text.length)

   for (let i = 0; i < text.length; i++) {
      view[i] = text.charCodeAt(i)
   }

   // Return the address we started at.
   return [addr, text.length];
}

function textArguments(...texts) {
    let args = [];

    for (let text of texts) {
        let [addr, length] = insertText(text);
        args.push(addr, length);
    }

    return args;
}


(async() => {
    const wasm = await WebAssembly.compile(readFileSync("wasi.wasm"));

    console.log("WASM compiled");

    module = await WebAssembly.instantiate(wasm, importObject);

    console.log("WASM instantiated");

    wasi.start(module);

    // instance.exports._start();

    // wasi.start(instance);

    console.log("WASI started");

    await module.exports.Initialize();

    await module.exports.Configure(...textArguments(
        "",
        "",
        "",
        "",
    ));

    module.exports.Sign(...textArguments("/data/test.ps1"));
})();