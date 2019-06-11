
function Create-RandomFile {
param([String] $OutFile, [UInt64] $size)
    $rand = New-Object System.Security.Cryptography.RNGCryptoServiceProvider
    $b = New-Object "System.Byte[]" $size
    $rand.GetBytes($b)
    [System.IO.File]::WriteAllBytes($OutFile, $b)
}			


$basePath = Resolve-Path '.'

Create-RandomFile -OutFile (Join-Path $basePath '.\random1k.bin') -size (1 * 1024)
Create-RandomFile -OutFile (Join-Path $basePath '.\random3k.bin') -size (3 * 1024)
Create-RandomFile -OutFile (Join-Path $basePath '.\random5k.bin') -size (5 * 1024)

Create-RandomFile -OutFile (Join-Path $basePath '.\random1M.bin') -size (1 * 1024 * 1024)
Create-RandomFile -OutFile (Join-Path $basePath '.\random3M.bin') -size (3 * 1024 * 1024)
Create-RandomFile -OutFile (Join-Path $basePath '.\random5M.bin') -size (5 * 1024 * 1024)
Create-RandomFile -OutFile (Join-Path $basePath '.\random20M.bin') -size (20 * 1024 * 1024)
Create-RandomFile -OutFile (Join-Path $basePath '.\random21M.bin') -size (21 * 1024 * 1024)

Create-RandomFile -OutFile (Join-Path $basePath '.\random1G.bin') -size (1 * 1024 * 1024 * 1024)
