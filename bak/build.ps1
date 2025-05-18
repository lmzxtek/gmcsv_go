
Write-Host "==== Version Selection ======" -ForegroundColor Cyan
Write-Host "  1. Release                 "
Write-Host "  2. Debug                   "
Write-Host "=============================" -ForegroundColor Cyan
$soft_choice = Read-Host " Enter your choice (default: 1 )"  
if ($soft_choice -eq "2") {
    $env:GIN_MODE="debug"; 
}else{
    $env:GIN_MODE="release"; 
}

Write-Host "======== Arch Selection =====" -ForegroundColor Cyan
Write-Host "  1. Amd64                   "
Write-Host "  2. Arm64                   "
Write-Host "  3. i386                    "
Write-Host "=============================" -ForegroundColor Cyan
$soft_choice = Read-Host " Enter your choice (default: 1 )"  
if ($soft_choice -eq "2") {
    $env:GOARCH="arm64"; 
}elseif ($soft_choice -eq "3") {
    $env:GOARCH="i386"; 
}else{
    $env:GOARCH="amd64"; 
}

Write-Host "==== Platform Selection =====" -ForegroundColor Cyan
Write-Host "  1. Windows                 "
Write-Host "  2. Linux                   "
Write-Host "=============================" -ForegroundColor Cyan
$soft_choice = Read-Host " Enter your choice (default: 1 )"  
if ($soft_choice -eq "2") {
    $env:GOOS="linux"; 
    # $file = "gm_csv-${env:GOOS}"; 
}else{
    $env:GOOS="windows"; 
    # $file = "gm_csv-windowns.exe"; 
}

$file = "gm_csv-${env:GOOS}-${env:GOARCH}"; 
if ($env:GOOS -eq "windows") { $file = "${file}.exe"; }

$bin_path = "./bin"
if (-not (Test-Path $bin_path)) { 
    New-Item -ItemType Directory -Force -Path $bin_path | Out-Null
    Write-Host "Created $bin_path directory" -ForegroundColor Green
}

$is_to_build = "true"

$fpath = "./bin/$file"
if (Test-Path $fpath) { 
    $choice = Write-Host "File already exists, overwrite? (Y/n)" -ForegroundColor Yellow
    if ($choice.ToLower() -eq "no" || $choice.ToLower() -eq "n") {
        $is_to_build = "false"
    }else{
        Write-Host "Aborted build !" -ForegroundColor Red
        exit 1
    }
}

if ($is_to_build -eq "true") {

    Write-Host "Building $file ..." -ForegroundColor Cyan

    go build -o $fpath -v -ldflags "-s -w" main.go

    Write-Host "Build Complete!" -ForegroundColor Green

}else{
    Write-Host "Skipping build" -ForegroundColor Yellow
}