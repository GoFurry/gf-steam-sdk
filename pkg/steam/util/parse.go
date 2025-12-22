package util

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

// bbReplacements 定义 Steam 自定义 BBCode 标签到 HTML 标签的映射规则
// 键为正则表达式(忽略大小写), 值为对应的 HTML 替换模板
// bbReplacements defines the mapping rules from Steam custom BBCode tags to HTML tags
// Key is regular expression (case-insensitive), value is corresponding HTML replacement template
var bbReplacements = map[string]string{
	// 文本样式 (Text styles)
	`(?i)\[b\](.*?)\[/b\]`:     `<strong>$1</strong>`, // 粗体 (Bold)
	`(?i)\[i\](.*?)\[/i\]`:     `<em>$1</em>`,         // 斜体 (Italic)
	`(?i)\[u\](.*?)\[/u\]`:     `<u>$1</u>`,           // 下划线 (Underline)
	`(?i)\[s\](.*?)\[/s\]`:     `<s>$1</s>`,           // 删除线 (Strikethrough)
	`(?i)\[h1\](.*?)\[/h1\]`:   `<h1>$1</h1>`,         // 一级标题 (Level 1 heading)
	`(?i)\[h2\](.*?)\[/h2\]`:   `<h2>$1</h2>`,         // 二级标题 (Level 2 heading)
	`(?i)\[h3\](.*?)\[/h3\]`:   `<h3>$1</h3>`,         // 三级标题 (Level 3 heading)
	`(?i)\[p\](.*?)\[/p\]`:     `$1<br>`,              // 段落 (Paragraph, converted to line break)
	`(?i)\[img\](.*?)\[/img\]`: `<img src="$1" />`,    // 图片 (Image)

	// 链接 (Links)
	`(?i)\[url\](.*?)\[/url\]`:       `<a href="$1">$1</a>`, // 普通链接 (Plain link, URL as both href and text)
	`(?i)\[url=(.*?)\](.*?)\[/url\]`: `<a href="$1">$2</a>`, // 带文本链接 (Link with custom text)
}

// ParseBBCode recursively parses Steam custom BBCode into HTML string 将 Steam 自定义 BBCode 递归解析为 HTML 字符串
// 参数说明 (Parameters):
//
//	input - 原始 Steam BBCode 文本 (Original Steam BBCode text)
//	num   - 嵌套解析次数, 用于处理多层嵌套的 BBCode 标签(如 [p][p]内容[/p][/p])
//	        (Number of nested parsing times, used to handle multi-layer nested BBCode tags (e.g. [p][p]content[/p][/p]))
//
// 返回值 (Returns):
//
//	解析后的 HTML 字符串 (Parsed HTML string)
func ParseBBCode(input string, num int) string {
	// 先替换 Steam 图片前缀 (Replace Steam image prefix first)
	input = strings.ReplaceAll(input, "{STEAM_CLAN_IMAGE}", "https://clan.fastly.steamstatic.com/images/")

	// 处理 img/video/youtube 特殊格式标签 (Process special format tags for img/video/youtube)
	input = regexp.MustCompile(`(?i)\[img\s+src="(.*?)"\]\s*\[/img\]`).ReplaceAllString(input, `<img src="$1" />`)
	input = regexp.MustCompile(`(?i)\[video\](.*?)\[/video\]`).ReplaceAllString(input, `<video src="$1" controls></video>`)
	input = regexp.MustCompile(`(?i)\[youtube\](.*?)\[/youtube\]`).ReplaceAllString(input, `<iframe src="https://www.youtube.com/embed/$1" frameborder="0" allowfullscreen></iframe>`)

	// 替换常规 BBCode 标签 最多解析至第n次嵌套 例如:[p][p][p]内容[/p][/p][/p]
	// Replace regular BBCode tags and parse nested tags for n times (e.g. [p][p][p]content[/p][/p][/p])
	for i := 0; i < num; i++ {
		for pattern, replacement := range bbReplacements {
			re := regexp.MustCompile(pattern)
			input = re.ReplaceAllString(input, replacement)
		}
	}

	// 递归解析列表 (Recursively parse lists)
	input = parseLists(input)

	// 将换行转换为 <br> (Convert line breaks to <br>)
	input = strings.ReplaceAll(input, "\n", "<br>")

	return input
}

// parseLists recursively parses [list] (unordered list) and [olist] (ordered list) tags in Steam BBCode
// 递归解析 Steam BBCode 中的 [list](无序列表)和 [olist](有序列表)标签
// 参数说明 (Parameters):
//
//	input - 包含列表标签的 BBCode 文本 (BBCode text containing list tags)
//
// 返回值 (Returns):
//
//	转换后的列表 HTML 字符串 (Converted list HTML string)
func parseLists(input string) string {
	// 匹配无序列表标签(忽略大小写和换行) (Match unordered list tags (case-insensitive and ignore line breaks))
	listPattern := regexp.MustCompile(`(?is)\[list\](.*?)\[/list\]`)
	// 匹配有序列表标签(忽略大小写和换行) (Match ordered list tags (case-insensitive and ignore line breaks))
	olistPattern := regexp.MustCompile(`(?is)\[olist\](.*?)\[/olist\]`)

	// 处理无序列表 (Process unordered lists)
	input = listPattern.ReplaceAllStringFunc(input, func(m string) string {
		content := listPattern.FindStringSubmatch(m)[1]
		items := splitListItems(content)
		var buf strings.Builder
		buf.WriteString("<ul>")
		for _, it := range items {
			// 递归解析列表项中的嵌套列表 (Recursively parse nested lists in list items)
			buf.WriteString("<li>" + parseLists(it) + "</li>")
		}
		buf.WriteString("</ul>")
		return buf.String()
	})

	// 处理有序列表 (Process ordered lists)
	input = olistPattern.ReplaceAllStringFunc(input, func(m string) string {
		content := olistPattern.FindStringSubmatch(m)[1]
		items := splitListItems(content)
		var buf strings.Builder
		buf.WriteString("<ol>")
		for _, it := range items {
			// 递归解析列表项中的嵌套列表 (Recursively parse nested lists in list items)
			buf.WriteString("<li>" + parseLists(it) + "</li>")
		}
		buf.WriteString("</ol>")
		return buf.String()
	})

	return input
}

// splitListItems safely splits BBCode list items and extracts valid list content based on [*] delimiter
// 安全拆分 BBCode 列表项, 基于 [*] 分隔符提取有效列表内容
func splitListItems(content string) []string {
	// 把 [*] 替换成一个特殊分隔符 (Replace [*] with a special delimiter)
	tmp := strings.ReplaceAll(content, "[*]", "\x00")
	parts := strings.Split(tmp, "\x00")
	var items []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			items = append(items, p)
		}
	}
	return items
}

// MarkdownToHTML converts Markdown format text to HTML string 将 Markdown 格式文本转换为 HTML 字符串
// 参数说明 (Parameters):
//
//	markdownContent - 原始 Markdown 文本内容 (Original Markdown text content)
func MarkdownToHTML(markdownContent string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	// 使用 goldmark 库转换 Markdown 为 HTML (Convert Markdown to HTML using goldmark library)
	err := md.Convert([]byte(markdownContent), &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
