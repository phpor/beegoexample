<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title></title>
</head>
<body>
已上传文件（点击下载）<br />
<ul>
    {{range .paths}}
    <li><a href='/api/file/download?name={{.}}'>{{.}}</a></li>
    {{end}}
</ul>
<form method="post" action="/api/file/upload" enctype="multipart/form-data">
    文件名：<input type="text" name="filename" value=""/><br/>
    <input type="file" name="file"/> <br/>
    <input type="submit" value="提交"/>
</form>
</body>
</html>