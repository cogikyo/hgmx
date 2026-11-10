import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const holes = el.querySelector('g[transform-origin="10px 14px"]');
  animate(holes, { rotate: [0, 360] }, { duration: τmid });
  const pin = el.querySelector('g[transform-origin="17.95px 22px"]');
  animate(pin, { rotate: [0, 12, -6, 0] }, { duration: τmid, delay: τmin });
}
