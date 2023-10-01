create table if not exists semaphore.project__host (
    id int(11) NOT NULL AUTO_INCREMENT
        primary key ,
    project_id int(11) NOT NULL,
    name varchar(255) NOT NULL,
    user_name varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    host_ip varchar(255) NOT NULL,
    KEY project__host__index (project_id),
    KEY project__host_ip_index (host_ip) USING BTREE,
    CONSTRAINT project__host_project_id_fk FOREIGN KEY (project_id) REFERENCES semaphore.project (id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


create table if not exists semaphore.project__host__inventory__rel (
    id int(11) unsigned NOT NULL AUTO_INCREMENT
        primary key ,
    project_id int(11) NOT NULL,
    host_id int(11) NOT NULL,
    inventory_id int(11) NOT NULL,
    UNIQUE KEY host_inv_uqi_idx (project_id,host_id,inventory_id) USING BTREE,
    KEY host_index (host_id),
    KEY inventory_index (inventory_id),
    CONSTRAINT host_fk FOREIGN KEY (host_id) REFERENCES semaphore.project__host (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT inventory_fk FOREIGN KEY (inventory_id) REFERENCES semaphore.project__inventory (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;