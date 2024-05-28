#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <db_name> <mysql_user> <mysql_password>"
    exit 1
fi



# Assign arguments to variables
db_name="$1"
mysql_user="$2"
mysql_password="$3"
echo "<br>删除 $db_name 的文件。"
dir_path="/var/www/$db_name"

if [ -d "$dir_path" ]; then
    rm -rf "$dir_path"
    
else
echo "<br><span style=\"color:red\">目录 "$dir_path" 不存在。</span>"
exit 1
fi
echo "<br>删除数据库：$db_name。"

if ! error_message=$(mysql -u "$mysql_user" -p"$mysql_password" -e "DROP DATABASE IF EXISTS \`$db_name\`;" 2>&1); then
    echo "<br><span style=\"color:red\">无法删除 $db_name 的数据库。<br> . 错误: $error_message </span>"
    rm -rf "/var/www/$db_name"
    exit 1
else
    echo "<br><br> <span style=\"color:chartreuse\">网站 $db_name 已成功删除。</span>"
fi


