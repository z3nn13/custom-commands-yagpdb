{{/* 
        Trigger: keygen
        Trigger Type: Command Type
        Usage: -keygen <length>

        Copyright (c): zen | ゼン#0008; 2021
        License: MIT
        Repository: https://github.com/z3nn13/custom-commands-yagpdb
        */}}

{{$help := cembed "title" "Keygen" "description" "\x60\x60\x60Keygen <Length:Whole number>\x60\x60\x60Generates a key base on your length"}}

{{$err := false}}{{$errMsg := ""}}
{{with .CmdArgs}}
        {{$length := index . 0|toInt}}
        {{ if gt $length 10000}}
                {{$errMsg = "Length must be under 10k Limit"}}
                {{$err = true}}
        {{ end }}
        {{if not $err}}
                {{ if $length }}
                        {{ $rLetters := split "abcdefghijklmnopqrstuvwxyz" "" }}
                        {{ $rNumbers := split "1234567890" "" }}
                        {{ $code := "" }}
                        {{ range seq 0 $length}}
                                {{- $x := randInt 2 -}}
                                {{- if eq $x 0 -}}
                                        {{- $code = print $code (index ($rNumbers|shuffle) 0) -}}
                                {{- else -}}
                                        {{- $capital := randInt 2 -}}
                                        {{- if eq $capital 0 -}}
                                                {{- $code = print $code ((index ($rLetters|shuffle) 0)|upper) -}}
                                        {{- else -}}
                                                {{- $code = print $code ((index ($rLetters|shuffle) 0)|lower) -}}
                                        {{- end -}}
                        {{- end}}{{end}}
                        {{if ge (len $code) 1973}}
                                {{ sendMessage nil (complexMessage "file" (print $code ))}}
                        {{else}}
                                {{ sendMessage nil (print "Code : " $code)}}
                        {{end}}
                {{ else }}        
                        {{ $errMsg = printf "Unknown Length %q : Length must be a whole number" (index . 0)}}
                        {{ $err = true}}
                {{ end}}
        {{ end }}
{{else}}
        {{$err = true}}
{{end}}

{{if $err}}
        {{with $errMsg}}
                {{sendMessage nil (complexMessage "content" . "embed" $help)}}
        {{else}}
                {{sendMessage nil $help}}
        {{end}}
{{end}}
