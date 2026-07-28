import unittest


class VisualContractTests(unittest.TestCase):
    def test_background_requires_all_four_layers(self):
        from assetpack.visual import VisualError, validate_background_layers

        with self.assertRaisesRegex(VisualError, "ASSET-BG-005"):
            validate_background_layers({"b00": object(), "b01": object(), "b02": object()})
