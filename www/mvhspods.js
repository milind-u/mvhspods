function runGo() {
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("mvhspods.wasm", go.importObject)).then(result => go.run(result.instance));
}

runGo();

function doMakePods() {
  const files = this.files;
  if (files.length !== 1) {
    console.error("Expected 1 file");
  }

  const reader = new FileReader();
  reader.onload = e => {
    const pods = makePods(e.target.result);

    // Write the pods as a csv to disk
    const elem = document.createElement('a');
    elem.href = "data:test/csv;charset=utf-8," + encodeURIComponent(pods);
    elem.download = "pods.csv";
    elem.style.display = "none";
    document.body.appendChild(elem);
    elem.click();
    document.body.removeChild(elem);
  };
  reader.readAsText(files[0]);
}
