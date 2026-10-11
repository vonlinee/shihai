import os
import subprocess
import sys
import unittest
from pathlib import Path


TEST_DIR = Path(__file__).resolve().parent
SCRIPT = TEST_DIR.parent / "shihai.py"


def run_script(*args: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        capture_output=True,
        text=True,
        check=False,
    )


class ShihaiScriptTests(unittest.TestCase):
    def test_help(self) -> None:
        result = run_script("--help")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("python scripts/shihai.py <command>", result.stdout)
        self.assertIn("start-development-frontend", result.stdout)

    def test_development_dry_run(self) -> None:
        result = run_script("start", "--mode", "development", "--backend-port", "9090", "--dry-run")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("Mode: development", result.stdout)
        self.assertIn("Backend port: 9090", result.stdout)
        self.assertIn("Frontend URL: http://localhost:3000", result.stdout)
        self.assertIn("Backend URL: http://localhost:9090", result.stdout)
        self.assertIn("Process count: 2", result.stdout)
        self.assertIn("go run ./cmd/server -port 9090", result.stdout)

    def test_standalone_dry_run(self) -> None:
        deploy_path = r"C:\tmp\shihai-python-test" if os.name == "nt" else "/tmp/shihai-python-test"
        result = run_script(
            "start", "--frontend-mode", "standalone", "--deploy-path", deploy_path,
            "--backend-port", "9090", "--frontend-port", "4173", "--no-redeploy", "--dry-run",
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("Frontend mode: standalone", result.stdout)
        self.assertIn("Process count: 2", result.stdout)
        self.assertIn("--port 4173", result.stdout)
        self.assertIn("Frontend URL: http://localhost:4173", result.stdout)
        self.assertIn("Backend URL: http://localhost:9090", result.stdout)

    def test_invalid_options(self) -> None:
        invalid_port = run_script("start", "--backend-port", "0", "--dry-run")
        self.assertNotEqual(invalid_port.returncode, 0)
        self.assertIn("Backend port must be an integer", invalid_port.stderr)

        incompatible = run_script("deploy", "--frontend-mode", "embedded", "--skip-frontend")
        self.assertNotEqual(incompatible.returncode, 0)
        self.assertIn("Embedded mode requires both frontend and backend", incompatible.stderr)


if __name__ == "__main__":
    unittest.main()
