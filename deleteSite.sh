#!/bin/bash
echo "success: Removing files for  $dir_path"
# Check if the correct number of arguments is provided
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <db_name> <mysql_user> <mysql_password>"
    exit 1
fi

# Assign arguments to variables
db_name="$1"
mysql_user="$2"
mysql_password="$3"

# Define the directory path
dir_path="/var/www/$db_name"

# Check if the directory exists and delete it
if [ -d "$dir_path" ]; then
    rm -rf "$dir_path"
    
else
    echo "failure: Directory $dir_path does not exist."
fi
echo "success: Removing database for  $dir_path"
# Drop the database
mysql -u "$mysql_user" -p"$mysql_password" -e "DROP DATABASE IF EXISTS \`$db_name\`;"
if [ $? -eq 0 ]; then
    echo "success: website: $db_name successfully removed."
else
    echo "failure: to drop database $db_name."
fi

