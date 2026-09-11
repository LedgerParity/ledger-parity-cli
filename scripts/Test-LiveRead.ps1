param([string]$Binary = (Join-Path $PSScriptRoot '..\ledger-parity.exe'))
$ErrorActionPreference = 'Stop'
# Read-only public testnet check. The expected application row is synthesized
# from a live observation, so this cannot establish real application adoption.
$endpoint = 'https://horizon-testnet.stellar.org'
$network = 'Test SDF Network ; September 2015'
$root = Invoke-RestMethod -Uri "$endpoint/" -TimeoutSec 25
if ($root.network_passphrase -ne $network) { throw 'Unexpected network' }
$page = Invoke-RestMethod -Uri "$endpoint/payments?limit=200&order=desc" -TimeoutSec 25
$payment = $page._embedded.records | Where-Object { $_.type -eq 'payment' -and $_.transaction_successful -and -not $_.from_muxed -and -not $_.to_muxed } | Select-Object -First 1
if (-not $payment) { throw 'No ordinary payment in the latest 200 operations; retry later' }
$dir = Join-Path $PSScriptRoot '..\.local-checks'
New-Item -ItemType Directory -Force $dir | Out-Null
$dir = (Resolve-Path -LiteralPath $dir).Path
$encoding = New-Object System.Text.UTF8Encoding($false)
$assetCode = $payment.asset_code
if ($payment.asset_type -eq 'native') { $assetCode = 'XLM' }
$expected = @{
 id='synthetic-expectation-from-live-observation';source_app='live-read-check';network=$network
 operation_type='payment';operation_id=[string]$payment.id;reference_id=[string]$payment.transaction_hash
 sender=$payment.from;recipient=$payment.to;amount=[string]$payment.amount
 asset=$assetCode;asset_type=$payment.asset_type;asset_issuer=[string]$payment.asset_issuer
 timestamp=$payment.created_at;status='completed'
}
$internalPath=Join-Path $dir 'live-internal.json'
[IO.File]::WriteAllText($internalPath, (ConvertTo-Json -InputObject @($expected) -Depth 8), $encoding)
$closed = [DateTimeOffset]::Parse($payment.created_at)
$reportPath=Join-Path $dir 'live-report.json'
$configPath=Join-Path $dir 'live-config.json'
$config = @{
 target_app=@{name='live-read-check';format='json';source_path=$internalPath;complete=$false}
 stellar=@{network=$network;accounts=@($payment.from);horizon_url=$endpoint}
 reconciliation=@{start=$closed.AddSeconds(-1).ToString('o');end=$closed.ToString('o');timeframe_tolerance_sec=0}
 output=@{format='json';file_path=$reportPath}
}
[IO.File]::WriteAllText($configPath, ($config | ConvertTo-Json -Depth 8), $encoding)
& $Binary --config $configPath
if ($LASTEXITCODE -ne 0) { throw "Live read failed with exit $LASTEXITCODE" }
$report=Get-Content -LiteralPath $reportPath -Raw | ConvertFrom-Json
if ($report.total_matched -ne 1 -or $report.total_unknown -ne 0 -or $report.total_discrepancies -ne 0) { throw 'Unexpected live report' }
$evidence=@{
 observed_at=[DateTimeOffset]::UtcNow.ToString('o');network=$network;endpoint=$endpoint
 operation_id=[string]$payment.id;transaction_hash=[string]$payment.transaction_hash;payment_timestamp=$payment.created_at
 matched=$report.total_matched;discrepancies=$report.total_discrepancies;unknown=$report.total_unknown
 coverage=$report.coverage
 method='GET latest public testnet payments; synthesize one expected row; independently refetch account payments through CLI and compare exact operation identity/amount/asset/direction'
 limitations='Read-only provider observation. Expected application data derived from that observation; not an independent business record, customer integration, transaction submission or production-readiness proof. Testnet can reset.'
}
[IO.File]::WriteAllText((Join-Path $dir 'live-evidence.json'), ($evidence | ConvertTo-Json -Depth 8), $encoding)
$evidence | ConvertTo-Json -Depth 8
