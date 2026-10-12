"""Regression checks for generated files that only appear after generation."""
import pathlib
import tempfile
import unittest

from run import generated_changes, sha


class GeneratedFilesTest(unittest.TestCase):
    def test_complete_generated_inventory(self):
        with tempfile.TemporaryDirectory() as directory:
            snapshot = pathlib.Path(directory)
            existing = snapshot / "existing_templ.go"
            existing.write_bytes(b"current generated output")
            hashes = {existing.name: sha(existing.read_bytes())}
            self.assertEqual(generated_changes(snapshot, hashes), [])

            # A deleted companion or a newly added template produces a path
            # absent from the pre-generation checkout inventory.
            added = snapshot / "new_templ.go"
            added.write_bytes(b"new generated output")
            self.assertEqual(generated_changes(snapshot, hashes), [added.name])

            added.unlink()
            existing.write_bytes(b"changed generated output")
            self.assertEqual(generated_changes(snapshot, hashes), [existing.name])

            existing.unlink()
            self.assertEqual(generated_changes(snapshot, hashes), [existing.name])


if __name__ == "__main__":
    unittest.main()
