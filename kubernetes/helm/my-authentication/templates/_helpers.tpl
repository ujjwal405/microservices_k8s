{{- define "my-authentication-name" }}
{{- printf "%s-%s" .Release.Name .Values.deploymentName -}}
{{- end }}

{{- define "my-label" -}}
app: {{ .Values.label }}
{{- end -}}


{{- define "my-authentication-service" }}
{{- printf "%s-%s" .Release.Name .Values.serviceName -}}
{{- end }}


