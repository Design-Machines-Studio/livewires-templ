"""Regression checks for generated files that only appear after generation."""
import pathlib
import tempfile
import unittest
from unittest import mock

import run
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


class BrowserSourceTest(unittest.TestCase):
    def test_browser_modules_change_an_already_dirty_fingerprint(self):
        with tempfile.TemporaryDirectory() as directory, mock.patch.object(run, "CHECKOUT", pathlib.Path(directory)):
            checkout = pathlib.Path(directory)
            gallery = checkout / "examples/gallery"
            (gallery / "tests").mkdir(parents=True)
            (gallery / "assets.lock.json").write_text('{"producerCommit":"fixture"}')
            modules = ["examples/gallery/playwright.config.mjs", "examples/gallery/tests/gallery.spec.mjs"]
            for name in modules:
                (checkout / name).write_text("export default {};\n")
            run.command(["git", "init", "-q"])
            run.command(["git", "add", "."])
            run.command(["git", "-c", "user.name=Gallery test", "-c", "user.email=gallery@example.invalid",
                         "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null", "commit", "-qm", "Fixture"])
            (checkout / "unrelated.txt").write_text("Already dirty")
            before, hashes = run.source()
            self.assertTrue(before["dirty"])
            for name in modules:
                self.assertIn(name, hashes)
                with (checkout / name).open("a") as module:
                    module.write("// Changed review program\n")
                after, _ = run.source()
                self.assertTrue(after["dirty"])
                self.assertNotEqual(before["fingerprint"], after["fingerprint"])
                before = after


if __name__ == "__main__":
    unittest.main()
