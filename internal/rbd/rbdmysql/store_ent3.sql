SELECT
    console.tenant_info.tenant_alias AS 团队名,
    console.tenant_info.tenant_name AS 团队英文名,
    console.service_group.group_name AS 应用名,
    console.service_group.k8s_app AS 应用英文名,
    console.tenant_service.service_cname AS 组件名,
    console.tenant_service.k8s_component_name AS 组件英文名,
    console.tenant_service.image AS 应用镜像,
    CONCAT(console.service_group.k8s_app,"-",console.tenant_service.k8s_component_name) AS deployment名称,
    region.tenant_services.replicas as 副本数,
    console.tenant_service.min_cpu AS 组件CPU,
    console.tenant_service.min_memory AS 组件内存,
    console.tenant_service.cmd AS 启动命令,
    console.component_k8s_attributes.attribute_value AS 高级资源配置
FROM console.tenant_info
         LEFT JOIN console.service_group
                   ON console.tenant_info.tenant_id = console.service_group.tenant_id COLLATE utf8mb4_unicode_ci
         LEFT JOIN console.service_group_relation
                   ON console.service_group.ID = console.service_group_relation.group_id COLLATE utf8mb4_unicode_ci
         LEFT JOIN console.tenant_service
                   ON console.service_group_relation.service_id = console.tenant_service.service_id COLLATE utf8mb4_unicode_ci
         LEFT JOIN region.tenant_services on region.tenant_services.service_id = console.tenant_service.service_id COLLATE utf8mb4_unicode_ci
         LEFT JOIN console.component_k8s_attributes
                   ON console.tenant_service.service_id = console.component_k8s_attributes.component_id AND console.component_k8s_attributes.name="resources" COLLATE utf8mb4_unicode_ci;