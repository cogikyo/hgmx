import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bars = ["M2", "M6", "M10", "M14", "M18", "M22"].map((x) => el.querySelector(`path[d^="${x} "]`));
  animate(bars, { scaleY: [1, 0.6, 1.2, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const x = el.querySelectorAll('path[d^="m17 17"], path[d^="m22 17"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
