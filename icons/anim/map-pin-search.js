import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const pin = el.querySelector("g:not(.search)");
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const lens = el.querySelector("g.search");
  animate(lens, { scale: [1.35, 1], translateX: [0, -2, 1, 0], translateY: [0, -2, -1, 0] }, { duration: τmax });
}
