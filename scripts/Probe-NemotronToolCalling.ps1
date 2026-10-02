Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

$SecretCandidates=@(
  'C:\ProgramData\CKBR\research-sidecar-yggdrasil\secrets\openrouter.dpapi',
  'C:\ProgramData\CKBR\research-sidecar\secrets\openrouter.dpapi'
)
$Secret=$SecretCandidates|Where-Object{Test-Path -LiteralPath $_ -PathType Leaf}|Select-Object -First 1
if([string]::IsNullOrWhiteSpace($Secret)){throw 'OPENROUTER_SECRET_MISSING'}

$Secure=(Get-Content -LiteralPath $Secret -Raw).Trim()|ConvertTo-SecureString
$Ptr=[Runtime.InteropServices.Marshal]::SecureStringToBSTR($Secure)
try{$Token=[Runtime.InteropServices.Marshal]::PtrToStringBSTR($Ptr)}finally{[Runtime.InteropServices.Marshal]::ZeroFreeBSTR($Ptr)}
if([string]::IsNullOrWhiteSpace($Token)){throw 'OPENROUTER_SECRET_EMPTY'}

$Headers=@{
  Authorization='Bearer '+$Token
  'Content-Type'='application/json'
  'X-OpenRouter-Metadata'='enabled'
}

$Body=[ordered]@{
  model='nvidia/nemotron-3-ultra-550b-a55b:free'
  temperature=0
  messages=@(
    @{role='system';content='You are a protocol test. Use the provided tool exactly once when asked.'},
    @{role='user';content='Call the read_file tool with path smoke.txt. Do not answer normally.'}
  )
  tools=@(
    @{
      type='function'
      function=@{
        name='read_file'
        description='Read a file.'
        parameters=@{
          type='object'
          properties=@{path=@{type='string'}}
          required=@('path')
          additionalProperties=$false
        }
      }
    }
  )
  tool_choice='auto'
  max_tokens=512
}
$Json=$Body|ConvertTo-Json -Depth 20 -Compress
$Resp=Invoke-RestMethod -Method Post -Uri 'https://openrouter.ai/api/v1/chat/completions' -Headers $Headers -Body $Json -TimeoutSec 120
$Msg=$Resp.choices[0].message
$Calls=@($Msg.tool_calls)
Write-Host "response_model=$($Resp.model)"
Write-Host "finish_reason=$($Resp.choices[0].finish_reason)"
Write-Host "tool_call_count=$($Calls.Count)"
if($Calls.Count-ne1){throw "NEMOTRON_TOOL_CALL_COUNT_INVALID count=$($Calls.Count)"}
if([string]$Calls[0].function.name-cne'read_file'){throw "NEMOTRON_TOOL_NAME_INVALID name=$($Calls[0].function.name)"}
$Args=[string]$Calls[0].function.arguments|ConvertFrom-Json
if([string]$Args.path-cne'smoke.txt'){throw "NEMOTRON_TOOL_ARG_INVALID path=$($Args.path)"}
Write-Host 'NEMOTRON_TOOL_CALLING=PASS'
