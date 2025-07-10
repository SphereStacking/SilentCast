; SilentCast AutoHotkey Script
; プレフィックスキー: Alt+S+C

; Alt+S+Cの後にキーを押すための設定
#NoEnv
SendMode Input
SetWorkingDir %A_ScriptDir%

; プレフィックスモードのフラグ
PrefixMode := false

; Alt+S+Cでプレフィックスモードに入る
!s::
If (A_PriorHotkey = "!s" and A_TimeSincePriorHotkey < 500)
{
    ; Alt+S+Cが完成
    PrefixMode := true
    ToolTip, SilentCast: Ready for command...
    SetTimer, ResetPrefix, 1000  ; 1秒後にリセット
}
return

!c::
If (A_PriorHotkey = "!s" and A_TimeSincePriorHotkey < 100)
{
    ; Alt+S+Cが完成
    PrefixMode := true
    ToolTip, SilentCast: Ready for command...
    SetTimer, ResetPrefix, 1000  ; 1秒後にリセット
}
return

; プレフィックスモードのリセット
ResetPrefix:
PrefixMode := false
ToolTip
return

; プレフィックスモード中のホットキー
#If PrefixMode

; N - メモ帳
n::
Run, notepad.exe
PrefixMode := false
ToolTip
return

; E - エクスプローラー
e::
Run, explorer.exe
PrefixMode := false
ToolTip
return

; C - コマンドプロンプト
c::
Run, cmd.exe
PrefixMode := false
ToolTip
return

; T - Windows Terminal
t::
Run, wt.exe
PrefixMode := false
ToolTip
return

; X - 電卓
x::
Run, calc.exe
PrefixMode := false
ToolTip
return

; B - ブラウザ（Edge）
b::
Run, msedge.exe
PrefixMode := false
ToolTip
return

; P - PowerShell
p::
Run, powershell.exe
PrefixMode := false
ToolTip
return

; V - VS Code
v::
Run, code
PrefixMode := false
ToolTip
return

; H - Hello メッセージ
h::
MsgBox, 0, SilentCast, Hello from SilentCast! 🪄
PrefixMode := false
ToolTip
return

; G,S - Git Status（例：Windows Terminalで実行）
g::
Input, NextKey, L1 T1
if (NextKey = "s")
{
    Run, wt.exe -d . -- git status
    PrefixMode := false
    ToolTip
}
return

#If

; Ctrl+Q で終了
^q::ExitApp