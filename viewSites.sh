#!/bin/bash

echo "<span style=\"color:deepskyblue\"> 该服务器上的站点列表：</span><br>"
echo "<span style=\"color:deepskyblue\"> _______________________</span><br>"
# Check if the /var/www directory exists
if [ ! -d "/var/www" ]; then
    echo "<span style=\"color:red\"> 无法访问主 Web 目录。</span><br>"
    exit 1
fi

# Loop through each directory in /var/www and print its name and disk space used
for dir in /var/www/*; do
    if [ -d "$dir" ]; then
        dir_name=$(basename "$dir")
        # Skip directories with "HTML" in their name
        if [[ "$dir_name" == *"html"* ]]; then
            continue
        fi
          if [[ "$dir_name" == *"siteBas"* ]]; then
            continue
        fi
        dir_size=$(du -sh "$dir" | awk '{print $1}')
        echo " $dir_name - $dir_size <br>"
    fi
done

# Print remaining disk space on the server
remaining_space=$(df -h / | awk 'NR==2 {print $4}')

echo "<br>"
echo "<span style=\"color:deepskyblue\"> _______________________</span><br>"
echo "<span style=\"color:deepskyblue\">  服务器剩余磁盘空间： $remaining_space"
