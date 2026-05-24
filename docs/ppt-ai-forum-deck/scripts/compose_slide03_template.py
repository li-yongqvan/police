"""Slide 03 template: Xu Cen style, reserved slots for PC + mobile screenshots."""
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

ROOT = Path(__file__).resolve().parents[1]
W, H = 2560, 1440
BG = (10, 10, 10)
SLOT_FILL = (22, 22, 24)
SLOT_BORDER = (55, 55, 58)
WHITE = (250, 250, 250)
GRAY = (140, 140, 140)
MUTED = (90, 90, 92)


def load_font(size: int, bold: bool = False):
    for path in (
        [r"C:\Windows\Fonts\msyhbd.ttc", r"C:\Windows\Fonts\msyh.ttc"]
        if bold
        else [r"C:\Windows\Fonts\msyh.ttc", r"C:\Windows\Fonts\simhei.ttf"]
    ):
        try:
            return ImageFont.truetype(path, size)
        except OSError:
            continue
    return ImageFont.load_default()


def rounded_rect(draw, xy, radius: int, fill, outline=None, width: int = 2):
    draw.rounded_rectangle(xy, radius=radius, fill=fill, outline=outline, width=width)


def main():
    canvas = Image.new("RGB", (W, H), BG)
    draw = ImageDraw.Draw(canvas)

    title_font = load_font(118, bold=True)
    sub_font = load_font(36)
    label_font = load_font(28)
    title, sub = "一端", "PC 演示登录  ·  手机社区浏览"
    tw = draw.textlength(title, font=title_font)
    draw.text(((W - tw) / 2, 72), title, font=title_font, fill=WHITE)
    sw = draw.textlength(sub, font=sub_font)
    draw.text(((W - sw) / 2, 210), sub, font=sub_font, fill=GRAY)

    content_top = 310
    content_h = H - content_top - 72

    # PC slot — landscape
    pc_w, pc_h = 1180, int(content_h * 0.92)
    pc_x, pc_y = 180, content_top + (content_h - pc_h) // 2
    rounded_rect(draw, (pc_x, pc_y, pc_x + pc_w, pc_y + pc_h), 32, SLOT_FILL, SLOT_BORDER, 2)

    # Mobile slot — portrait
    mob_w, mob_h = 520, int(content_h * 0.95)
    mob_x = W - mob_w - 220
    mob_y = content_top + (content_h - mob_h) // 2
    rounded_rect(draw, (mob_x, mob_y, mob_x + mob_w, mob_y + mob_h), 40, SLOT_FILL, SLOT_BORDER, 2)

    draw.text((pc_x, pc_y - 46), "PC", font=label_font, fill=GRAY)
    draw.text((mob_x, mob_y - 46), "手机", font=label_font, fill=GRAY)

    out = ROOT / "origin_image" / "slide_03.png"
    out.parent.mkdir(parents=True, exist_ok=True)
    canvas.save(out, quality=95)
    print(out)


if __name__ == "__main__":
    main()
