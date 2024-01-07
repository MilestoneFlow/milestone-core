#!/bin/bash
set -e

mongosh <<EOF
use admin
db.createUser({
    user: "flowAdmin",
    pwd: "milestoneFlow123",
    roles: [
        {
            role: "userAdminAnyDatabase",
            db: "admin",
        }, "readWriteAnyDatabase"
    ],
});

use flowDb
db.createCollection("users", { capped: false });
db.createCollection("customers", { capped: false });
db.createCollection("flows", { capped: false });
EOF