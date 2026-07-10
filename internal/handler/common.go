package handler

import "html/template"

// InjectInitialData 供模板使用，生成包含初始 JSON 的 script 标签（以 template.HTML 返回避免转义）。
func InjectInitialData(initialJSON string) template.HTML {
	return template.HTML(`<script id="initial-data" type="application/json">` + initialJSON + `</script>`)
}
