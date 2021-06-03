{{/*
        Recommended Trigger: \A(?:\-|<@!?204255221017214977>)\s*(?:r(?:ole)?\s*i(?:nfo)?)(?: +|\z)
        Trigger Type: Regex
        Usage: -roleinfo [ID/Mention/Name/Position]

        Copyright (c): zen | ゼン#0008; 2021
        License: MIT
        Repository: https://github.com/z3nn13/custom-commands-yagpdb
        
*/}}

{{/* Global Variables */}}
{{$help := cembed "title" "RoleInfo/RInfo/Ri" "description" (print "\x60\x60\x60Roleinfo [Role:Name/Mention/ID/Position]\x60\x60\x60Shows information about a role")}}
{{$role := ""}}{{$errorMsg := ""}}{{$err := false}}{{$Args := .StrippedMsg}}
{{$guildRoles := cslice.AppendSlice .Guild.Roles}} {{$listroles := split (exec "listroles" ) "\n"}}

{{/* Checking Input */}}
{{if .CmdArgs}}
        {{ if $roleID := reFind `\d{17,19}` $Args|toInt64}}{{/* Mention or ID */}}
                {{ with .Guild.Role $roleID}}
                        {{ $role = .}}
                {{ else }}
                        {{ $errorMsg = "Invalid ID/Role Does Not Exist"}}{{$err = true}}
                {{ end }}
        {{else}}
                {{$found := false}}
                {{/* Name Input */}}
                {{range $guildRoles}}
                        {{- if or (eq (.Name|lower) ($Args|lower)) (eq (.Name|lower) (index $.CmdArgs 0|lower)) -}}
                                {{- $role = .}}{{$found = true -}}
                        {{- end -}}
                {{- end}}
                {{/* Position Input */}}
                {{ if $index := toInt $Args}}
                        {{if le $index (len $guildRoles)}}
                                {{ $role = add $index 1 | index $listroles | reFind `\d{17,19}` | toInt64 | .Guild.Role }}{{$found = true}}
                        {{end}}
                {{end}}
                {{if not $found}}
                        {{$errorMsg =printf "%q not recognized: Invalid Name/Position" $Args}}{{$err = true}}
                {{end}}
        {{ end }}
{{else}}
        {{$err = true}}
{{end}}

{{/* Preparing Embed */}}
{{ $ex := or (and (reFind "a_" $.Guild.Icon) "gif" ) "png" }}
{{ $icon := print "https://cdn.discordapp.com/icons/" $.Guild.ID "/" $.Guild.Icon "." $ex "?size=1024" }}
{{ $embed := sdict "author" 
( sdict  "name" "Role Info" "icon_url" "https://images-ext-2.discordapp.net/external/G67VOLJZEh_p_JpowOPIRo4LimCe4KNMj7X5Azffufc/https/cdn.discordapp.com/emojis/764251327814696970.png")
"thumbnail" (sdict "url" $icon)}}

{{with $role}}
        {{ $createdAt := div .ID 4194304 | add 1420070400000 | mult 1000000 | toDuration | (newDate 1970 1 1 0 0 0).Add }}
        {{ $mentionable := or (and $role.Mentionable "`Yes`") "`No`" }}
        
        {{/* Fetching Position */}}
        {{ $pos := 0 }}{{ $up := "" }}{{ $down := "" }}
        {{ range $i,$v := $listroles }}{{ if reFind (str $role.ID) $v }}{{ $pos = sub $i 2 }}{{ end }}{{ end }}
        
        {{/* Note: .Role.Position had a wonky order, so this was a workaround */}}
        {{ if eq $pos 0 }}
                {{ $up = "-----------\n" }}
        {{ else }}
                {{ $up = printf "> #%d • %s\n" ($pos) ((sub $pos 1 |index $guildRoles).Mention) }}                
        {{ end }}

        {{ if eq $pos (len $guildRoles|add -1) }}
                {{ $down = "-----------\n"}}
        {{ else }}
                {{ $down = printf "> #%d • %s\n" (add $pos 2) (add $pos 1|index $guildRoles).Mention }}
        {{ end }}
        {{ $final_pos := printf "%s> **#%d • %s**\n%s\n> `.Position` = %d\n> (Total Roles: **%d**)" $up (add $pos 1) .Mention $down .Position (len $guildRoles)}}


        {{$fields := cslice
        ( sdict "name" "• Name" "value" (print .Mention) "inline" true) (sdict "name" "• ID" "value" (print .ID) "inline" true)
        ( sdict "name" "• Created At" "value" (printf "%s\n%s ago" ($createdAt.Format "Monday, January 2, 2006 at 3:04 PM MST") (humanizeTimeSinceDays $createdAt)) "inline" true)
        ( sdict "name" "• Position ↓" "value" $final_pos "inline" true)
        ( sdict "name" "• Color" "value" (printf "#%x" .Color ) "inline" true)
        ( sdict "name" "• Mentionable" "value" $mentionable "inline" true)}}
        
        {{/* Credits To Satty#9361 */}}
        {{ $pbit := .Permissions}}
        {{ $perms := cslice "Create Invite" "Kick Members" "Ban Members" "Administrator" "Manage Channels" "Manage Server" "Add Reactions" "View Audit Log" "Priority Speaker" "Video" "View Channels" "Send Messages" "Send TTS Messages" "Manage Messages" "Embed Links" "Attach Files" "Read Message History" "Mention @everyone" "Use External Emoji" "View Server Insights" "Connect" "Speak" "Mute Members" "Deafen Members" "Move Members" "Use Voice Activity" "Change Nickname" "Manage Nicknames" "Manage Roles" "Manage Webhooks" "Use Slash Commands" "Request to Speak"}}
        {{ $enabled := cslice}}{{ $disabled := cslice}}
        {{ range seq 0 (len $perms) }}
                {{- if mod $pbit 2 -}}{{- $enabled = $enabled.Append (print "<:c_:838811489581137961> " (index $perms .)) -}}
                {{- else -}}{{- $disabled = $disabled.Append (print "<:x_:832257188168859649> " (index $perms .)) -}}{{- end -}}
                {{- $pbit = div $pbit 2 -}}
        {{- end }}
        {{/* Arranging them in columns */}}
        {{ $split := split (print (joinStr "\n" $enabled.StringSlice) "\n" (joinStr "\n" $disabled.StringSlice)) "\n"}}
        {{ $fields = $fields.AppendSlice (cslice
        (sdict "name" "• Permissions" "value" (joinStr "\n" (slice $split 0 11)) "inline" true)
        (sdict "name" "​" "value" (joinStr "\n" (slice $split 11 22)) "inline" true)
        (sdict "name" "​" "value" (joinStr "\n" (slice $split 22)) "inline" true))}}


        {{$embed.Set "color" .Color}}
        {{$embed.Set "footer" (sdict "text" (print "Triggered By • " $.User.String) "icon_url" ($.User.AvatarURL "256"))}}
        {{$embed.Set "fields" $fields}}
        {{ sendMessage nil (cembed $embed)}}
{{ end }}

{{ if $err}}
        {{ with $errorMsg}}
                {{sendMessage nil (complexMessage "content" . "embed" $help)}}
        {{ else }}
                {{sendMessage nil $help}}
        {{ end}}
{{ end}}
