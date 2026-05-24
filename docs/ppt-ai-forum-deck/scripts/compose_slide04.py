"""Compose slide 04: PC + mobile UI on Xu Cen style black canvas."""
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

ROOT = Path(__file__).resolve().parents[1]
W, H = 2560, 1440
BG = (10, 10, 10)
WHITE = (250, 250, 250)
GRAY = (140, 140, 140)


def load_font(size: int, bold: bool = False):
    candidates = [
        r"C:\Windows\Fonts\msyhbd.ttc" if bold else r"C:\Windows\Fonts\msyh.ttc",
        r"C:\Windows\Fonts\simhei.ttf",
    ]
    for path in candidates:
        try:
            return ImageFont.truetype(path, size)
        except OSError:
            continue
    return ImageFont.load_default()


def fit(img: Image.Image, max_w: int, max_h: int) -> Image.Image:
    ratio = min(max_w / img.width, max_h / img.height)
    nw, nh = int(img.width * ratio), int(img.height * ratio)
    return img.resize((nw, nh), Image.Resampling.LANCZOS)


def rounded_mask(size, radius: int) -> Image.Image:
    mask = Image.new("L", size, 0)
    draw = ImageDraw.Draw(mask)
    draw.rounded_rectangle((0, 0, size[0] - 1, size[1] - 1), radius=radius, fill=255)
    return mask


def paste_rounded(canvas: Image.Image, img: Image.Image, xy, radius: int = 28):
    out = Image.new("RGBA", img.size, (0, 0, 0, 0))
    out.paste(img, (0, 0))
    out.putalpha(rounded_mask(img.size, radius))
    base = canvas.convert("RGBA")
    base.paste(out, xy, out)
    canvas.paste(base.convert("RGB"))


def main():
    pc = Image.open(ROOT / "assets" / "pc-demo.png").convert("RGBA")
    mobile = Image.open(ROOT / "assets" / "mobile-demo.png").convert("RGBA")

    canvas = Image.new("RGB", (W, H), BG)
    draw = ImageDraw.Draw(canvas)

    title_font = load_font(118, bold=True)
    sub_font = load_font(36)
    label_font = load_font(28)

    title = "一端"
    sub = "PC 演示登录  ·  手机社区浏览"
    tw = draw.textlength(title, font=title_font)
    draw.text(((W - tw) / 2, 72), title, font=title_font, fill=WHITE)
    sw = draw.textlength(sub, font=sub_font)
    draw.text(((W - sw) / 2, 210), sub, font=sub_font, fill=GRAY)

    content_top = 300
    content_h = H - content_top - 80

    pc_fit = fit(pc, 1180, content_h)
    mob_fit = fit(mobile, 520, content_h - 40)

    pc_x = 180
    pc_y = content_top + (content_h - pc_fit.height) // 2
    mob_x = W - mob_fit.width - 220
    mob_y = content_top + (content_h - mob_fit.height) // 2 - 20

    paste_rounded(canvas, pc_fit, (pc_x, pc_y), radius=32)
    paste_rounded(canvas, mob_fit, (mob_x, mob_y), radius=40)

    draw.text((pc_x, pc_y - 46), "PC", font=label_font, fill=GRAY)
    draw.text((mob_x, mob_y - 46), "手机", font=label_font, fill=GRAY)

    out = ROOT / "origin_image" / "slide_03.png"
    out.parent.mkdir(parents=True, exist_ok=True)
    canvas.save(out, quality=95)
    print(out)


if __name__ == "__main__":
    main()
