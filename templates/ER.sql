{{- range $TableInfo := .tableInfos}}
CREATE TABLE {{$TableInfo.TableName}} {{$TableInfo.TableComment.String}} (
    {{ range $ColumnInfo := $TableInfo.ColumnInfos -}}
    {{- if eq $ColumnInfo.Name "ID" -}}
    {{$ColumnInfo.Name}} {{$ColumnInfo.Kind}} PRIMARY KEY {{$ColumnInfo.Comment.String}},
    {{- else if or (eq $ColumnInfo.Precision.String "") (eq $ColumnInfo.Precision.String "0")}}
    {{$ColumnInfo.Name}} {{$ColumnInfo.Kind}}({{$ColumnInfo.Length.String}}) {{$ColumnInfo.Comment.String}},
    {{- else}}
    {{$ColumnInfo.Name}} {{$ColumnInfo.Kind}}({{$ColumnInfo.Length.String}},{{$ColumnInfo.Precision.String}}) {{$ColumnInfo.Comment.String}},
    {{- end -}}
    {{- end }}
);
{{end}}