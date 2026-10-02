param(
    [Parameter(Mandatory=$true)][string]$RepoPath,
    [Parameter(Mandatory=$true)][ValidateSet('Wingless','Yggdrasil')][string]$Program,
    [Parameter(Mandatory=$true)][string]$ActiveBranch,
    [Parameter(Mandatory=$true)][string]$StartSha,
    [Parameter(Mandatory=$true)][string]$PromptPath,
    [Parameter(Mandatory=$true)][string]$LastMessagePath,
    [ValidateSet('openrouter','nvidia')][string]$Provider='openrouter',
    [string]$Model='',
    [int]$MaxRounds=48,
    [switch]$ProtocolProbe,
    [switch]$SelfTest
)

Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'

if([string]::IsNullOrWhiteSpace($Model)){
    $Model=if($Provider-ceq'nvidia'){'moonshotai/kimi-k3'}else{'nvidia/nemotron-3-ultra-550b-a55b:free'}
}

function Resolve-RepoPath([string]$Relative,[bool]$MustExist=$false){
    if([string]::IsNullOrWhiteSpace($Relative)){throw 'NEMOTRON_PATH_EMPTY'}
    if([IO.Path]::IsPathRooted($Relative)){throw "NEMOTRON_PATH_ROOTED path=$Relative"}
    $Root=[IO.Path]::GetFullPath($RepoPath).TrimEnd('\')
    $Full=[IO.Path]::GetFullPath((Join-Path $Root ($Relative -replace '/','\')))
    $InsideRoot=$Full.Equals($Root,[StringComparison]::OrdinalIgnoreCase) -or $Full.StartsWith($Root+'\',[StringComparison]::OrdinalIgnoreCase)
    if(-not$InsideRoot){throw "NEMOTRON_PATH_ESCAPE path=$Relative"}
    if($MustExist -and -not(Test-Path -LiteralPath $Full)){throw "NEMOTRON_PATH_MISSING path=$Relative"}
    return $Full
}

function Get-Relative([string]$Full){
    $Root=[IO.Path]::GetFullPath($RepoPath).TrimEnd('\')+'\'
    $F=[IO.Path]::GetFullPath($Full)
    if(-not$F.StartsWith($Root,[StringComparison]::OrdinalIgnoreCase)){throw 'NEMOTRON_RELATIVE_PATH_ESCAPE'}
    return $F.Substring($Root.Length).Replace('\','/')
}

function Test-WritablePath([string]$Relative){
    $P=$Relative.Replace('\','/')
    if($P.Contains('..')){return $false}
    if($Program-ceq'Wingless'){
        return $P-ceq'.wingless/qualification-request.json' -or
               $P.StartsWith('docs/experiments/') -or
               $P.StartsWith('unitary/') -or
               $P.StartsWith('cmd/') -or
               $P.StartsWith('scripts/')
    }
    return $P-ceq'.yggdrasil/qualification-request.json' -or
           $P-ceq'.yggdrasil/isolated-run.json' -or
           $P.StartsWith('research/experiments/') -or
           $P.StartsWith('research/applications/plane/')
}

function Assert-Writable([string]$Relative){
    if(-not(Test-WritablePath $Relative)){throw "NEMOTRON_WRITE_PATH_FORBIDDEN path=$Relative"}
    [void](Resolve-RepoPath $Relative $false)
}

function Invoke-GitRaw([string[]]$GitArgs){
    $PriorErrorAction=$ErrorActionPreference
    try{
        # Native git may emit non-fatal warnings on stderr (for example line-ending
        # normalization) while still returning exit code 0. Capture those bytes as
        # tool output and judge success strictly by the process exit code.
        $ErrorActionPreference='Continue'
        $Out=@(& git -C $RepoPath @GitArgs 2>&1 | ForEach-Object {[string]$_})
        $Code=$LASTEXITCODE
    }
    finally{
        $ErrorActionPreference=$PriorErrorAction
    }
    return [ordered]@{exit_code=$Code;output=($Out -join [Environment]::NewLine)}
}

function Convert-ToolResult($Object){
    if($Object -is [string]){return $Object}
    return ($Object|ConvertTo-Json -Depth 30 -Compress)
}

function Invoke-Tool([string]$Name,[string]$ArgumentsJson){
    $A=if([string]::IsNullOrWhiteSpace($ArgumentsJson)){[pscustomobject]@{}}else{$ArgumentsJson|ConvertFrom-Json}
    switch($Name){
        'list_files' {
            $Rel=if($A.PSObject.Properties.Name -contains 'path'){[string]$A.path}else{'.'}
            $Depth=if($A.PSObject.Properties.Name -contains 'max_depth'){[Math]::Min([Math]::Max([int]$A.max_depth,0),6)}else{3}
            $Base=Resolve-RepoPath $Rel $true
            if(-not(Test-Path -LiteralPath $Base -PathType Container)){throw "NEMOTRON_LIST_NOT_DIRECTORY path=$Rel"}
            $BaseDepth=($Base.TrimEnd('\').Split('\')).Count
            $Rows=Get-ChildItem -LiteralPath $Base -Force -Recurse -ErrorAction SilentlyContinue |
                Where-Object {(($_.FullName.Split('\')).Count-$BaseDepth)-le$Depth} |
                Select-Object -First 600 |
                ForEach-Object {
                    [ordered]@{path=(Get-Relative $_.FullName);type=if($_.PSIsContainer){'dir'}else{'file'};bytes=if($_.PSIsContainer){0}else{$_.Length}}
                }
            return Convert-ToolResult @($Rows)
        }
        'read_file' {
            $Rel=[string]$A.path
            $Full=Resolve-RepoPath $Rel $true
            if(-not(Test-Path -LiteralPath $Full -PathType Leaf)){throw "NEMOTRON_READ_NOT_FILE path=$Rel"}
            $Lines=@(Get-Content -LiteralPath $Full)
            $Start=if($A.PSObject.Properties.Name -contains 'start_line'){[Math]::Max([int]$A.start_line,1)}else{1}
            $End=if($A.PSObject.Properties.Name -contains 'end_line'){[Math]::Min([int]$A.end_line,$Lines.Count)}else{[Math]::Min($Start+399,$Lines.Count)}
            if($Start-gt$Lines.Count){return ''}
            $Text=($Lines[($Start-1)..($End-1)] -join [Environment]::NewLine)
            if($Text.Length-gt65536){$Text=$Text.Substring(0,65536)}
            return $Text
        }
        'search_text' {
            $Query=[string]$A.query
            if([string]::IsNullOrWhiteSpace($Query)){throw 'NEMOTRON_SEARCH_QUERY_EMPTY'}
            $Rel=if($A.PSObject.Properties.Name -contains 'path'){[string]$A.path}else{'.'}
            $Base=Resolve-RepoPath $Rel $true
            $Max=if($A.PSObject.Properties.Name -contains 'max_results'){[Math]::Min([Math]::Max([int]$A.max_results,1),100)}else{50}
            $Files=if(Test-Path -LiteralPath $Base -PathType Leaf){@((Get-Item -LiteralPath $Base))}else{
                @(Get-ChildItem -LiteralPath $Base -File -Recurse -ErrorAction SilentlyContinue | Where-Object{$_.Length-le2097152})
            }
            $Rows=New-Object Collections.Generic.List[object]
            foreach($F in $Files){
                try{
                    foreach($M in @(Select-String -LiteralPath $F.FullName -SimpleMatch -Pattern $Query -ErrorAction Stop)){
                        $Rows.Add([ordered]@{path=(Get-Relative $F.FullName);line=$M.LineNumber;text=([string]$M.Line).Trim()})
                        if($Rows.Count-ge$Max){break}
                    }
                }catch{}
                if($Rows.Count-ge$Max){break}
            }
            return Convert-ToolResult @($Rows)
        }
        'git_status' {
            return Convert-ToolResult (Invoke-GitRaw @('status','--short','--branch'))
        }
        'git_log' {
            $N=if($A.PSObject.Properties.Name -contains 'max_count'){[Math]::Min([Math]::Max([int]$A.max_count,1),40)}else{12}
            return Convert-ToolResult (Invoke-GitRaw @('log',("-"+$N),'--oneline','--decorate'))
        }
        'git_diff' {
            $Ref=if($A.PSObject.Properties.Name -contains 'ref'){[string]$A.ref}else{''}
            if($Ref -and $Ref -notmatch '^[0-9A-Za-z._/^~-]+$'){throw "NEMOTRON_GIT_REF_INVALID ref=$Ref"}
            $GitDiffArgs=New-Object Collections.Generic.List[string]
            $GitDiffArgs.Add('diff')
            if($Ref){$GitDiffArgs.Add($Ref)}
            $R=Invoke-GitRaw @($GitDiffArgs.ToArray())
            if($R.output.Length-gt65536){$R.output=$R.output.Substring(0,65536)}
            return Convert-ToolResult $R
        }
        'write_file' {
            $Rel=[string]$A.path
            Assert-Writable $Rel
            $Content=[string]$A.content
            if($Content.Length-gt1048576){throw "NEMOTRON_WRITE_TOO_LARGE path=$Rel chars=$($Content.Length)"}
            $Full=Resolve-RepoPath $Rel $false
            $Dir=Split-Path -Parent $Full
            if($Dir){New-Item -ItemType Directory -Force -Path $Dir|Out-Null}
            [IO.File]::WriteAllText($Full,$Content,(New-Object Text.UTF8Encoding($false)))
            return "WROTE $Rel chars=$($Content.Length)"
        }
        'delete_file' {
            $Rel=[string]$A.path
            Assert-Writable $Rel
            $Full=Resolve-RepoPath $Rel $false
            if(Test-Path -LiteralPath $Full -PathType Container){throw "NEMOTRON_DELETE_DIRECTORY_FORBIDDEN path=$Rel"}
            if(Test-Path -LiteralPath $Full -PathType Leaf){Remove-Item -LiteralPath $Full -Force}
            return "DELETED $Rel"
        }
        'git_restore' {
            $Paths=@($A.paths|ForEach-Object{[string]$_})
            if($Paths.Count-eq0){throw 'NEMOTRON_RESTORE_PATHS_EMPTY'}
            foreach($P in $Paths){Assert-Writable $P}
            $R=Invoke-GitRaw (@('restore','--')+$Paths)
            if($R.exit_code-ne0){throw "NEMOTRON_GIT_RESTORE_FAILED output=$($R.output)"}
            return 'RESTORED '+($Paths -join ',')
        }
        'git_commit' {
            $Message=[string]$A.message
            if([string]::IsNullOrWhiteSpace($Message) -or $Message.Length-gt200){throw 'NEMOTRON_COMMIT_MESSAGE_INVALID'}
            $Paths=@($A.paths|ForEach-Object{[string]$_})
            if($Paths.Count-eq0){throw 'NEMOTRON_COMMIT_PATHS_EMPTY'}
            foreach($P in $Paths){Assert-Writable $P}
            [void](Invoke-GitRaw @('reset'))
            $Add=Invoke-GitRaw (@('add','--')+$Paths)
            if($Add.exit_code-ne0){throw "NEMOTRON_GIT_ADD_FAILED output=$($Add.output)"}
            $Staged=Invoke-GitRaw @('diff','--cached','--name-only')
            if($Staged.exit_code-ne0 -or [string]::IsNullOrWhiteSpace($Staged.output)){throw 'NEMOTRON_COMMIT_NOTHING_STAGED'}
            foreach($Line in @([regex]::Split($Staged.output,'\r?\n'))){
                if([string]::IsNullOrWhiteSpace($Line)){continue}
                if(-not(Test-WritablePath $Line)){throw "NEMOTRON_STAGED_PATH_FORBIDDEN path=$Line"}
            }
            $Commit=Invoke-GitRaw @('commit','-m',$Message)
            if($Commit.exit_code-ne0){throw "NEMOTRON_GIT_COMMIT_FAILED output=$($Commit.output)"}
            $Head=(Invoke-GitRaw @('rev-parse','HEAD'))
            return Convert-ToolResult ([ordered]@{committed=$true;head=$Head.output.Trim();files=$Staged.output})
        }
        'run_repo_process' {
            $Kind=[string]$A.kind
            [string[]]$ProcessArgs=@($A.args|ForEach-Object{[string]$_})
            $Timeout=if($A.PSObject.Properties.Name -contains 'timeout_seconds'){[Math]::Min([Math]::Max([int]$A.timeout_seconds,1),600)}else{180}
            $Exe='';$FinalArgs=@()
            switch($Kind){
                'go_test' {
                    $Exe=(Get-Command go -ErrorAction Stop).Source
                    $FinalArgs=@('test')+$ProcessArgs
                    foreach($X in $ProcessArgs){if($X -match '^-?(exec|toolexec|overlay|vettool)(=|$)'){throw "NEMOTRON_GO_FLAG_FORBIDDEN arg=$X"}}
                }
                'go_run' {
                    $Exe=(Get-Command go -ErrorAction Stop).Source
                    if($ProcessArgs.Count-eq0 -or $ProcessArgs[0] -notmatch '^(\./)?cmd/'){throw 'NEMOTRON_GO_RUN_PATH_INVALID'}
                    $FinalArgs=@('run')+$ProcessArgs
                }
                'powershell_file' {
                    if($ProcessArgs.Count-eq0){throw 'NEMOTRON_POWERSHELL_FILE_REQUIRED'}
                    $ScriptRel=$ProcessArgs[0]
                    if($ScriptRel -notmatch '^scripts/[A-Za-z0-9._/-]+\.ps1$'){throw "NEMOTRON_POWERSHELL_PATH_INVALID path=$ScriptRel"}
                    $ScriptFull=Resolve-RepoPath $ScriptRel $true
                    $Exe='powershell.exe'
                    $Tail=if($ProcessArgs.Count-gt1){@($ProcessArgs[1..($ProcessArgs.Count-1)])}else{@()}
                    $FinalArgs=@('-NoProfile','-NonInteractive','-ExecutionPolicy','Bypass','-File',$ScriptFull)+$Tail
                }
                'python_file' {
                    if($Program-cne'Yggdrasil'){throw 'NEMOTRON_PYTHON_NOT_ALLOWED_FOR_PROGRAM'}
                    if($ProcessArgs.Count-eq0){throw 'NEMOTRON_PYTHON_FILE_REQUIRED'}
                    $ScriptRel=$ProcessArgs[0]
                    if($ScriptRel -notmatch '^research/applications/plane/[A-Za-z0-9._/-]+\.py$'){throw "NEMOTRON_PYTHON_PATH_INVALID path=$ScriptRel"}
                    $Exe='C:\ProgramData\CKBR\research-sidecar-yggdrasil\python312\python.exe'
                    if(-not(Test-Path -LiteralPath $Exe -PathType Leaf)){$Exe=(Get-Command python -ErrorAction Stop).Source}
                    $Tail=if($ProcessArgs.Count-gt1){@($ProcessArgs[1..($ProcessArgs.Count-1)])}else{@()}
                    $FinalArgs=@((Resolve-RepoPath $ScriptRel $true))+$Tail
                }
                'git_diff_check' {
                    $Exe=(Get-Command git -ErrorAction Stop).Source
                    $FinalArgs=@('-C',$RepoPath,'diff','--check')
                }
                default {throw "NEMOTRON_PROCESS_KIND_INVALID kind=$Kind"}
            }
            $OutFile=Join-Path $env:RUNNER_TEMP ("nemotron-out-"+[guid]::NewGuid().ToString('N')+".txt")
            $ErrFile=Join-Path $env:RUNNER_TEMP ("nemotron-err-"+[guid]::NewGuid().ToString('N')+".txt")
            try{
                $P=Start-Process -FilePath $Exe -ArgumentList $FinalArgs -WorkingDirectory $RepoPath -RedirectStandardOutput $OutFile -RedirectStandardError $ErrFile -NoNewWindow -PassThru
                if(-not$P.WaitForExit($Timeout*1000)){
                    try{$P.Kill()}catch{}
                    throw "NEMOTRON_PROCESS_TIMEOUT kind=$Kind timeout_seconds=$Timeout"
                }
                $Stdout=if(Test-Path -LiteralPath $OutFile){[IO.File]::ReadAllText($OutFile)}else{''}
                $Stderr=if(Test-Path -LiteralPath $ErrFile){[IO.File]::ReadAllText($ErrFile)}else{''}
                if($Stdout.Length-gt50000){$Stdout=$Stdout.Substring(0,50000)}
                if($Stderr.Length-gt20000){$Stderr=$Stderr.Substring(0,20000)}
                return Convert-ToolResult ([ordered]@{exit_code=$P.ExitCode;stdout=$Stdout;stderr=$Stderr})
            }finally{
                Remove-Item -LiteralPath $OutFile,$ErrFile -Force -ErrorAction SilentlyContinue
            }
        }
        default {throw "NEMOTRON_TOOL_UNKNOWN name=$Name"}
    }
}

if(-not(Test-Path -LiteralPath $RepoPath -PathType Container)){throw 'NEMOTRON_REPO_MISSING'}
if(-not(Test-Path -LiteralPath $PromptPath -PathType Leaf)){throw 'NEMOTRON_PROMPT_MISSING'}
$CurrentBranch=(& git -C $RepoPath branch --show-current 2>$null|Select-Object -First 1)
if(([string]$CurrentBranch).Trim()-cne$ActiveBranch){throw "NEMOTRON_BRANCH_MISMATCH expected=$ActiveBranch actual=$CurrentBranch"}
$CurrentHead=(& git -C $RepoPath rev-parse HEAD).Trim()
if($CurrentHead-cne$StartSha){throw "NEMOTRON_START_SHA_MISMATCH expected=$StartSha actual=$CurrentHead"}

if($SelfTest){
    $Status=(Invoke-Tool 'git_status' '{}')|ConvertFrom-Json
    if([int]$Status.exit_code-ne0){throw "NEMOTRON_SELFTEST_GIT_STATUS_FAILED output=$($Status.output)"}
    $Log=(Invoke-Tool 'git_log' '{"max_count":2}')|ConvertFrom-Json
    if([int]$Log.exit_code-ne0){throw "NEMOTRON_SELFTEST_GIT_LOG_FAILED output=$($Log.output)"}
    $Listing=Invoke-Tool 'list_files' '{"path":".","max_depth":1}'
    if([string]::IsNullOrWhiteSpace([string]$Listing)){throw 'NEMOTRON_SELFTEST_LIST_EMPTY'}
    Write-Host 'NEMOTRON_AGENT_SELFTEST=PASS'
    return
}

$Endpoint=$null
if($Provider-ceq'nvidia'){
    $Token=[string]$env:NVIDIA_API_KEY
    if([string]::IsNullOrWhiteSpace($Token)){throw 'NEMOTRON_NVIDIA_SECRET_MISSING'}
    $Endpoint='https://integrate.api.nvidia.com/v1/chat/completions'
}else{
    $SecretCandidates=@(
        'C:\ProgramData\CKBR\research-sidecar-yggdrasil\secrets\openrouter.dpapi',
        'C:\ProgramData\CKBR\research-sidecar\secrets\openrouter.dpapi'
    )
    $Secret=$SecretCandidates|Where-Object{Test-Path -LiteralPath $_ -PathType Leaf}|Select-Object -First 1
    if([string]::IsNullOrWhiteSpace($Secret)){throw 'NEMOTRON_OPENROUTER_SECRET_MISSING'}
    $Secure=(Get-Content -LiteralPath $Secret -Raw).Trim()|ConvertTo-SecureString
    $Ptr=[Runtime.InteropServices.Marshal]::SecureStringToBSTR($Secure)
    try{$Token=[Runtime.InteropServices.Marshal]::PtrToStringBSTR($Ptr)}finally{[Runtime.InteropServices.Marshal]::ZeroFreeBSTR($Ptr)}
    if([string]::IsNullOrWhiteSpace($Token)){throw 'NEMOTRON_OPENROUTER_SECRET_EMPTY'}
    $Endpoint='https://openrouter.ai/api/v1/chat/completions'
}

$Prompt=[IO.File]::ReadAllText($PromptPath)
$System=@"
You are the repo-local autonomous research worker. You have a constrained tool interface to inspect and modify exactly one research repository.
Use the minimum tools necessary and never inspect unrelated files just to gather more context. Never attempt to access paths outside the repository or authority outside the prompt.
Follow the scientific preregistration/freeze/no-post-result-tuning rules exactly.
Make local Git commits using git_commit. Never push.
You MUST terminate by calling the finish tool exactly once when the task is complete or legitimately blocked. Do not keep exploring after the prompt's requirements are satisfied. Do not send a normal final answer instead of finish.
"@

$Tools=@(
 @{type='function';function=@{name='list_files';description='List files/directories under a repository-relative directory.';parameters=@{type='object';properties=@{path=@{type='string'};max_depth=@{type='integer';minimum=0;maximum=6}};additionalProperties=$false}}},
 @{type='function';function=@{name='read_file';description='Read a bounded line range from a repository file.';parameters=@{type='object';properties=@{path=@{type='string'};start_line=@{type='integer';minimum=1};end_line=@{type='integer';minimum=1}};required=@('path');additionalProperties=$false}}},
 @{type='function';function=@{name='search_text';description='Search literal text in repository files.';parameters=@{type='object';properties=@{query=@{type='string'};path=@{type='string'};max_results=@{type='integer';minimum=1;maximum=100}};required=@('query');additionalProperties=$false}}},
 @{type='function';function=@{name='git_status';description='Return git status for the research checkout.';parameters=@{type='object';properties=@{};additionalProperties=$false}}},
 @{type='function';function=@{name='git_log';description='Return recent commit log.';parameters=@{type='object';properties=@{max_count=@{type='integer';minimum=1;maximum=40}};additionalProperties=$false}}},
 @{type='function';function=@{name='git_diff';description='Return repository diff, optionally against a safe git ref.';parameters=@{type='object';properties=@{ref=@{type='string'}};additionalProperties=$false}}},
 @{type='function';function=@{name='write_file';description='Write a UTF-8 file only within the program allowed research paths.';parameters=@{type='object';properties=@{path=@{type='string'};content=@{type='string'}};required=@('path','content');additionalProperties=$false}}},
 @{type='function';function=@{name='delete_file';description='Delete a file only within allowed research paths.';parameters=@{type='object';properties=@{path=@{type='string'}};required=@('path');additionalProperties=$false}}},
 @{type='function';function=@{name='git_restore';description='Restore allowed research paths from HEAD.';parameters=@{type='object';properties=@{paths=@{type='array';items=@{type='string'};minItems=1}};required=@('paths');additionalProperties=$false}}},
 @{type='function';function=@{name='git_commit';description='Stage exactly the supplied allowed research paths and commit them locally.';parameters=@{type='object';properties=@{message=@{type='string'};paths=@{type='array';items=@{type='string'};minItems=1}};required=@('message','paths');additionalProperties=$false}}},
 @{type='function';function=@{name='run_repo_process';description='Run a bounded project-local verification process. kind is one of go_test, go_run, powershell_file, python_file, git_diff_check.';parameters=@{type='object';properties=@{kind=@{type='string';enum=@('go_test','go_run','powershell_file','python_file','git_diff_check')};args=@{type='array';items=@{type='string'}};timeout_seconds=@{type='integer';minimum=1;maximum=600}};required=@('kind','args');additionalProperties=$false}}},
 @{type='function';function=@{name='finish';description='Finish the bounded research task. Call exactly once when complete or legitimately blocked. summary must include any required terminal marker from the prompt.';parameters=@{type='object';properties=@{summary=@{type='string';minLength=1};status=@{type='string';enum=@('complete','blocked')}};required=@('summary','status');additionalProperties=$false}}}
)
if($ProtocolProbe){
    $Tools=@($Tools|Where-Object{[string]$_.function.name -in @('git_status','finish')})
}

$Messages=New-Object Collections.Generic.List[object]
$Messages.Add([ordered]@{role='system';content=$System})
$Messages.Add([ordered]@{role='user';content=$Prompt})

$Headers=@{
    Authorization='Bearer '+$Token
    'Content-Type'='application/json'
}
if($Provider-ceq'openrouter'){$Headers['X-OpenRouter-Metadata']='enabled'}

for($Round=1;$Round-le$MaxRounds;$Round++){
    if($Round-eq([Math]::Max(2,$MaxRounds-2))){
        $Messages.Add([ordered]@{role='user';content='The tool budget is nearly exhausted. Stop broad exploration. Complete only the essential remaining work, then call finish. If the task cannot be completed under the frozen rules, call finish with status blocked and explain why.'})
    }
    $ToolChoice=if($Round-eq$MaxRounds){
        [ordered]@{type='function';function=[ordered]@{name='finish'}}
    }else{'auto'}
    $BodyObject=[ordered]@{
        model=$Model
        temperature=if($Provider-ceq'nvidia'){1}else{0}
        messages=$Messages.ToArray()
        tools=$Tools
        tool_choice=$ToolChoice
        max_tokens=if($Provider-ceq'nvidia'){16384}else{4096}
    }
    if($Provider-ceq'nvidia'){$BodyObject.reasoning_effort='low'}
    $Body=$BodyObject|ConvertTo-Json -Depth 50 -Compress

    $Resp=$null
    for($Attempt=1;$Attempt-le3;$Attempt++){
        try{
            $Resp=Invoke-RestMethod -Method Post -Uri $Endpoint -Headers $Headers -Body $Body -TimeoutSec 180
            break
        }catch{
            $StatusCode=$null
            try{$StatusCode=[int]$_.Exception.Response.StatusCode}catch{}
            $IsRateLimit=($StatusCode-eq429 -or $_.Exception.Message -match '(?i)429|too many requests')
            if($IsRateLimit -and $Attempt-lt3){
                $Delay=if($Attempt-eq1){5}else{15}
                Write-Host "NEMOTRON_PROVIDER_RATE_LIMIT provider=$Provider round=$Round attempt=$Attempt retry_seconds=$Delay"
                Start-Sleep -Seconds $Delay
                continue
            }
            throw "NEMOTRON_PROVIDER_REQUEST_FAILED provider=$Provider round=$Round attempt=$Attempt status=$StatusCode error=$($_.Exception.Message)"
        }
    }
    if($null-eq$Resp){throw "NEMOTRON_PROVIDER_RESPONSE_MISSING provider=$Provider round=$Round"}
    if($null-eq$Resp.choices -or @($Resp.choices).Count-lt1){throw "NEMOTRON_RESPONSE_CHOICES_MISSING round=$Round"}
    $Msg=$Resp.choices[0].message
    $ToolCallsProperty=$Msg.PSObject.Properties['tool_calls']
    [object[]]$Calls=@()
    if($null-ne$ToolCallsProperty -and $null-ne$ToolCallsProperty.Value){
        $Calls=@($ToolCallsProperty.Value)
    }
    $ContentProperty=$Msg.PSObject.Properties['content']
    $ReasoningProperty=$Msg.PSObject.Properties['reasoning_content']
    $Assistant=[ordered]@{role='assistant';content=if($null-eq$ContentProperty -or $null-eq$ContentProperty.Value){''}else{[string]$ContentProperty.Value}}
    if($null-ne$ReasoningProperty -and $null-ne$ReasoningProperty.Value){
        $Assistant.reasoning_content=[string]$ReasoningProperty.Value
    }
    if($Calls.Count-gt0){
        $ToolCallRows=New-Object Collections.Generic.List[object]
        foreach($Call in $Calls){
            $ToolCallRows.Add([ordered]@{
                id=[string]$Call.id
                type='function'
                function=[ordered]@{
                    name=[string]$Call.function.name
                    arguments=[string]$Call.function.arguments
                }
            })
        }
        $Assistant.tool_calls=$ToolCallRows.ToArray()
    }
    $Messages.Add($Assistant)

    if($Calls.Count-eq0){
        $Final=[string]$Assistant.content
        $Reasoning=if($Assistant.Contains('reasoning_content')){[string]$Assistant.reasoning_content}else{''}
        if([string]::IsNullOrWhiteSpace($Final)){
            if($Provider-ceq'nvidia' -and -not[string]::IsNullOrWhiteSpace($Reasoning) -and $Round-lt$MaxRounds){
                $Messages.Add([ordered]@{role='user';content='Continue from your preserved reasoning. You must now call one of the supplied tools; call finish when complete.'})
                continue
            }
            throw "NEMOTRON_EMPTY_FINAL round=$Round"
        }
        if($ProtocolProbe){throw "NEMOTRON_PROTOCOL_TOOL_CALL_REQUIRED round=$Round"}
        [IO.File]::WriteAllText($LastMessagePath,$Final,(New-Object Text.UTF8Encoding($false)))
        Write-Host "NEMOTRON_AGENT_PASS provider=$Provider rounds=$Round model=$Model"
        return
    }

    foreach($Call in $Calls){
        $Name=[string]$Call.function.name
        Write-Host "NEMOTRON_TOOL round=$Round name=$Name"
        if($Name-ceq'finish'){
            $FinishArgs=([string]$Call.function.arguments)|ConvertFrom-Json
            $Summary=[string]$FinishArgs.summary
            if([string]::IsNullOrWhiteSpace($Summary)){throw 'NEMOTRON_FINISH_SUMMARY_EMPTY'}
            [IO.File]::WriteAllText($LastMessagePath,$Summary,(New-Object Text.UTF8Encoding($false)))
            Write-Host "NEMOTRON_AGENT_PASS provider=$Provider rounds=$Round model=$Model status=$([string]$FinishArgs.status)"
            return
        }
        try{
            $Result=Invoke-Tool $Name ([string]$Call.function.arguments)
            $ToolContent=[string]$Result
        }catch{
            $ToolContent="TOOL_ERROR: $($_.Exception.Message)"
        }
        if($ToolContent.Length-gt70000){$ToolContent=$ToolContent.Substring(0,70000)}
        $Messages.Add([ordered]@{role='tool';tool_call_id=[string]$Call.id;content=$ToolContent})
    }
}
throw "NEMOTRON_AGENT_MAX_ROUNDS_EXCEEDED max_rounds=$MaxRounds"
