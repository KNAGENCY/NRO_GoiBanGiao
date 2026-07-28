"""Kiểm tra closure nền b00..b03 trước khi render contact/full-map."""


class VisualError(ValueError):
    pass


def validate_background_layers(layers) -> None:
    if set(layers) != {"b00", "b01", "b02", "b03"}:
        raise VisualError("ASSET-BG-005: background phải đủ b00..b03")
