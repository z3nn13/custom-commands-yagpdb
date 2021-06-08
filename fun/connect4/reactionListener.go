{{/*
    Trigger Type: Reaction ; Added Reactions Only
    
    Copyright (c): zen | ゼン#0008; 2021
    License: MIT
    Repository: https://github.com/z3nn13/custom-commands-yagpdb
*/}}

{{- define "board_maker" -}}
    {{- $emojis := sdict "0" "⚫" "1" "🔴" "2" "🟡" "10" "🔵" -}}
    {{- $board := cslice.AppendSlice .board -}}
    {{- $visual := "" -}}
    {{- range $board -}}
        {{- range . -}}
                {{- $visual = printf "%s%s " $visual (str . | $emojis.Get) -}}
        {{- end -}}
        {{- $visual = print $visual "\n" -}}
    {{- end -}}
    {{- $visual = print $visual "1️⃣ 2️⃣ 3️⃣ 4️⃣ 5️⃣ 6️⃣ 7️⃣" -}}
    {{- $embed := sdict "author" (sdict "name" (print .cPlayer.Username "'s turn") "icon_url" (.cPlayer.AvatarURL "256")) 
        "title" "Connect4" 
        "description" $visual
        "fields" (cslice 1 2)
        "color" 0x1c25a3 "footer" (sdict "text" "Powered by • Yagpdb.xyz" ) -}}
    {{- $embed.fields.Set (sub .turn 1) (sdict "name" (print "Player " .turn) "value" (print "> " .cPlayer.Mention "\nToken: " (str .turn|$emojis.Get)) "inline" true) -}}
    {{- $embed.fields.Set (sub .nextTurn 1 ) (sdict "name" (print "Player " .nextTurn) "value" .nPlayer.Mention "inline" true ) -}}    
    {{- .Set "embed" $embed -}}
{{- end}}

{{ $data := sdict}}{{$players := cslice}}{{$cPlayer := ""}}{{$turn := ""}}{{$nextTurn := ""}}{{$nPlayer := ""}}{{$input := ""}}{{$flag := false}}
{{ $store := sdict  "1" (sdict "emoji" "🔴" "color" 0xff4d12) "2" (sdict "emoji" "🟡" "color" 0xfff457)}}

{{if $db := dbGet 2021 "connect4"}}
    
    {{/* Global variables */}}
    {{ $data = sdict $db.Value}}
    {{ $players = $players.AppendSlice $data.players}}
    {{ $turn = $data.turn}}
    {{ $cPlayer = (index $players $turn|toInt64|getMember).User}}
    {{ $nextTurn = sub 3 $turn}}
    {{ $nPlayer = (index $players $nextTurn|toInt64|getMember).User}}

    {{ if eq ($data.msgID|toInt64) .Message.ID}}
            {{ deleteMessageReaction nil .Message.ID .User.ID (print (or (and ($k:=.Reaction.Emoji).Animated "a:") "") $k.Name (or (and $k.ID (print ":" $k.ID)) ""))}}
            {{/* quit reaction */}}
            {{ if and (eq .Reaction.Emoji.ID 844556617085485058) (in $players (str .User.ID)) }}
                    {{ dbDel 2021 "connect4"}}
                    {{ deleteAllMessageReactions nil .Message.ID }}
                    {{ $tempData := (dict "board" $data.board "turn" $turn "nextTurn" $nextTurn "cPlayer" $cPlayer "nPlayer" $nPlayer)}}
                    {{ template "board_maker" $tempData}}
                    {{ $otherPlayer := or (and (eq $cPlayer.ID .User.ID) $nPlayer) $cPlayer}}
                    {{ $tempData.embed.set "author" (sdict "name" "Game Over" "icon_url" ($otherPlayer.AvatarURL "256"))}}
                    {{ $tempData.embed.set "color" 0xffdc42}}
                    {{ editMessage nil .Message.ID (complexMessageEdit "content" (print "> You have left the game\n**Winner: **" $otherPlayer.Mention " !") "embed" (cembed $tempData.embed))}}
            {{ else }}
                    {{ if eq $cPlayer.ID $.User.ID}}
                            {{ if and ($temp := (index (toRune .Reaction.Emoji.Name) 0|printf "%c"|toInt)) (le $temp 7)}}
                                    {{$input = sub $temp 1}}
                                    {{$flag = true}}
                            {{end}}
                    {{ end }}
            {{ end }}
    {{ end }}
{{ end }}

{{if $flag}}
        {{ $tempData := (dict "board" $data.board "input" $input "turn" $turn "nextTurn" $nextTurn "cPlayer" $cPlayer "nPlayer" $nPlayer "full" false)}}
        {{ template "slot_checker" $tempData}}
        {{ template "board_maker" $tempData}}
        
        {{ if $tempData.full}}
                {{ editMessage nil .Message.ID (complexMessageEdit "content" (printf "> %s, %q is already full.\nPlease react another slot (1-7)" .User.Mention ($input|add 1))
                "embed" (cembed $tempData.embed))}}
        {{else}}
            {{ template "win_checker" $tempData}}
            {{ $msg := ""}}
            {{ if $tempData.gameWon }}
                    {{ $tempData.embed.Set "color" (str $turn|$store.Get).color}}
                    {{ $tempData.embed.author.Set "name" (printf "Game Over • %s Wins !" (str $turn|$store.Get).emoji)}}
                    {{ $msg = printf "**%s** vs **%s**\n **Winner:** %s" $nPlayer.Mention .User.Mention .User.Mention}}
                    {{ deleteAllMessageReactions nil .Message.ID}}
                    {{ dbDel 2021 "connect4" }}
            {{ else if $tempData.gameTie }}
                    {{ $tempData.embed.Set  "color" 0xba19ff}}
                    {{ $msg = print "Owo what's this, the match is a draw"}}
                    {{ deleteAllMessageReactions nil .Message.ID}}
                    {{ dbDel 2021 "connect4"}}
            {{ else}}
                    {{ $tempData.embed.Set "color" (str $nextTurn|$store.Get).color}}
                    {{ $tempData.embed.Set "author" (sdict "name" (print $nPlayer.Username "'s turn") "icon_url" ($nPlayer.AvatarURL "256"))}}
                    {{ $tempData.embed.fields.Set (sub $nextTurn 1) (sdict "name" (print "Player " $nextTurn) "value" (print "> " $nPlayer.Mention "\nToken: " (str $nextTurn|$store.Get).emoji) "inline" true)}}
                    {{ $tempData.embed.fields.Set (sub $turn 1 ) (sdict "name" (print "Player " $turn) "value" $cPlayer.Mention "inline" true)}}
                    {{ $msg = printf "> **%s** dropped token in slot %d⃣\n%s, Please pick a slot" .User.Username ($input|add 1) $nPlayer.Mention}}
                    {{ $data.Set "turn" ($nextTurn|toInt) }}
                    {{ $data.Set "board" $tempData.board }}
                    {{ dbSet 2021 "connect4" $data }}
            {{ end }}
            {{ editMessage nil .Message.ID (complexMessageEdit "content" $msg "embed" (cembed $tempData.embed))}}
        {{end}}
{{end}}


{{- define "slot_checker" -}}
    {{- $board := cslice.AppendSlice .board }}{{ $turn := .turn }}{{ $input := .input }}{{ $full := false }}{{$position := 0 -}}{{ $found := false}}
    {{- $verti := cslice -}}
    {{- range $board -}}
        {{- $verti = $verti.Append (index . $input) -}}
    {{- end -}}
    {{- range $i,$v := $verti -}}
        {{- if not $found -}}
            {{- if and (eq $i 0) (not (eq $v 10)) $v -}}
                {{- $full = true -}}
            {{- else -}}
                {{- if and (not $v) (eq $i (len $verti|add -1)) -}}
                    {{- $position = $i -}}
                    {{- $found = true -}}
                {{- else if and (eq $v 10|not) $v -}}
                    {{- $position = sub $i 1 -}}
                    {{- $found = true -}}
                {{- else if not $v -}}
                    {{- $convert := cslice.AppendSlice (index $board $i)}}
                    {{- $convert.Set $input 10 -}}
                    {{- $board.Set $i $convert}}
                {{- end -}}
            {{- end -}}
        {{- end -}}
    {{- end}}
    {{- range $rowIndex,$row := $board -}}
        {{- range $col,$v := $row -}}
            {{- if and (eq $col $input|not) (eq $v 10) -}}
                {{- $convert := cslice.AppendSlice (index $board $rowIndex) -}}
                {{- $convert.Set $col 0 -}}
                {{- $board.Set $rowIndex $convert -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
    {{- if and $found (not $full) -}}
        {{- $convert := cslice.AppendSlice (index $board $position) -}}
        {{- $convert.Set $input $turn -}}
        {{- $board.Set $position $convert -}}
    {{- end -}}
    {{- .Set "board" $board -}}
    {{- .Set "full" $full -}}
    {{- .Set "position" $position -}}
{{- end -}}

{{- define "win_checker" -}}
    {{- $gameWon := false}}{{$board := cslice.AppendSlice .board -}}{{$turn := .turn}}{{$input := .input}}{{$position := .position -}}
    {{/* horizontal checking */}}
    {{- $check  := index $board $position -}}
    {{- range $col,$v := $check -}}
        {{- if and $v (not (eq $v 10)) (lt $col (len $check|add -3)) -}}
            {{- if and (eq $v (add $col 1|index $check)) (eq $v (add $col 2|index $check)) (eq $v (add $col 3|index $check)) -}}
                {{- $gameWon = true -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
    {{/* vertical checking */}}
    {{- $verti := cslice -}}
    {{- range $board -}}
        {{- $verti = $verti.Append (index . $input) -}}
    {{- end -}}
    {{- range $i,$v := $verti -}}
        {{- if and $v (not (eq $v 10)) (lt $i (len $verti|add -3)) -}}
            {{- if and (eq $v (add $i 1|index $verti)) (eq $v (add $i 2|index $verti)) (eq $v (add $i 3|index $verti)) -}}
                {{- $gameWon = true -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}

    {{/* diagonal checking */}}
    {{- $marker := cslice -}}
    {{- $total := cslice -}}
    {{- range $rowIndex, $row := $board -}}
        {{- range $col,$v := $row -}}
            {{- if and $v (eq $v 10|not) -}}
                {{$total = $total.Append $v -}}
            {{- end -}}
            {{- if eq $v $turn -}}
                {{- $marker = $marker.Append (printf "%s%s" (str $col) (str $rowIndex)) -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
    {{- range $marker -}}
        {{- if and (in $marker (.|toInt|add 11|str)) (in $marker (.|toInt|add 22|str)) (in $marker (.|toInt|add 33|str)) -}}
            {{- $gameWon = true -}}
        {{- else if and (in $marker (.|toInt|add -9|str)) (in $marker (.|toInt|add -18|str)) (in $marker (.|toInt|add -27|str)) -}}
            {{- $gameWon = true -}}
        {{- end -}}
    {{- end -}}
    {{- if and (not $gameWon) (eq (len $total) 42) -}}
    {{- .Set "gameTie" true -}}
    {{- end -}}
    {{- .Set "gameWon" $gameWon -}}
{{- end }}
