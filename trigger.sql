CREATE TRIGGER annotation_mapvalue_notify_event
AFTER INSERT OR UPDATE OR DELETE ON annotation_mapvalue
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
