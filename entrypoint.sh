#!/bin/bash
PGPASSWORD=$CONFIG_omero_db_pass psql -h db -U $CONFIG_omero_db_user -d $CONFIG_omero_db_name -f /root/notify.sql
PGPASSWORD=$CONFIG_omero_db_pass psql -h db -U $CONFIG_omero_db_user -d $CONFIG_omero_db_name -f /root/trigger.sql
# /usr/local/bin/pgweb --bind 0.0.0.0 --ssl disable --host db --user $CONFIG_omero_db_user --db $CONFIG_omero_db_name --pass $CONFIG_omero_db_pass
/usr/local/bin/omero2cnb db $CONFIG_omero_db_name $CONFIG_omero_db_user $CONFIG_omero_db_pass $CONFIG_omero_web_server
