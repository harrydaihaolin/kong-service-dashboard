CREATE INDEX idx_services_name ON services (name);
CREATE INDEX idx_service_versions_service_id ON service_versions(service_id);
# Foreign keys typically benefit from an index to optimize join queries and ensure fast lookups for related rows in the service_versions table.