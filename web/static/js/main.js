

function ClearRule() {
    document.getElementById("yaraRule").value = "";
}

function TestRule() {
    var rule = document.getElementById("yaraRule").value;
    var form = document.getElementById("yara");
    form.submit();
}