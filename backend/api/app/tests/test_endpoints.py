from fastapi.testclient import TestClient

from app.main import app


def test_healthcheck():
    with TestClient(app) as client:
        res = client.get("/healthcheck")
        assert res.status_code == 200
