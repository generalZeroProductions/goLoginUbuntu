#!/bin/bash

drop_table() {
    local db_name="$1"
    local mysql_user="$3"
    local mysql_password="$4"

    mysql -u "$mysql_user" -p"$mysql_password" -e "DROP DATABASE IF EXISTS $db_name;" 2>/dev/null

    if ! error_message=$(mysql -u "$mysql_user" -p"$mysql_password" -e "USE ${db_name}; $drop_table_sql" 2>&1); then
        echo "<span style=\"color:red\"> 删除数据库失败</span> . Error: $error_message"

    fi
}

if [ -z "$1" ]; then
    echo "<span style=\"color:red\"> 未提供网站名称</span>"
    exit 1
fi

if [ -z "$2" ]; then
    echo "<span style=\"color:red\"> 未提供用户名</span>"
    exit 1
fi

if [ -z "$3" ]; then
    echo "<span style=\"color:red\"> 没有提供密码</span>"
    exit 1
fi

# Assign input arguments to variables
dir_name="$1"
mysql_user="$2"
mysql_password="$3"
ip_address="4"
server_ip=$(hostname -I | awk '{print $1}')

if [ -z "$server_ip" ]; then
    echo "<span style=\"color:red\"> 无法检索服务器 IP 地址</span>"
    exit 1
fi
# Check if the directory already exists
if [ -d "/var/www/$dir_name" ]; then
    echo "<span style=\"color:red\"> 目录已存在 IP 地址</span>"
    exit 1
fi

# Alternative check for directory existence (redundant, can be removed if the first check is sufficient)
if stat -t "/var/www/$dir_name" >/dev/null 2>&1; then
    echo "<span style=\"color:red\"> 目录已存在 IP 地址</span>"
    exit 1
fi
echo "开始网站创建<br>"
# Create the directory
mkdir "/var/www/$dir_name"
if [ $? -ne 0 ]; then
    echo "<span style=\"color:red\"> 创建目录失败</span>"
    exit 1
fi

source_dir="/var/www/siteBase/siteBase"

# Define the destination directory
destination_dir="/var/www/$dir_name"

cp -r "$source_dir"/.* "$destination_dir"
# Change into the directory

cd "/var/www/$dir_name"
if [ $? -ne 0 ]; then
    echo "<span style=\"color:red\"> 无法进入目录</span>"
    rm -rf "/var/www/$dir_name"
    exit 1
fi

echo "为新网站建立数据库<br>"
# Create a new MySQL database
if ! error_message=$(mysql -u "$mysql_user" -p"$mysql_password" -e "CREATE DATABASE ${dir_name};" 2>&1); then
    echo "<span style=\"color:red\"> 创建数据库失败</span> . 错误: $error_message"
    rm -rf "/var/www/$dir_name"
    exit 1
fi

echo "配置新网站<br>"
# Copy .env.example to .env
env_file=".env"
cp .env.example "$env_file"
if [ $? -ne 0 ]; then
    echo "<span style=\"color:red\"> 无法将 .env.example 复制到 .env</span>"
    drop_database "$dir_name" "$mysql_user" "$mysql_password"
fi

# # Update .env file with the required values


sed -i "s/^APP_NAME=.*/APP_NAME=${dir_name}/" "$env_file"
sed -i "s/^APP_ENV=.*/APP_ENV=local/" "$env_file"
sed -i "s/^APP_KEY=.*/APP_KEY=base64:fmLIgNQQQa7xbLP8ZmfRAxPDWve3B8FeKfRN+YKlN9M=/" "$env_file"
sed -i "s/^APP_DEBUG=.*/APP_DEBUG=true/" "$env_file"
sed -i "s|^APP_URL=.*|APP_URL=http://$server_ip|" "$env_file"
sed -i "s/^DB_CONNECTION=.*/DB_CONNECTION=mysql/" "$env_file"
sed -i "s/^DB_HOST=.*/DB_HOST=127.0.0.1/" "$env_file"
sed -i "s/^DB_PORT=.*/DB_PORT=3306/" "$env_file"
sed -i "s/^DB_DATABASE=.*/DB_DATABASE=${dir_name}/" "$env_file"
sed -i "s/^DB_USERNAME=.*/DB_USERNAME=${mysql_user}/" "$env_file"
sed -i "s/^DB_PASSWORD=.*/DB_PASSWORD=${mysql_password}/" "$env_file"
if [ $? -ne 0 ]; then
    echo "<span style=\"color:red\"> 更新 .env 文件失败</span>"
    drop_database "$dir_name" "$mysql_user" "$mysql_password"
    exit 1

fi
echo "正在为新网站填充数据库<br>"



# Run Laravel migrations
echo "Attempting to change ownership of /var/www/$dir_name/storage to www-data:www-data"

# Check if the directory exists
if [ -d "/var/www/$dir_name/storage" ]; then
    echo "Directory exists: /var/www/$dir_name/storage"
else
    echo "Directory does not exist: /var/www/$dir_name/storage"
    exit 1
fi

php artisan migrate
if [ $? -ne 0 ]; then
    echo "<span style=\"color:red\"> 无法运行迁移</span>"
    drop_database "$dir_name" "$mysql_user" "$mysql_password"
    exit 1
fi

# # Seed the database
php artisan db:seed --class=InitDbSeeder
if [ $? -ne 0 ]; then
    echo "<span style=\"color:red\">未能为新数据库设定种子</span>"
    drop_database "$dir_name" "$mysql_user" "$mysql_password"
    exit 1
fi
sudo chown -R www-data:www-data "/var/www/$dir_name/storage" 2> chown_error.log

# Debug: Verify the command worked
if [ $? -eq 0 ]; then
    echo "Ownership changed successfully."
else
    echo "Failed to change ownership. Check chown_error.log for details."
    cat chown_error.log
fi

systemctl restart nginx
echo "<br>
<br><span style=\"color:chartreuse\">完成创建网站$dir_name。如果您的 DNS 已正确设置为在以下位置查看，<br>则 www.$dir_name.com 将可见并可供编辑</span>"
