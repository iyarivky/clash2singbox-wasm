// https://github.com/torch2424/wasm-by-example/blob/master/demo-util/
export const wasmBrowserInstantiate = async (wasmModuleUrl, importObject) => {
    let response = undefined;
  
    // Check if the browser supports streaming instantiation
    if (WebAssembly.instantiateStreaming) {
      // Fetch the module, and instantiate it as it is downloading
      response = await WebAssembly.instantiateStreaming(
        fetch(wasmModuleUrl),
        importObject
      );
    } else {
      // Fallback to using fetch to download the entire module
      // And then instantiate the module
      const fetchAndInstantiateTask = async () => {
        const wasmArrayBuffer = await fetch(wasmModuleUrl).then(response =>
          response.arrayBuffer()
        );
        return WebAssembly.instantiate(wasmArrayBuffer, importObject);
      };
  
      response = await fetchAndInstantiateTask();
    }
  
    return response;
  };

  const go = new Go(); // Defined in wasm_exec.js. Don't forget to add this in your index.html.

  const runWasmAdd = async () => {
    // Get the importObject from the go instance.
    const importObject = go.importObject;
  
    // Instantiate our wasm module
    const wasmModule = await wasmBrowserInstantiate("./main.wasm", importObject);
  
    // Allow the wasm_exec go instance, bootstrap and execute our wasm module
    go.run(wasmModule.instance);
  
    // Call the Add function export from wasm, save the result
    const configFile = `port: 7890
    socks-port: 7891
    redir-port: 7892
    mixed-port: 7893
    tproxy-port: 7895
    ipv6: false
    mode: rule
    log-level: silent
    allow-lan: true
    external-controller: 0.0.0.0:9090
    secret: ""
    bind-address: "*"
    unified-delay: true
    profile:
      store-selected: true
    dns:
      enable: true
      ipv6: false
      enhanced-mode: redir-host
      listen: 0.0.0.0:7874
      nameserver:
        - 8.8.8.8
        - 1.0.0.1
        - https://dns.google/dns-query
      fallback:
        - 1.1.1.1
        - 8.8.4.4
        - https://cloudflare-dns.com/dns-query
        - 112.215.203.254
      default-nameserver:
        - 8.8.8.8
        - 1.1.1.1
        - 112.215.203.254
    proxies:
      - name: RELAY-104.17.253.39-0991
        server: cfcdn2.sanfencdn.net
        port: 2052
        type: vmess
        uuid: 78d6fa88-0d27-414f-8f51-1e16186d588f
        alterId: 0
        cipher: auto
        tls: false
        skip-cert-verify: true
        servername: us8.sanfencdn2.com
        network: ws
        ws-opts:
          path: /
          headers:
            Host: us8.sanfencdn2.com
        udp: true
    proxy-groups:
      - name: FASTSSH-SSHKIT-HOWDY
        type: select
        proxies:
          - RELAY-104.17.253.39-0991
          - DIRECT
    rules:
      - MATCH,FASTSSH-SSHKIT-HOWDY
    `
    const addResult = wasmModule.instance.exports.configClash(configFile);
  
    // Set the result onto the body
    document.body.textContent = `Hello World! addResult: ${addResult}`;
  };
  runWasmAdd();