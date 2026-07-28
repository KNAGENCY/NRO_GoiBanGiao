using Game1;

namespace NroE2E
{
    public static class NroUiOverlay
    {
        private static Image avatarFrame;
        private static Image dialogFrame;
        private static Image disabledButton;
        private static Image itemSlot;
        private static Image selectedItemSlot;

        public static bool PaintDialogFrame(
            mGraphics graphics,
            int x,
            int y,
            int width,
            int height)
        {
            Image image = Load(
                ref dialogFrame,
                "/naruto/ui/dialog-frame.png");
            if (image == null)
            {
                return false;
            }
            graphics.drawImageScale(image, x, y, width, height, 0);
            return true;
        }

        public static bool PaintAvatarFrame(
            mGraphics graphics,
            int x,
            int y)
        {
            Image image = Load(
                ref avatarFrame,
                "/naruto/ui/avatar-frame.png");
            if (image == null)
            {
                return false;
            }
            graphics.drawImage(image, x, y, StaticObj.BOTTOM_LEFT);
            return true;
        }

        public static bool PaintButtonDisabled(
            mGraphics graphics,
            int x,
            int y,
            int width,
            int height)
        {
            return PaintScaled(
                graphics,
                ref disabledButton,
                "/naruto/ui/button-disabled.png",
                x,
                y,
                width,
                height);
        }

        public static bool PaintItemSlot(
            mGraphics graphics,
            int x,
            int y,
            int width,
            int height)
        {
            return PaintScaled(
                graphics,
                ref itemSlot,
                "/naruto/ui/item-slot.png",
                x,
                y,
                width,
                height);
        }

        public static bool PaintSelectedItemSlot(
            mGraphics graphics,
            int x,
            int y,
            int width,
            int height)
        {
            return PaintScaled(
                graphics,
                ref selectedItemSlot,
                "/naruto/ui/item-slot-selected.png",
                x,
                y,
                width,
                height);
        }

        private static bool PaintScaled(
            mGraphics graphics,
            ref Image cached,
            string key,
            int x,
            int y,
            int width,
            int height)
        {
            Image image = Load(ref cached, key);
            if (image == null)
            {
                return false;
            }
            graphics.drawImageScale(image, x, y, width, height, 0);
            return true;
        }

        private static Image Load(ref Image cached, string key)
        {
            if (cached == null)
            {
                cached = GameCanvas.loadImage(key);
            }
            return cached;
        }
    }
}
