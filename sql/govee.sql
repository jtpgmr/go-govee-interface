create schema if not exists govee_smart_tech;

-- Govee devices models registered in the system
create table if not exists govee_smart_tech.device_models(
    serial_id int4 NOT NULL GENERATED ALWAYS AS identity primary key,
    id uuid not null unique,
    model_type int2 not null, -- enums -> 1: LPS (light, plugs, switches), 2: Appliances
    model varchar(10) not null unique,
    created_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NOT NULL
);

create table if not exists govee_smart_tech.device_model_configs(
    serial_id int4 NOT NULL GENERATED ALWAYS AS identity primary key,
    id uuid not null unique,
    device_model_id uuid not null unique references govee_smart_tech.device_models(id),
    min_range uint4 NOT NULL,
    max_range uint4 NOT NULL,
    support_cmds text[] not null,
    created_at timestamptz NOT NULL DEFAULT now(),
    modified_at timestamptz NULL,
    deleted_at timestamptz NULL
);

create table if not exists govee_smart_tech.devices(
    serial_id int4 NOT NULL GENERATED ALWAYS AS identity primary key,
    id uuid not null unique,
    model_id uuid not null references govee_smart_tech.device_models(id),
    name varchar(50) not null,
    mac_address varchar(50) not null unique,
    created_at timestamptz NOT NULL DEFAULT now(),
    modified_at timestamptz NOT NULL,
    deleted_at timestamptz NOT NULL
);

create or replace view 
    govee_smart_tech.devices_summary
as 
    select 
        d.serial_id,
        d.id,
        d.name,
        d.mac_address,
        case d.model_type
            when 1 then 'LPS'
            when 2 then 'Appliance'
            else 'Unknown Model Type'
        end model_type,
        jsonb_build_object(
            'range',
            jsonb_build_object(
                'min', dmc.min_range,
                'max', dmc.max_range
            ),
            'support_cmds', dmc.support_cmds
        ) as model_configs,
        d.created_at,
        d.modified_at,
        d.deleted_at
    from 
        govee_smart_tech.devices d
    left join
        govee_smart_tech.device_models dm
    on
        dm.id = d.model_id
    left join
        govee_smart_tech.device_model_configs dmc
    on
        dmc.device_model_id = d.model_id
	group by 
        d.serial_id,
        d.id,
        d.name,
        d.mac_address,
        d.created_at,
        d.modified_at,
        d.deleted_at
	order by 
        d.serial_id
;


