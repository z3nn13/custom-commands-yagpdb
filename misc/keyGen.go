{{/* 
        Trigger: keygen
        Trigger Type: Command Type
        Usage: -keygen <length>
        */}}

{{$help := cembed "title" "Keygen" "description" "```Keygen <Length:Whole number>```Generates a key base on your length"}}

{{$err := false}}{{$errMsg := ""}}
{{with .CmdArgs}}
        {{with $length := index . 0|toInt}}
                {{ if gt . 10000}}
                        {{$errMsg = "Length must be under 10k Limit"}}
                {{ else }}
                        {{ $rLetters := split "abcdefghijklmnopqrstuvwxyz" "" }}
                        {{ $code := "" }}
                        {{ range seq 0 .}}
                                {{- if randInt 2}}
                                        {{- $code = randInt 10 | print $code -}}
                                {{- else -}}
                                        {{- if randInt 2 -}}
                                                {{- $code = index ($rLetters|shuffle) 0 | upper | print $code -}}
                                        {{- else -}}
                                                {{- $code = index ($rLetters|shuffle) 0 | lower | print $code -}}
                                        {{- end -}}
                        {{- end}}{{end}}
                        {{if ge (len $code) 1975}}
                                {{ sendMessage nil (complexMessage "file" $code)}}
                        {{else}}
                                {{ sendMessage nil (print "🔑 **Key Generated :**\n> `" $code "`")}}
                        {{end}}
                {{ end}}
        {{ else }}
                {{ $errMsg = printf "Unknown Length %q : Length must be a whole number" (index . 0)}}
        {{ end}}
{{else}}
        {{sendMessage nil $help}}
{{end}}

{{with $errMsg}}
        {{sendMessage nil (complexMessage "content" . "embed" $help)}}
{{end}}
