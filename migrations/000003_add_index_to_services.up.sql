CREATE INDEX idx_services_name ON services(service_name);
CREATE INDEX idx_services_versions_service_id ON service_versions(service_id);