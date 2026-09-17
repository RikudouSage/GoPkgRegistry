const define = function (name, value) {
  window.runtimeVariables ??= {};
  window.runtimeVariables[name] = value;
}
