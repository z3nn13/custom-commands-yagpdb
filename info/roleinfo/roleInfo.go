{{/*
        
        Recommended Trigger: \A(?:\-|<@!?204255221017214977>)\s*(?:r(?:ole)?i(?:nfo)?)(?: +|\z)
        Trigger Type: Regex
        Usage: -roleinfo <Role:ID/Mention/Name/Position> [-p]
        Aliases: ri, rolei, rinfo

        Copyright (c): zen | ゼン#0008; 2021
        License: MIT
        Repository: https://github.com/z3nn13/custom-commands-yagpdb
        
*/}}

{{/* Global Variables */}}
{{$help := cembed "title" "RoleInfo/RInfo/Ri" "description" (print "```Roleinfo [Role:Name/Mention/ID/Position]``````[-p p:Switch - Include Permissions]```Shows Information about a role")}}
{{$guildRoles := cslice.AppendSlice .Guild.Roles}} {{$listroles := split (exec "listroles" ) "\n"}}
{{$role := ""}}{{$errorMsg := ""}}{{$err := false}}{{$StrippedMsg := .StrippedMsg}}{{$permFlag := false}}

{{/* Checking Input */}}
{{if .CmdArgs}}
        {{if reFind `\s+\-(?:p(?:erm(?:ission)?)?s?)(?:\s+|\z)` $StrippedMsg}}
                {{$StrippedMsg = reReplace `\s+\-(?:p(?:erm(?:ission)?)?s?)(?:\s+|\z)` $StrippedMsg ""}}
                {{$permFlag = true}}
        {{end}}
        {{ if $roleID := reFind `\d{17,19}` $StrippedMsg|toInt64}}{{/* Mention or ID */}}
                {{ with .Guild.GetRole $roleID}}
                        {{ $role = .}}
                {{ else }}
                        {{ $errorMsg = "Invalid ID/Role Does Not Exist"}}{{$err = true}}
                {{ end }}
        {{else}}
                {{$found := false}}
                {{/* Name Input */}}
                {{range $guildRoles}}
                        {{- if eq (.Name|lower) ($StrippedMsg|lower) (index $.CmdArgs 0|lower) -}}
                                {{- $role = .}}{{$found = true -}}
                        {{- end -}}
                {{- end}}
                {{/* Position Input */}}
                {{ if $index := toInt $StrippedMsg}}
                        {{if le $index (len $guildRoles)}}
                                {{ $role = add $index 1 | index $listroles | reFind `\d{17,19}` | toInt64 | .Guild.GetRole}}{{$found = true}}
                        {{end}}
                {{end}}
                {{if not $found}}
                        {{$errorMsg =printf "%q not recognized: Invalid Name/Position" $StrippedMsg}}{{$err = true}}
                {{end}}
        {{ end }}
{{else}}
        {{$err = true}}
{{end}}

{{/* Preparing Embed */}}
{{ $ex := or (and (reFind "a_" .Guild.Icon) "gif" ) "png" }}
{{ $icon := print "https://cdn.discordapp.com/icons/" .Guild.ID "/" .Guild.Icon "." $ex "?size=1024" }}
{{ $embed := sdict "author"
( sdict  "name" "Role Info" "icon_url" "https://images-ext-2.discordapp.net/external/G67VOLJZEh_p_JpowOPIRo4LimCe4KNMj7X5Azffufc/https/cdn.discordapp.com/emojis/764251327814696970.png")
"thumbnail" (sdict "url" $icon)}}

{{with $role}}
        {{ $createdAt := div .ID 4194304 | add 1420070400000 | mult 1000000 | toDuration | (newDate 1970 1 1 0 0 0).Add }}
        {{ $mentionable := or (and .Mentionable "`Yes`") "`No`" }}
        {{ $hoist := or (and .Hoist "`Yes`") "`No`" }}
        {{ $managed := or (and .Managed "`Yes`") "`No`" }}

        {{ $pos := 0 }}{{ $up := "" }}{{ $down := "" }}
        {{ range $i,$v := $listroles }}{{ if reFind (str $role.ID) $v }}{{ $pos = sub $i 2 }}{{ end }}{{ end }}
        
        {{/* .Role.Position had a wonky order, so this was a workaround */}}
        {{ if eq $pos 0 }}
                {{ $up = "-----------\n" }}
        {{ else }}
                {{ $up = printf "> #%d • %s\n" $pos ((sub $pos 1 |index $guildRoles).ID|mentionRoleID) }}
        {{ end }}
        
        {{ if eq $pos (len $guildRoles|add -1) }}
                {{ $down = "-----------\n"}}
        {{ else }}
                {{ $down = printf "> #%d • %s\n" (add $pos 2) ((add $pos 1|index $guildRoles).ID|mentionRoleID) }}
        {{ end }}
        {{ $final_pos := printf "%s> **#%d • %s**\n%s\n> `.Position` = %d\n> (Total Roles: **%d**)" $up (add $pos 1) (.ID|mentionRoleID) $down .Position (len $guildRoles)}}


        {{ $fields := cslice
        ( sdict "name" "• Name" "value" (.ID|mentionRoleID) "inline" true)
        ( sdict "name" "• ID" "value" (str .ID) "inline" true)
        ( sdict "name" "• Others" "value" (printf "> **Hoist** • %s\n> **Managed** • %s\n> **Mentionable** • %s" $hoist $managed $mentionable) "inline" true)
        ( sdict "name" "• Position ↓" "value" $final_pos "inline" true)
        ( sdict "name" "• Color" "value" (printf "#%x" .Color|upper) "inline" true)
        ( sdict "name" "• Created At" "value" (print ($createdAt.Format "📆 January 2, 2006\n") "️️️⏱️ " (currentTime.Sub $createdAt|humanizeDurationSeconds) " ago") "inline" true)}}
        {{$embed.Set "footer" (sdict "text" (print "Triggered By • " $.User.String " • Use `-p` flag to view perms ") "icon_url" ($.User.AvatarURL "256"))}}

        {{/* Credits To Satty#9361 for this bit*/}}
        {{if $permFlag}}
                {{ $pbit := .Permissions}}
                {{ $perms := cslice "Create Invite" "Kick Members" "Ban Members" "Administrator" "Manage Channels" "Manage Server" "Add Reactions" "View Audit Log" "Priority Speaker" "Video" "View Channels" "Send Messages" "Send TTS Messages" "Manage Messages" "Embed Links" "Attach Files" "Read Message History" "Mention @everyone" "Use External Emoji" "View Server Insights" "Connect" "Speak" "Mute Members" "Deafen Members" "Move Members" "Use Voice Activity" "Change Nickname" "Manage Nicknames" "Manage Roles" "Manage Webhooks" "Manage Emojis" "Use Slash Commands" "Request to Speak"}}
                {{ $enabled := cslice}}{{ $disabled := cslice}}
                {{ range seq 0 (len $perms) }}
                        {{- if mod $pbit 2 -}}{{- $enabled = $enabled.Append (print "`✅` " (index $perms .)) -}}
                        {{- else -}}{{- $disabled = $disabled.Append (print "`✖️` " (index $perms .)) -}}{{- end -}}
                        {{- $pbit = div $pbit 2 -}}
                {{- end }}
                {{/* Arranging them in columns */}}
                {{ $combined := ($enabled.AppendSlice $disabled).StringSlice }}
                {{ $fields = $fields.AppendSlice (cslice
                (sdict "name" "• Permissions" "value" (joinStr "\n" (slice $combined 0 12)) "inline" true)
                (sdict "name" "​" "value" (joinStr "\n" (slice $combined 12 24)) "inline" true)
                (sdict "name" "​" "value" (joinStr "\n" (slice $combined 24)) "inline" true))}}
                {{$embed.Set "footer" (sdict "text" (print "Triggered By • " $.User.String) "icon_url" ($.User.AvatarURL "256"))}}
        {{end}}
        {{$embed.Set "color" .Color}}
        {{$embed.Set "fields" $fields}}
        {{sendMessage nil (cembed $embed)}}
{{ end }}

{{if $err}}
        {{with $errorMsg}}
                {{sendMessage nil (complexMessage "content" . "embed" $help)}}
        {{ else }}
                {{sendMessage nil $help}}
        {{ end}}
{{ end}}
